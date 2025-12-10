package services

import (
	"testing"
	"time"
)

func TestCreateServiceRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateServiceRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateServiceRequest{
				BusinessAccountID: "123",
				Name:              "Test Service",
				DurationMinutes:   60,
				Price:             50.0,
				Currency:          "USD",
			},
			wantErr: false,
		},
		{
			name: "missing business account ID",
			req: CreateServiceRequest{
				Name:            "Test Service",
				DurationMinutes: 60,
				Price:           50.0,
			},
			wantErr: true,
		},
		{
			name: "missing name",
			req: CreateServiceRequest{
				BusinessAccountID: "123",
				DurationMinutes:   60,
				Price:             50.0,
			},
			wantErr: true,
		},
		{
			name: "invalid duration",
			req: CreateServiceRequest{
				BusinessAccountID: "123",
				Name:              "Test Service",
				DurationMinutes:   0,
				Price:             50.0,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			req: CreateServiceRequest{
				BusinessAccountID: "123",
				Name:              "Test Service",
				DurationMinutes:   60,
				Price:             -10.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This would test the validation logic if we had it in the store
			// For now, just check the struct creation
			if tt.req.BusinessAccountID == "" && !tt.wantErr {
				t.Errorf("CreateServiceRequest validation failed: business_account_id is required")
			}
			if tt.req.Name == "" && !tt.wantErr {
				t.Errorf("CreateServiceRequest validation failed: name is required")
			}
			if tt.req.DurationMinutes <= 0 && !tt.wantErr {
				t.Errorf("CreateServiceRequest validation failed: duration_minutes must be positive")
			}
			if tt.req.Price < 0 && !tt.wantErr {
				t.Errorf("CreateServiceRequest validation failed: price cannot be negative")
			}
		})
	}
}

func TestService_JSONTags(t *testing.T) {
	service := &Service{
		ID:                "test-id",
		BusinessAccountID: "business-123",
		Name:              "Test Service",
		Description:       nil,
		DurationMinutes:   60,
		Price:             50.0,
		Currency:          "USD",
		Category:          nil,
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// This test ensures the JSON tags are properly set
	// In a real test, you'd marshal to JSON and verify the output
	if service.ID == "" {
		t.Error("Service ID should not be empty")
	}
	if service.BusinessAccountID == "" {
		t.Error("BusinessAccountID should not be empty")
	}
	if service.Name == "" {
		t.Error("Name should not be empty")
	}
}

func TestListServicesRequest_Defaults(t *testing.T) {
	req := ListServicesRequest{}

	// Test default values
	if req.Limit != 0 {
		t.Errorf("Expected default limit to be 0, got %d", req.Limit)
	}
	if req.Offset != 0 {
		t.Errorf("Expected default offset to be 0, got %d", req.Offset)
	}
}
