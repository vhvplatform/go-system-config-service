package repository

import (
	"context"
	"time"

	"github.com/vhvplatform/go-system-config-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	watchSubscriptionCollection = "watch_subscriptions"
)

// WatchRepository handles watch subscription persistence
type WatchRepository struct {
	db *mongo.Database
}

// NewWatchRepository creates a new watch repository
func NewWatchRepository(db *mongo.Database) *WatchRepository {
	return &WatchRepository{db: db}
}

// Create creates a new watch subscription
func (r *WatchRepository) Create(ctx context.Context, subscription *domain.WatchSubscription) error {
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()
	subscription.FailureCount = 0

	result, err := r.db.Collection(watchSubscriptionCollection).InsertOne(ctx, subscription)
	if err != nil {
		return err
	}
	subscription.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a subscription by ID
func (r *WatchRepository) FindByID(ctx context.Context, id string) (*domain.WatchSubscription, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var subscription domain.WatchSubscription
	err = r.db.Collection(watchSubscriptionCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&subscription)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &subscription, err
}

// FindBySubscriberID finds a subscription by subscriber ID
func (r *WatchRepository) FindBySubscriberID(ctx context.Context, subscriberID string) (*domain.WatchSubscription, error) {
	var subscription domain.WatchSubscription
	err := r.db.Collection(watchSubscriptionCollection).FindOne(
		ctx,
		bson.M{"subscriber_id": subscriberID},
	).Decode(&subscription)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &subscription, err
}

// Update updates a subscription
func (r *WatchRepository) Update(ctx context.Context, subscription *domain.WatchSubscription) error {
	subscription.UpdatedAt = time.Now()

	filter := bson.M{"_id": subscription.ID}
	update := bson.M{"$set": subscription}

	_, err := r.db.Collection(watchSubscriptionCollection).UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a subscription
func (r *WatchRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(watchSubscriptionCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// List lists all subscriptions with pagination
func (r *WatchRepository) List(ctx context.Context, page, perPage int) ([]*domain.WatchSubscription, int64, error) {
	filter := bson.M{}

	// Count total
	total, err := r.db.Collection(watchSubscriptionCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate pagination
	skip := int64((page - 1) * perPage)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(perPage)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.db.Collection(watchSubscriptionCollection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var subscriptions []*domain.WatchSubscription
	if err = cursor.All(ctx, &subscriptions); err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

// GetActiveSubscriptions gets all active subscriptions
func (r *WatchRepository) GetActiveSubscriptions(ctx context.Context) ([]*domain.WatchSubscription, error) {
	filter := bson.M{"status": "active"}

	cursor, err := r.db.Collection(watchSubscriptionCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []*domain.WatchSubscription
	if err = cursor.All(ctx, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// IncrementFailureCount increments the failure count for a subscription
func (r *WatchRepository) IncrementFailureCount(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$inc": bson.M{"failure_count": 1}}

	_, err = r.db.Collection(watchSubscriptionCollection).UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)
	return err
}

// ResetFailureCount resets the failure count for a subscription
func (r *WatchRepository) ResetFailureCount(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"failure_count": 0,
			"last_notified": now,
		},
	}

	_, err = r.db.Collection(watchSubscriptionCollection).UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)
	return err
}
