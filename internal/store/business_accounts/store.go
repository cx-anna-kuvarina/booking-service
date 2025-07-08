package business_accounts

import (
	"context"
	"encoding/json"
)

type BusinessAccount struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	BusinessType string          `json:"businessType"`
	Location     string          `json:"location"`
	Links        json.RawMessage `json:"links"`
}

type Store interface {
	CreateBusinessAccount(ctx context.Context, account *BusinessAccount, userID string) error
	UpdateBusinessAccount(ctx context.Context, account *BusinessAccount) error
	DeleteBusinessAccount(ctx context.Context, businessAccountID string) error
	UserOwnsBusinessAccount(ctx context.Context, businessAccountID, userID string) (bool, error)
}
