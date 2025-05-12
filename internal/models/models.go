package models

import (
	"time"

	"github.com/google/uuid"
)

// user models

// User represents a user entity in the application
type User struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// CreateUser represents the fields required to create a new user
type CreateUser struct {
	Email string
}

// account models

// Account represents an account entity in the application
type Account struct {
	ID            uuid.UUID  `json:"id"`
	AccountNumber string     `json:"account_number"`
	UserID        string     `json:"user_id"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

// CreateAccount represents the fields required to create a new account
type CreateAccount struct {
	AccountNumber string
	UserID        string
}

// transaction models

// TransactionPurpose is the purpose of the transaction, "credit/debit"
type TransactionPurpose string

const (
	CREDIT TransactionPurpose = "credit"
	DEBIT  TransactionPurpose = "debit"
)

// TransactionLine represents a single line of ledger entry in the system
type TransactionLine struct {
	ID            string             `json:"id"`
	AccountID     string             `json:"account_id"`
	TransactionID string             `json:"transaction_id"`
	Purpose       TransactionPurpose `json:"purpose"`
	Amount        uint64             `json:"amount"`
	CreatedAt     *time.Time         `json:"created_at"`
}

// Transaction contains all transaction lines and relevant information about a transaction
type Transaction struct {
	ID        uuid.UUID         `json:"id"`
	Reference string            `json:"reference"`
	Lines     []TransactionLine `json:"lines"`
	CreatedAt *time.Time        `json:"created_at"`
}

// CreateTransactionLine represents the required fields needed to create a line of transaction
type CreateTransactionLine struct {
	AccountID     string
	TransactionID string
	Purpose       TransactionPurpose
	Amount        uint64
}

// CreateTransaction holds the information needed to create a new transaction in the system. The transactions lines should be 2 or more
type CreateTransaction struct {
	Lines []CreateTransactionLine
}

// APIResponse is the standard application response for both success and error messages
type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
