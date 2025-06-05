package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrshabel/sgbank/internal/models"
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB, logger *slog.Logger) *TransactionRepository {
	return &TransactionRepository{db: db, logger: logger}
}

// GetTx returns a database transaction that can be used for all operations
func (r *TransactionRepository) GetTx(ctx context.Context) (tx *sql.Tx, err error) {
	tx, err = r.db.BeginTx(ctx, &sql.TxOptions{})
	return tx, err
}

// CreateTransaction adds a new transaction to the database
func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx *sql.Tx, data *models.CreateTransaction) (*models.Transaction, error) {
	query := `
		INSERT INTO transactions (reference)
		VALUES ($1, $2, $3)
		RETURNING id, reference, created_at
	`

	// retrieve transaction details
	var transaction models.Transaction

	// create transaction
	if err := tx.QueryRowContext(ctx, query, data.Reference).Scan(&transaction.ID, &transaction.Reference, &transaction.CreatedAt); err != nil {
		return nil, err
	}

	// create transaction lines
	query = `
		INSERT INTO transaction_lines (account_id, transaction_id, purpose, amount)
		VALUES ($1, $2, $3)
		RETURNING id, account_id, transaction_id, purpose, amount, created_at
	`

	// create transaction lines
	rows, err := r.db.QueryContext(ctx, query, transaction.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var line models.TransactionLine
		if err := rows.Scan(&line.ID, &line.AccountID, &line.TransactionID, &line.Purpose, &line.Amount, &line.CreatedAt); err != nil {
			return nil, err
		}
		transaction.Lines = append(transaction.Lines, line)
	}
	return &transaction, nil
}

// GetBalanceByAccountID retrieves an non-deleted transaction by their ID
func (r *TransactionRepository) GetBalanceByAccountID(ctx context.Context, acctID uuid.UUID) (uint64, error) {
	query := `
	 SELECT 
	 SUM(CASE WHEN purpose = $1 THEN amount END) AS credit_balance,
	 SUM(CASE WHEN purpose = $2 THEN amount END) AS debit_balance
	 FROM transaction_lines
	 WHERE account_id = $1
	 `

	// retrieve balances
	var creditBalance, debitBalance uint64
	if err := r.db.QueryRowContext(ctx, query, models.CREDIT, models.DEBIT, acctID).Scan(creditBalance, debitBalance); err != nil {
		return 0, err
	}

	normalBalance := creditBalance - debitBalance
	return normalBalance, nil
}

// GetTransactionByID retrieves an non-deleted transaction by their ID
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	query := `
	 SELECT 
	 id,
	 reference,
	 created_at,
	 lines.id AS line_id,
	 lines.account_id AS line_account_id,
	 lines.purpose AS line_purpose,
	 lines.amount AS line_amount,
	 lines.created_at AS line_created_at
	 FROM transactions
	 JOIN transaction_lines
	 ON transactions.id = transaction_lines.transaction_id
	 WHERE id = $1
	 `

	//  retrieve transaction with lines
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	var transaction models.Transaction

	for rows.Next() {
		var line models.TransactionLine
		if err := rows.Scan(&transaction.ID, &transaction.Reference, &transaction.CreatedAt, &line.ID, &line.AccountID, &line.Purpose, &line.Amount, &line.CreatedAt); err != nil {
			return nil, err
		}
		transaction.Lines = append(transaction.Lines, line)
	}
	return &transaction, nil
}

// GetTransactionByUserId retrieves all non-deleted transactions belonging to a user
func (r *TransactionRepository) GetTransactionsByAccountID(ctx context.Context, accountId uuid.UUID) ([]*models.Transaction, error) {
	query := `
	 SELECT 
	 id,
	 reference,
	 created_at,
	 lines.id AS line_id,
	 lines.account_id AS line_account_id,
	 lines.purpose AS line_purpose,
	 lines.amount AS line_amount,
	 lines.created_at AS line_created_at
	 FROM transactions
	 JOIN transaction_lines AS lines
	 ON transactions.id = transaction_lines.transaction_id
	 WHERE lines.account_id = $1
	 ORDER BY created_at DESC
	 `

	rows, err := r.db.QueryContext(ctx, query, accountId)
	if err != nil {
		return nil, err
	}
	// hold group transaction lines in order of how they were returned from the db as grouped by their id
	groupedTx := make(map[uuid.UUID]*models.Transaction)
	var orderedTx []uuid.UUID

	for rows.Next() {
		var transaction models.Transaction
		var line models.TransactionLine

		if err := rows.Scan(&transaction.ID, &transaction.Reference, &transaction.CreatedAt, &line.ID, &line.AccountID, &line.Purpose, &line.Amount, &line.CreatedAt); err != nil {
			return nil, err
		}

		// add new line or append line to existing transaction
		existingTx, exists := groupedTx[transaction.ID]
		if !exists {
			transaction.Lines = append(transaction.Lines, line)
			groupedTx[transaction.ID] = &transaction
		} else {
			existingTx.Lines = append(existingTx.Lines, line)
		}
		orderedTx = append(orderedTx, transaction.ID)
	}

	// order transactions
	transactions := make([]*models.Transaction, 0, len(orderedTx))
	for _, id := range orderedTx {
		transactions = append(transactions, groupedTx[id])
	}

	return transactions, nil
}
