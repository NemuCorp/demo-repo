package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/NemuCorp/demo-repo/server/db"
	"github.com/NemuCorp/demo-repo/server/logger"
	"github.com/NemuCorp/demo-repo/server/myerrors"
)

type AuthHandler struct {
	DB *db.AuthDB
}

func NewAuthHandler(authDB *db.AuthDB) *AuthHandler {
	return &AuthHandler{DB: authDB}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error.Println("failed to hash password:", err)
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	user, err := h.DB.CreateUser(req.Email, string(passwordHash))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			JSONError(c, http.StatusConflict, myerrors.ErrEmailTaken.Error())
			return
		}
		logger.Error.Println("failed to create user:", err)
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusCreated, gin.H{"user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.DB.GetUserByEmail(req.Email)
	if err != nil {
		JSONError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		JSONError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		logger.Error.Println("failed to generate session token:", err)
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}
	plainToken := hex.EncodeToString(token)

	hash := sha256.Sum256([]byte(plainToken))
	sessionHash := hex.EncodeToString(hash[:])

	session, err := h.DB.CreateSession(user.ID, sessionHash, time.Now().Add(24*time.Hour))
	if err != nil {
		logger.Error.Println("failed to create session:", err)
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, gin.H{
		"token":   plainToken,
		"session": session,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error())
		return
	}

	if err := h.DB.DeleteUserSessions(userID); err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, gin.H{"message": "logged out"})
}
