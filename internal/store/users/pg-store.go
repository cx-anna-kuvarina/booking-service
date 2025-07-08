package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const tableName = "users"

type PgStore struct {
	readPool  *pgxpool.Pool
	writePool *pgxpool.Pool
}

// type check
var _ Store = NewStore(nil, nil)

func NewStore(readPool, writePool *pgxpool.Pool) *PgStore {
	return &PgStore{
		readPool:  readPool,
		writePool: writePool,
	}
}

func (s *PgStore) CreateUser(ctx context.Context, u *User) error {
	if u.ID == "" {
		id, _ := uuid.NewUUID()
		u.ID = id.String()
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, username, email, firstname, lastname, phone)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, tableName)

	_, err := s.writePool.Exec(ctx, query, u.ID, u.Username, u.Email, u.FirstName, u.LastName, u.Phone)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *PgStore) GetByID(ctx context.Context, userID string) (*User, error) {
	query := fmt.Sprintf(`SELECT FROM %s id, username, email, firstname, lastname, phone where id = $1`, tableName)

	var user User
	err := s.readPool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (s *PgStore) GetUserIdByEmail(ctx context.Context, email string) (string, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE email = $1`, tableName)

	var id string
	err := s.readPool.QueryRow(ctx, query, email).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get user id by email: %w", err)
	}

	return id, nil
}

func (s *PgStore) GetUser(ctx context.Context, userID string) (*User, error) {
	query := fmt.Sprintf(`SELECT id, username, email, firstname, lastname, phone FROM %s WHERE id = $1`, tableName)

	var user User
	err := s.readPool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *PgStore) UpdateUser(ctx context.Context, u *User) error {
	query := fmt.Sprintf(`UPDATE %s SET username = $1, firstname = $2, lastname = $3, phone = $4 WHERE id = $5`, tableName)

	_, err := s.writePool.Exec(ctx, query, u.Username, u.FirstName, u.LastName, u.Phone, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
