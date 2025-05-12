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
	"github.com/mrshabel/sgbank/internal/utils"
)

// errors
var (
	ErrAccountExists   = errors.New("account already exist")
	ErrAccountNotFound = errors.New("account not found")
)

// AccountHandler contains http handlers for account-related endpoints
type AccountHandler struct {
	accountRepo *repository.AccountRepository
	logger      *slog.Logger
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountRepo *repository.AccountRepository, logger *slog.Logger) *AccountHandler {
	return &AccountHandler{
		accountRepo: accountRepo,
		logger:      logger,
	}
}

// create account

// CreateAccountRequest represents the account request payload
type CreateAccountRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

// CreateAccount handles new account creation
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var body CreateAccountRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		// log error
		h.logError("invalid request body", err)
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: err.Error(),
		})
		return
	}

	// TODO: generate unique account number
	accountNumber := utils.GenerateAccountNumber(10)

	account, err := h.accountRepo.CreateAccount(c.Request.Context(), &models.CreateAccount{AccountNumber: accountNumber, UserID: body.UserID})
	if err != nil {
		// log error
		h.logError("failed to create account", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "failed to create account",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Account created successfully",
		Data:    account,
	})
}

// GetAccountURI represents the path params of the GetAccount request
type GetAccountURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetAccount handles account retrieval
func (h *AccountHandler) GetAccount(c *gin.Context) {
	var params GetAccountURI
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
	account, err := h.accountRepo.GetAccountByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logError("account not found", err)
			c.JSON(http.StatusNotFound, models.APIResponse{
				Message: "Account not found",
			})
			return
		}

		// log error
		h.logError("failed to retrieve account", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Failed to retrieve account",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Account retrieved successfully",
		Data:    account,
	})
}

// GetAccountURI represents the path params of the GetAccount request
type GetUserAccountsQuery struct {
	UserID string `form:"user_id" binding:"required,uuid"`
}

// GetUserAccounts handles accounts retrieval for a specific user
func (h *AccountHandler) GetUserAccounts(c *gin.Context) {
	var params GetUserAccountsQuery
	if err := c.ShouldBindQuery(&params); err != nil {
		// log error
		h.logError("invalid request params", err)
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: err.Error(),
		})
		return
	}

	// parse uuid
	userID, _ := uuid.Parse(params.UserID)
	accounts, err := h.accountRepo.GetAccountsByUserID(c.Request.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logError("user accounts not found", err)
			c.JSON(http.StatusNotFound, models.APIResponse{
				Message: "User Accounts not found",
			})
			return
		}

		// log error
		h.logError("failed to retrieve user accounts", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Failed to retrieve user accounts",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Accounts retrieved successfully",
		Data:    accounts,
	})
}

// DisableAccountURI represents the path params of the GetAccount request
type DisableAccountURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// DisableAccount handles the deletion of an account
func (h *AccountHandler) DisableAccount(c *gin.Context) {
	var params GetAccountURI
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
	account, err := h.accountRepo.DisableAccountByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logError("account not found", err)
			c.JSON(http.StatusNotFound, models.APIResponse{
				Message: "Account not found",
			})
			return
		}

		// log error
		h.logError("failed to disable account", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Failed to disable account",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Account disabled successfully",
		Data:    account,
	})
}

func (h *AccountHandler) logError(message string, err error) {
	h.logger.Error(message, "error", err)
}

// RegisterAccountHandlers adds all the handler methods to the provided http router
func RegisterAccountHandlers(h *AccountHandler, router *gin.Engine, logger *slog.Logger) {
	r := router.Group("/accounts")
	r.POST("", h.CreateAccount)
	r.GET("", h.GetUserAccounts)
	r.GET("/:id", h.GetAccount)
	r.PATCH("/:id/disable", h.DisableAccount)
}
