package bookings

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	BusinessID string    `json:"business_id"`
	ServiceID  string    `json:"service_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateBookingRequest struct {
	UserID     string    `json:"user_id"`
	BusinessID string    `json:"business_id"`
	ServiceID  string    `json:"service_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

type Store interface {
	GetBooking(ctx context.Context, id string) (*Booking, error)
	CreateBooking(ctx context.Context, req CreateBookingRequest) (*Booking, error)
}

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) GetBooking(ctx context.Context, id string) (*Booking, error) {
	query := `
		SELECT id, user_id, business_id, service_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`

	var booking Booking
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.BusinessID,
		&booking.ServiceID,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &booking, nil
}

func (s *store) CreateBooking(ctx context.Context, req CreateBookingRequest) (*Booking, error) {
	query := `
		INSERT INTO bookings (
			id, user_id, business_id, service_id, start_time, end_time, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id, user_id, business_id, service_id, start_time, end_time, status, created_at, updated_at
	`

	now := time.Now()
	booking := &Booking{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		BusinessID: req.BusinessID,
		ServiceID:  req.ServiceID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Status:     "pending", // Initial status
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := s.db.QueryRowContext(ctx, query,
		booking.ID,
		booking.UserID,
		booking.BusinessID,
		booking.ServiceID,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		booking.CreatedAt,
		booking.UpdatedAt,
	).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.BusinessID,
		&booking.ServiceID,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return booking, nil
}
