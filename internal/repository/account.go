package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrshabel/sgbank/internal/models"
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *sql.DB, logger *slog.Logger) *AccountRepository {
	return &AccountRepository{db: db, logger: logger}
}

// GetTx returns a database transaction that can be used for all operations
func (r *AccountRepository) GetTx(ctx context.Context) (tx *sql.Tx, err error) {
	tx, err = r.db.BeginTx(ctx, &sql.TxOptions{})
	return tx, err
}

// CreateAccount adds a new account to the database
func (r *AccountRepository) CreateAccount(ctx context.Context, data *models.CreateAccount) (*models.Account, error) {
	query := `
		INSERT INTO accounts (account_number, user_id)
		VALUES ($1, $2)
		RETURNING id, account_number, user_id, created_at, updated_at, deleted_at
	`

	// retrieve account details
	var account models.Account
	if err := r.db.QueryRowContext(ctx, query, data.AccountNumber, data.UserID).Scan(&account.ID, &account.AccountNumber, &account.UserID, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt); err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAccountByID retrieves a non-deleted account by their ID
func (r *AccountRepository) GetAccountByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	query := `
	 SELECT id, account_number, user_id, created_at, updated_at FROM accounts
	 WHERE deleted_at IS NULL AND id = $1
	 `
	var account models.Account
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&account.ID, &account.AccountNumber, &account.UserID, &account.CreatedAt, &account.UpdatedAt); err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAccountByAcctNumbers retrieves all non-deleted accounts belonging associated with the given account numbers
func (r *AccountRepository) GetAccountsByAcctNumbers(ctx context.Context, acctNums []string) ([]*models.Account, error) {
	query := `
	 SELECT id, account_number, user_id, created_at, updated_at FROM accounts
	 WHERE deleted_at IS NULL AND account_number IN ($1)
	 `

	var accounts []*models.Account
	rows, err := r.db.QueryContext(ctx, query, acctNums)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.AccountNumber, &account.UserID, &account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// GetAccountByUserId retrieves all non-deleted accounts belonging to a user
func (r *AccountRepository) GetAccountsByUserID(ctx context.Context, userId uuid.UUID) ([]*models.Account, error) {
	query := `
	 SELECT id, account_number, user_id, created_at, updated_at FROM accounts
	 WHERE deleted_at IS NULL AND user_id  = $1
	 `

	var accounts []*models.Account
	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.AccountNumber, &account.UserID, &account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// DisableAccountByID marks an account as deleted in the system
func (r *AccountRepository) DisableAccountByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	query := `
	 UPDATE accounts
	 SET deleted_at = NOW()
	 WHERE deleted_at IS NULL AND id  = $1
	 RETURNING id, account_number, user_id, created_at, updated_at, deleted_at
	 `

	var account models.Account
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&account.ID, &account.AccountNumber, &account.UserID, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt); err != nil {
		return nil, err
	}

	return &account, nil
}
