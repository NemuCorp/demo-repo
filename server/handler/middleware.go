package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/NemuCorp/demo-repo/server/db"
	"github.com/NemuCorp/demo-repo/server/logger"
	"github.com/NemuCorp/demo-repo/server/myerrors"
)

func AuthMiddleware(authDB *db.AuthDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			JSONError(c, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		token := header[7:]
		sessionHash := sha256Sum(token)

		session, err := authDB.GetSession(sessionHash)
		if errors.Is(err, myerrors.ErrSessionExpired) {
			JSONError(c, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error())
			c.Abort()
			return
		}
		if err != nil {
			logger.Error.Println("session lookup failed:", err)
			JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
			c.Abort()
			return
		}

		user, err := authDB.GetUserByID(session.UserID)
		if err != nil {
			logger.Error.Println("user lookup failed:", err)
			JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("is_admin", user.IsAdmin)
		c.Next()
	}
}

// OptionalAuthMiddleware sets user context if a valid Bearer token is provided,
// but does not fail or abort the request if no token or an invalid token is given.
func OptionalAuthMiddleware(authDB *db.AuthDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.Next()
			return
		}

		token := header[7:]
		sessionHash := sha256Sum(token)

		session, err := authDB.GetSession(sessionHash)
		if err != nil {
			c.Next()
			return
		}

		user, err := authDB.GetUserByID(session.UserID)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("is_admin", user.IsAdmin)
		c.Next()
	}
}

func sha256Sum(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, ok := IsAdmin(c)
		if !ok || !isAdmin {
			JSONError(c, http.StatusForbidden, myerrors.ErrForbidden.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

func IsAdmin(c *gin.Context) (bool, bool) {
	val, exists := c.Get("is_admin")
	if !exists {
		return false, false
	}
	isAdmin, ok := val.(bool)
	return isAdmin, ok
}
