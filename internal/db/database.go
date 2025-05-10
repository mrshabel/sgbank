package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

func New(connStr string, logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// verify that database connection is still alive
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// connection pool setup
	db.SetMaxOpenConns(16)

	logger.Debug("connected to database successfully")

	// migrations setup
	migration_queries := `
		// users
		CREATE TABLE users IF NOT EXISTS (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		// accounts
		CREATE TABLE accounts IF NOT EXISTS (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);

		// transactions
		// TODO: block updates on transactions
		CREATE TABLE transactions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			reference VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		// transaction lines
		// TODO: block updates on transaction lines
		CREATE TABLE transaction_lines (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			account_id UUID NOT NULL REFERENCES accounts(id),
			transaction_id UUID NOT NULL REFERENCES transactions(id),
			purpose VARCHAR(50) NOT NULL,
			AMOUNT VARCHAR(50) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(account_id, transaction_id)
		);


	`

	if _, err := db.Exec(migration_queries); err != nil {
		return nil, err
	}

	return db, nil
}
