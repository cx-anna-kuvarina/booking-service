package business_accounts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	businessAccountsTable     = "business_accounts"
	userBusinessAccountsTable = "user_business_accounts"
)

var (
	NotFoundError = errors.New("business account not found")
)

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

func (s *PgStore) CreateBusinessAccount(ctx context.Context, account *BusinessAccount, userID string) error {
	tx, err := s.writePool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert business account
	query := fmt.Sprintf(`
		INSERT INTO %s (id, name, business_type, location, links)
		VALUES ($1, $2, $3, $4, $5)
	`, businessAccountsTable)

	_, err = tx.Exec(ctx, query, account.ID, account.Name, account.BusinessType, account.Location, account.Links)
	if err != nil {
		return fmt.Errorf("failed to create business account: %w", err)
	}

	// Create user-business account relationship
	relQuery := fmt.Sprintf(`
		INSERT INTO %s (business_account_id, user_id)
		VALUES ($1, $2)
	`, userBusinessAccountsTable)

	_, err = tx.Exec(ctx, relQuery, account.ID, userID)
	if err != nil {
		return fmt.Errorf("failed to create user-business account relationship: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PgStore) UpdateBusinessAccount(ctx context.Context, account *BusinessAccount) error {
	query := fmt.Sprintf(`
		UPDATE %s SET name = $1, business_type = $2, location = $3, links = $4 WHERE id = $5
	`, businessAccountsTable)

	_, err := s.writePool.Exec(ctx, query, account.Name, account.BusinessType, account.Location, account.Links, account.ID)
	if err != nil {
		return fmt.Errorf("failed to update business account: %w", err)
	}
	return nil
}

func (s *PgStore) DeleteBusinessAccount(ctx context.Context, businessAccountID string) error {
	tx, err := s.writePool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete user-business account relationships
	relDelQuery := fmt.Sprintf(`
		DELETE FROM %s WHERE business_account_id = $1
	`, userBusinessAccountsTable)
	_, err = tx.Exec(ctx, relDelQuery, businessAccountID)
	if err != nil {
		return fmt.Errorf("failed to delete user-business account relationships: %w", err)
	}

	// Delete business account
	accDelQuery := fmt.Sprintf(`
		DELETE FROM %s WHERE id = $1
	`, businessAccountsTable)
	_, err = tx.Exec(ctx, accDelQuery, businessAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return NotFoundError
		}
		return fmt.Errorf("failed to delete business account: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PgStore) UserOwnsBusinessAccount(ctx context.Context, businessAccountID, userID string) (bool, error) {
	query := `SELECT COUNT(1) FROM user_business_accounts WHERE business_account_id = $1 AND user_id = $2`
	var count int
	err := s.readPool.QueryRow(ctx, query, businessAccountID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *PgStore) GetBusinessAccount(ctx context.Context, businessAccountID string) (*BusinessAccount, error) {
	query := `SELECT id, name, business_type, location, links FROM business_accounts WHERE id = $1`
	var account BusinessAccount
	err := s.readPool.QueryRow(ctx, query, businessAccountID).Scan(&account.ID, &account.Name, &account.BusinessType, &account.Location, &account.Links)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFoundError
		}
		return nil, err
	}
	return &account, nil
}
