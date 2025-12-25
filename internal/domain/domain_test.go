package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCountry_Validation(t *testing.T) {
	tests := []struct {
		name    string
		country Country
		wantErr bool
	}{
		{
			name: "Valid country",
			country: Country{
				Code:   "US",
				Name:   map[string]string{"en": "United States", "vi": "Hoa Kỳ"},
				Status: "active",
			},
			wantErr: false,
		},
		{
			name: "Empty code",
			country: Country{
				Code:   "",
				Name:   map[string]string{"en": "United States"},
				Status: "active",
			},
			wantErr: true,
		},
		{
			name: "Empty name",
			country: Country{
				Code:   "US",
				Name:   map[string]string{},
				Status: "active",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.country.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAppComponent_Validation(t *testing.T) {
	tests := []struct {
		name      string
		component AppComponent
		wantErr   bool
	}{
		{
			name: "Valid component",
			component: AppComponent{
				ID:       primitive.NewObjectID(),
				TenantID: "tenant-123",
				Code:     "dashboard",
				Name:     "Dashboard",
				Status:   "active",
			},
			wantErr: false,
		},
		{
			name: "Empty tenant ID",
			component: AppComponent{
				TenantID: "",
				Code:     "dashboard",
				Name:     "Dashboard",
				Status:   "active",
			},
			wantErr: true,
		},
		{
			name: "Empty code",
			component: AppComponent{
				TenantID: "tenant-123",
				Code:     "",
				Name:     "Dashboard",
				Status:   "active",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.component.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPaginationRequest_Defaults(t *testing.T) {
	req := &PaginationRequest{
		Page:    0,
		PerPage: 0,
	}

	req.SetDefaults()

	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 30, req.PerPage)
}

func TestPaginationRequest_CustomPagination(t *testing.T) {
	req := &PaginationRequest{
		Page:    2,
		PerPage: 50,
	}

	req.SetDefaults()

	assert.Equal(t, 2, req.Page)
	assert.Equal(t, 50, req.PerPage)
}

func TestPaginationRequest_MaxPerPage(t *testing.T) {
	req := &PaginationRequest{
		Page:    1,
		PerPage: 200,
	}

	req.SetDefaults()

	assert.Equal(t, 1, req.Page)
	assert.LessOrEqual(t, req.PerPage, 100)
}
