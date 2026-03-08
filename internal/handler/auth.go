package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/piotrzalecki/budget-api/internal/repo"
	"github.com/piotrzalecki/budget-api/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user and returns a session token.
//
// @Summary Login
// @Description Authenticate with email and password, receive a session token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	req, ok := GetValidatedRequest[model.LoginRequest](c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateToken()
	if err != nil {
		h.logger.Error("failed to generate token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	expiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	session, err := h.repo.CreateSession(c.Request.Context(), repo.CreateSessionParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: sql.NullTime{Time: expiresAt, Valid: true},
	})
	if err != nil {
		h.logger.Error("failed to create session")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var expiresAtPtr *time.Time
	if session.ExpiresAt.Valid {
		t := session.ExpiresAt.Time
		expiresAtPtr = &t
	}

	c.JSON(http.StatusOK, gin.H{
		"data": model.LoginResponse{
			Token:     session.Token,
			ExpiresAt: expiresAtPtr,
			UserID:    user.ID,
			Email:     user.Email,
		},
		"error": nil,
	})
}

// Logout invalidates the current session token.
//
// @Summary Logout
// @Description Invalidate the current session token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 204
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	if err := h.repo.DeleteSession(c.Request.Context(), token); err != nil {
		h.logger.Error("failed to delete session")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
