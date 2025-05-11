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
	ErrUserExists   = errors.New("user already exist")
	ErrUserNotFound = errors.New("user not found")
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
	Email string `json:"email" binding:"required,email"`
}

// CreateUser handles new user creation
func (h *UserHandler) CreateUser(c *gin.Context) {
	var body CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		// log error
		h.logError("invalid request body", err)
		c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
			Message: err.Error(),
		})
		return
	}

	user, err := h.userRepo.CreateUser(c.Request.Context(), &models.CreateUser{Email: body.Email})
	if err != nil {
		// log error
		h.logError("failed to create user", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUserURI represents the path params of the GetUser request
type GetUserURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetUser handles user retrieval
func (h *UserHandler) GetUser(c *gin.Context) {
	var params GetUserURI
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
	user, err := h.userRepo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logError("user not found", err)
			c.JSON(http.StatusNotFound, models.APIResponse{
				Message: "User not found",
			})
			return
		}

		// log error
		h.logError("failed to retrieve user", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Failed to retrieve user",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "User retrieved successfully",
		Data:    user,
	})
}

func (h *UserHandler) logError(message string, err error) {
	h.logger.Error(message, "error", err)
}

// RegisterUserHandlers adds all the handler methods to the provided http router
func RegisterUserHandlers(h *UserHandler, router *gin.Engine, logger *slog.Logger) {
	r := router.Group("/users")
	r.POST("", h.CreateUser)
	r.GET("/:id", h.GetUser)
}
