package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrshabel/sgbank/internal/models"
	"github.com/mrshabel/sgbank/internal/repository"
)

// errors
var (
	ErrTransactionExists   = errors.New("transaction already exist")
	ErrTransactionNotFound = errors.New("transaction not found")
)

// TransactionHandler contains http handlers for transaction-related endpoints
type TransactionHandler struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
	logger          *slog.Logger
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository, logger *slog.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		logger:          logger,
	}
}

// create transaction

// CreateTransactionRequest represents the transaction request payload
type CreateTransactionRequest struct {
	Reference string `json:"reference" binding:"required"`
	Sender    string `json:"sender" binding:"required"`
	Recipient string `json:"recipient" binding:"required"`
	Amount    uint64 `json:"amount" binding:"required,gt=0"`
	// Purpose   string `json:"purpose" binding:"required"`
}

// validate ensures that the input correctly matches the ledger rules

// CreateTransaction handles new transaction creation
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var body CreateTransactionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.logError("invalid request body", err)
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: err.Error(),
		})
		return
	}

	// Start database transaction
	repoTx, err := h.transactionRepo.GetTx(c.Request.Context())
	if err != nil {
		h.logError("failed to obtain database transaction", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "failed to create transaction",
		})
		return
	}
	defer repoTx.Rollback()

	// validate accounts exist
	accounts, err := h.validateAccounts(c, body.Sender, body.Recipient)
	if err != nil {
		return
	}

	// create transaction lines based on transaction type
	lines, err := h.createTransactionLines(c, body, accounts)
	if err != nil {
		return
	}

	// create the transaction
	transaction, err := h.transactionRepo.CreateTransaction(c.Request.Context(), repoTx, &models.CreateTransaction{
		Reference: body.Reference,
		Lines:     lines,
	})
	if err != nil {
		h.logError("failed to create transaction", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "failed to process transaction",
		})
		return
	}

	// commit transaction
	if err := repoTx.Commit(); err != nil {
		h.logError("failed to commit transaction", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "failed to process transaction",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Transaction processed successfully",
		Data:    transaction,
	})
}

// validateAccounts retrieves the associated accounts in the system.
func (h *TransactionHandler) validateAccounts(c *gin.Context, sender, recipient string) (map[string]*models.Account, error) {
	accounts, err := h.accountRepo.GetAccountsByAcctNumbers(c.Request.Context(), []string{sender, recipient})
	if err != nil {
		h.logError("failed to retrieve related accounts", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "failed to create transaction",
		})
		return nil, err
	}

	if len(accounts) < 2 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Sender/Recipient account does not exist",
		})
		return nil, errors.New("accounts not found")
	}

	// map account numbers to accounts
	accountMap := make(map[string]*models.Account)
	for _, acct := range accounts {
		accountMap[acct.AccountNumber] = acct
	}

	return accountMap, nil
}

func (h *TransactionHandler) createTransactionLines(c *gin.Context, body CreateTransactionRequest, accounts map[string]*models.Account) ([]models.CreateTransactionLine, error) {
	var lines []models.CreateTransactionLine

	// skip balance checks for transfer from root accounts (deposits)
	if body.Sender == models.RootAccount {
		// system transactions: Credit recipient, Debit root account
		lines = []models.CreateTransactionLine{
			{
				AccountID: accounts[body.Recipient].ID,
				Purpose:   models.CREDIT,
				Amount:    body.Amount,
			},
			{
				AccountID: accounts[models.RootAccount].ID,
				Purpose:   models.DEBIT,
				Amount:    body.Amount,
			},
		}
	} else {
		// verify that sender has sufficient balance for inter-account transfers (non-root)
		senderAcct := accounts[body.Sender]
		balance, err := h.transactionRepo.GetBalanceByAccountID(c.Request.Context(), senderAcct.ID)
		if err != nil {
			h.logError("failed to retrieve sender balance", err)
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Message: "failed to process transaction",
			})
			return nil, err
		}

		if body.Amount > balance {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Message: "Insufficient balance",
			})
			return nil, errors.New("insufficient balance")
		}

		// double entry: Debit sender, Credit recipient
		lines = []models.CreateTransactionLine{
			{
				AccountID: senderAcct.ID,
				Purpose:   models.DEBIT,
				Amount:    body.Amount,
			},
			{
				AccountID: accounts[body.Recipient].ID,
				Purpose:   models.CREDIT,
				Amount:    body.Amount,
			},
		}
	}

	return lines, nil
}

// GetTransactionURI represents the path params of the GetTransaction request
type GetTransactionURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetTransaction handles transaction retrieval
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	var params GetTransactionURI
	if err := c.ShouldBindUri(&params); err != nil {
		// log error
		h.logError("invalid request params", err)
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: err.Error(),
		})
		return
	}

	// parse uuid
	id, _ := uuid.Parse(params.ID)
	transaction, err := h.transactionRepo.GetTransactionByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logError("transaction not found", err)
			c.JSON(http.StatusNotFound, models.APIResponse{
				Message: "Transaction not found",
			})
			return
		}

		// log error
		h.logError("failed to retrieve transaction", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Failed to retrieve transaction",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Transaction retrieved successfully",
		Data:    transaction,
	})
}

type GetAccountTransactionsQuery struct {
	AccountID string `form:"account_id" binding:"required,uuid"`
}

// GetAccountTransactions handles transactions retrieval for a specific user
func (h *TransactionHandler) GetAccountTransactions(c *gin.Context) {
	var params GetAccountTransactionsQuery
	if err := c.ShouldBindQuery(&params); err != nil {
		// log error
		h.logError("invalid request params", err)
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: err.Error(),
		})
		return
	}

	// parse uuid
	userID, _ := uuid.Parse(params.AccountID)
	transactions, err := h.transactionRepo.GetTransactionsByAccountID(c.Request.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logError("account transactions not found", err)
			c.JSON(http.StatusNotFound, models.APIResponse{
				Message: "Account Transactions not found",
			})
			return
		}

		// log error
		h.logError("failed to retrieve account transactions", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Failed to retrieve account transactions",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Transactions retrieved successfully",
		Data:    transactions,
	})
}

func (h *TransactionHandler) logError(message string, err error) {
	h.logger.Error(message, "error", err)
}

// RegisterTransactionHandlers adds all the handler methods to the provided http router
func RegisterTransactionHandlers(h *TransactionHandler, router *gin.Engine, logger *slog.Logger) {
	r := router.Group("/transactions")
	r.POST("", h.CreateTransaction)
	r.GET("", h.GetAccountTransactions)
	r.GET("/:id", h.GetTransaction)
}
