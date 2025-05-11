package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrshabel/sgbank/internal/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB, logger *slog.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

// GetTx returns a database transaction that can be used for all operations
func (r *UserRepository) GetTx(ctx context.Context) (tx *sql.Tx, err error) {
	tx, err = r.db.BeginTx(ctx, &sql.TxOptions{})
	return tx, err
}

// CreateUser adds a new user to the database. This is a password-less user
func (r *UserRepository) CreateUser(ctx context.Context, data *models.CreateUser) (*models.User, error) {
	query := `
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING *
	`

	// retrieve user details
	var user models.User
	if err := r.db.QueryRowContext(ctx, query, data.Email).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
	 SELECT id, email, created_at, updated_at FROM users
	 WHERE id = ($1)
	 `
	r.logger.Debug("user id", "id", id.String())
	var user models.User
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = ($1)`

	var user models.User
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

// rollbackWithErr is a helper function for rolling back a transaction
// func (r *UserRepository) rollbackWithErr(tx *sql.Tx, err error) error {
// 	if err := tx.Rollback(); err != nil {
// 		return err
// 	}
// 	return err
// }
