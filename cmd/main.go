package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-shared-go/config"
	"github.com/longvhv/saas-shared-go/logger"
	"github.com/longvhv/saas-shared-go/mongodb"
	"github.com/longvhv/saas-shared-go/redis"
	"github.com/longvhv/saas-framework-go/services/system-config-service/internal/handler"
	"github.com/longvhv/saas-framework-go/services/system-config-service/internal/repository"
	"github.com/longvhv/saas-framework-go/services/system-config-service/internal/router"
	"github.com/longvhv/saas-framework-go/services/system-config-service/internal/service"
	"github.com/longvhv/saas-framework-go/services/system-config-service/migrations"
	"go.uber.org/zap"
	grpcServer "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Sync()

	log.Info("Starting System Config Service", zap.String("environment", cfg.Environment))

	// Initialize MongoDB
	mongoClient, err := mongodb.NewClient(context.Background(), mongodb.Config{
		URI:         cfg.MongoDB.URI,
		Database:    cfg.MongoDB.Database,
		MaxPoolSize: cfg.MongoDB.MaxPoolSize,
		MinPoolSize: cfg.MongoDB.MinPoolSize,
	})
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Close(context.Background())

	// Initialize Redis
	redisClient, err := redis.NewClient(redis.Config{
		Addr:     cfg.Redis.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Seed initial data
	log.Info("Seeding initial data...")
	if err := migrations.SeedData(mongoClient.Database()); err != nil {
		log.Warn("Failed to seed data (may already exist)", zap.Error(err))
	} else {
		log.Info("Data seeded successfully")
	}

	// Initialize repositories
	appComponentRepo := repository.NewAppComponentRepository(mongoClient.Database())
	countryRepo := repository.NewCountryRepository(mongoClient.Database())
	
	// Initialize services
	appComponentService := service.NewAppComponentService(appComponentRepo, redisClient, log)
	countryService := service.NewCountryService(countryRepo, redisClient, log)
	
	// Initialize handlers
	appComponentHandler := handler.NewAppComponentHandler(appComponentService, log)
	countryHandler := handler.NewCountryHandler(countryService, log)

	// Start gRPC server
	grpcPort := os.Getenv("SYSTEM_CONFIG_SERVICE_PORT")
	if grpcPort == "" {
		grpcPort = "50055"
	}
	go startGRPCServer(log, grpcPort)

	// Start HTTP server
	httpPort := os.Getenv("SYSTEM_CONFIG_SERVICE_HTTP_PORT")
	if httpPort == "" {
		httpPort = "8085"
	}
	startHTTPServer(appComponentHandler, countryHandler, log, httpPort)
}

func startGRPCServer(log *logger.Logger, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}

	grpcSrv := grpcServer.NewServer()

	// Register health check service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcSrv, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	log.Info("gRPC server listening", zap.String("port", port))
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}

func startHTTPServer(appComponentHandler *handler.AppComponentHandler, countryHandler *handler.CountryHandler, log *logger.Logger, port string) {
	gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter(appComponentHandler, countryHandler, log)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Info("HTTP server listening", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
