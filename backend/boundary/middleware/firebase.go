package middleware

import (
	"backend/models"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type FirebaseMiddleware struct {
	SuperUserEncoded string
	Client           *auth.Client
}

func (m *FirebaseMiddleware) AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if strings.HasPrefix(authHeader, "Basic") {
		if strings.Split(authHeader, " ")[1] != m.SuperUserEncoded {
			c.JSON(http.StatusForbidden, "invalid basic auth")
			return
		}
		//base64.StdEncoding.EncodeToString(bytesconv.StringToBytes(base))
	} else if strings.HasPrefix(authHeader, "Bearer") {
		token := strings.Split(authHeader, " ")[1]
		if firebaseToken, err := m.Client.VerifyIDToken(c, token); err != nil {
			c.JSON(http.StatusUnauthorized, "unable to verify id token")
			return
		} else {
			c.Set(models.FirebaseContextKey, firebaseToken.UID)
		}
	} else {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.Next()
}
