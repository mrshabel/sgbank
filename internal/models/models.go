package models

import "time"

// user models

// User represents a user entity in the application
type User struct {
	ID        string     `json:"id"`
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
	ID        string
	UserID    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// CreateAccount represents the fields required to create a new account
type CreateAccount struct {
	ID     string
	UserID string
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
	ID            string
	AccountID     string
	TransactionID string
	Purpose       TransactionPurpose
	Amount        string
	CreatedAt     *time.Time
}

// Transaction contains all transaction lines and relevant information about a transaction
type Transaction struct {
	ID        string
	Reference string
	Lines     []TransactionLine
	CreatedAt *time.Time
}

// CreateTransactionLine represents the required fields needed to create a line of transaction
type CreateTransactionLine struct {
	AccountID     string
	TransactionID string
	Purpose       TransactionPurpose
	Amount        string
}

// CreateTransaction holds the information needed to create a new transaction in the system. The transactions lines should be 2 or more
type CreateTransaction struct {
	Lines []CreateTransactionLine
}
