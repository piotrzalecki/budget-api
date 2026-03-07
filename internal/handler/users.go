package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/piotrzalecki/budget-api/pkg/model"
)

// ListUsers returns all users.
//
// @Summary List users
// @Description List all users
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.UserResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.repo.ListUsers(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := make([]model.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, userToResponse(u))
	}
	c.JSON(http.StatusOK, gin.H{"data": resp, "error": nil})
}

// CreateUser creates a new user.
//
// @Summary Create user
// @Description Create a new user with hashed password
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "User details"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	req, ok := GetValidatedRequest[model.CreateUserRequest](c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	user, err := h.repo.CreateUser(c.Request.Context(), repo.CreateUserParams{
		Email:     req.Email,
		PwHash:    string(hash),
		IsService: req.IsService,
	})
	if err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": userToResponse(user), "error": nil})
}

// GetUserByID returns a user by ID.
//
// @Summary Get user
// @Description Get a user by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userToResponse(user), "error": nil})
}

// UpdateUser updates a user's email and/or password.
//
// @Summary Update user
// @Description Update a user's email and/or password
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body model.UpdateUserRequest true "Fields to update"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /users/{id} [patch]
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	req, ok := GetValidatedRequest[model.UpdateUserRequest](c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	existing, err := h.repo.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		h.logger.Error("failed to get user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	params := repo.UpdateUserParams{
		ID:     id,
		Email:  existing.Email,
		PwHash: existing.PwHash,
	}
	if req.Email != nil {
		params.Email = *req.Email
	}
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			h.logger.Error("failed to hash password", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		params.PwHash = string(hash)
	}

	user, err := h.repo.UpdateUser(c.Request.Context(), params)
	if err != nil {
		h.logger.Error("failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userToResponse(user), "error": nil})
}

// DeleteUser deletes a user and all their sessions.
//
// @Summary Delete user
// @Description Delete a user and invalidate all their sessions
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.repo.DeleteAllSessionsByUserID(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete sessions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if err := h.repo.DeleteUser(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func userToResponse(u repo.User) model.UserResponse {
	resp := model.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		IsService: u.IsService,
	}
	if u.CreatedAt.Valid {
		t := u.CreatedAt.Time
		resp.CreatedAt = &t
	}
	return resp
}
