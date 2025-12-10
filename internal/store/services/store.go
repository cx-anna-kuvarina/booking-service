package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	ID                string    `json:"id"`
	BusinessAccountID string    `json:"business_account_id"`
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	DurationMinutes   int       `json:"duration_minutes"`
	Price             float64   `json:"price"`
	Currency          string    `json:"currency"`
	Category          *string   `json:"category,omitempty"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateServiceRequest struct {
	BusinessAccountID string  `json:"business_account_id"`
	Name              string  `json:"name"`
	Description       *string `json:"description,omitempty"`
	DurationMinutes   int     `json:"duration_minutes"`
	Price             float64 `json:"price"`
	Currency          string  `json:"currency"`
	Category          *string `json:"category,omitempty"`
}

type UpdateServiceRequest struct {
	Name            *string  `json:"name,omitempty"`
	Description     *string  `json:"description,omitempty"`
	DurationMinutes *int     `json:"duration_minutes,omitempty"`
	Price           *float64 `json:"price,omitempty"`
	Currency        *string  `json:"currency,omitempty"`
	Category        *string  `json:"category,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}

type ListServicesRequest struct {
	BusinessAccountID *string `json:"business_account_id,omitempty"`
	Category          *string `json:"category,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
	Limit             int     `json:"limit"`
	Offset            int     `json:"offset"`
}

type ListServicesResponse struct {
	Services []*Service `json:"services"`
	Total    int64      `json:"total"`
}

type Store interface {
	CreateService(ctx context.Context, req CreateServiceRequest) (*Service, error)
	GetService(ctx context.Context, id string) (*Service, error)
	UpdateService(ctx context.Context, id string, req UpdateServiceRequest) (*Service, error)
	DeleteService(ctx context.Context, id string) error
	ListServices(ctx context.Context, req ListServicesRequest) (*ListServicesResponse, error)
	GetServicesByBusinessAccount(ctx context.Context, businessAccountID string) ([]*Service, error)
}

type PgStore struct {
	readPool  *pgxpool.Pool
	writePool *pgxpool.Pool
}

var _ Store = NewStore(nil, nil)

func NewStore(readPool, writePool *pgxpool.Pool) *PgStore {
	return &PgStore{
		readPool:  readPool,
		writePool: writePool,
	}
}

func (s *PgStore) CreateService(ctx context.Context, req CreateServiceRequest) (*Service, error) {
	query := `
		INSERT INTO services (
			id, business_account_id, name, description, duration_minutes, 
			price, currency, category, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, business_account_id, name, description, duration_minutes, 
			price, currency, category, is_active, created_at, updated_at
	`

	now := time.Now()
	service := &Service{
		ID:                uuid.New().String(),
		BusinessAccountID: req.BusinessAccountID,
		Name:              req.Name,
		Description:       req.Description,
		DurationMinutes:   req.DurationMinutes,
		Price:             req.Price,
		Currency:          req.Currency,
		Category:          req.Category,
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if service.Currency == "" {
		service.Currency = "USD"
	}

	err := s.writePool.QueryRow(ctx, query,
		service.ID,
		service.BusinessAccountID,
		service.Name,
		service.Description,
		service.DurationMinutes,
		service.Price,
		service.Currency,
		service.Category,
		service.IsActive,
		service.CreatedAt,
		service.UpdatedAt,
	).Scan(
		&service.ID,
		&service.BusinessAccountID,
		&service.Name,
		&service.Description,
		&service.DurationMinutes,
		&service.Price,
		&service.Currency,
		&service.Category,
		&service.IsActive,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *PgStore) GetService(ctx context.Context, id string) (*Service, error) {
	query := `
		SELECT id, business_account_id, name, description, duration_minutes, 
			price, currency, category, is_active, created_at, updated_at
		FROM services
		WHERE id = $1
	`

	var service Service
	err := s.readPool.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.BusinessAccountID,
		&service.Name,
		&service.Description,
		&service.DurationMinutes,
		&service.Price,
		&service.Currency,
		&service.Category,
		&service.IsActive,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &service, nil
}

func (s *PgStore) UpdateService(ctx context.Context, id string, req UpdateServiceRequest) (*Service, error) {
	// First get the current service to merge with updates
	currentService, err := s.GetService(ctx, id)
	if err != nil {
		return nil, err
	}
	if currentService == nil {
		return nil, errors.New("service not found")
	}

	// Build dynamic update query
	query := `
		UPDATE services SET 
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			duration_minutes = COALESCE($3, duration_minutes),
			price = COALESCE($4, price),
			currency = COALESCE($5, currency),
			category = COALESCE($6, category),
			is_active = COALESCE($7, is_active),
			updated_at = $8
		WHERE id = $9
		RETURNING id, business_account_id, name, description, duration_minutes, 
			price, currency, category, is_active, created_at, updated_at
	`

	now := time.Now()
	err = s.writePool.QueryRow(ctx, query,
		req.Name,
		req.Description,
		req.DurationMinutes,
		req.Price,
		req.Currency,
		req.Category,
		req.IsActive,
		now,
		id,
	).Scan(
		&currentService.ID,
		&currentService.BusinessAccountID,
		&currentService.Name,
		&currentService.Description,
		&currentService.DurationMinutes,
		&currentService.Price,
		&currentService.Currency,
		&currentService.Category,
		&currentService.IsActive,
		&currentService.CreatedAt,
		&currentService.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return currentService, nil
}

func (s *PgStore) DeleteService(ctx context.Context, id string) error {
	query := `DELETE FROM services WHERE id = $1`

	result, err := s.writePool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("service not found")
	}

	return nil
}

func (s *PgStore) ListServices(ctx context.Context, req ListServicesRequest) (*ListServicesResponse, error) {
	// Build dynamic WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if req.BusinessAccountID != nil {
		whereClause += " AND business_account_id = $" + string(rune(argIndex+'0'))
		args = append(args, *req.BusinessAccountID)
		argIndex++
	}

	if req.Category != nil {
		whereClause += " AND category = $" + string(rune(argIndex+'0'))
		args = append(args, *req.Category)
		argIndex++
	}

	if req.IsActive != nil {
		whereClause += " AND is_active = $" + string(rune(argIndex+'0'))
		args = append(args, *req.IsActive)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM services " + whereClause
	var total int64
	err := s.readPool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get services with pagination
	servicesQuery := `
		SELECT id, business_account_id, name, description, duration_minutes, 
			price, currency, category, is_active, created_at, updated_at
		FROM services 
		` + whereClause + `
		ORDER BY created_at DESC
		LIMIT $` + string(rune(argIndex+'0')) + ` OFFSET $` + string(rune(argIndex+1+'0'))

	args = append(args, req.Limit, req.Offset)

	rows, err := s.readPool.Query(ctx, servicesQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*Service
	for rows.Next() {
		var service Service
		err := rows.Scan(
			&service.ID,
			&service.BusinessAccountID,
			&service.Name,
			&service.Description,
			&service.DurationMinutes,
			&service.Price,
			&service.Currency,
			&service.Category,
			&service.IsActive,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}

	return &ListServicesResponse{
		Services: services,
		Total:    total,
	}, nil
}

func (s *PgStore) GetServicesByBusinessAccount(ctx context.Context, businessAccountID string) ([]*Service, error) {
	query := `
		SELECT id, business_account_id, name, description, duration_minutes, 
			price, currency, category, is_active, created_at, updated_at
		FROM services
		WHERE business_account_id = $1 AND is_active = true
		ORDER BY name ASC
	`

	rows, err := s.readPool.Query(ctx, query, businessAccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*Service
	for rows.Next() {
		var service Service
		err := rows.Scan(
			&service.ID,
			&service.BusinessAccountID,
			&service.Name,
			&service.Description,
			&service.DurationMinutes,
			&service.Price,
			&service.Currency,
			&service.Category,
			&service.IsActive,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}

	return services, nil
}
