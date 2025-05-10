package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mrshabel/sgbank/internal/repository"
)

// errors
var (
	ErrUserExists = errors.New("user already exist")
)

// UserHandler contains http handlers for user-related endpoints
type UserHandler struct {
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *repository.UserRepository, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		logger:   logger,
	}
}

// create user

// CreateUserRequest represents the user request payload
type CreateUserRequest struct {
	Email string `json:"email"`
}

// CreateUserHandler handles new user creation
func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	ctx := context.Background()
	// start the transaction
	tx, err := h.userRepo.GetTx(ctx)
}

// RegisterUserHandlers adds all the handler methods to the provided http router
func (h *UserHandler) RegisterUserHandlers(router *gin.Engine, logger *slog.Logger) {
	r := router.Group("/users")
	r.POST("/", h.CreateUserHandler)
}
