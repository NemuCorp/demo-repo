package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/NemuCorp/demo-repo/server/db"
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
		hash := sha256.Sum256([]byte(token))
		sessionHash := hex.EncodeToString(hash[:])

		session, err := authDB.GetSession(sessionHash)
		if err != nil {
			JSONError(c, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}
