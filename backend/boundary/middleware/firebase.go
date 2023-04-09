package middleware

import (
	"backend/models"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type FirebaseMiddleware struct {
	SuperUserEncoded string
	Client           *auth.Client
}

func NewFirebaseMiddleware(superUserString string, client *auth.Client) FirebaseMiddleware {
	return FirebaseMiddleware{SuperUserEncoded: superUserString, Client: client}
}

func (m *FirebaseMiddleware) AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if strings.HasPrefix(authHeader, "Basic") {
		if strings.Split(authHeader, " ")[1] != m.SuperUserEncoded {
			c.AbortWithStatusJSON(http.StatusForbidden, "invalid basic auth")
			return
		} else {
			c.Next()
			return
		}
		//base64.StdEncoding.EncodeToString(bytesconv.StringToBytes(base))
	} else if m.Client == nil {
		log.Printf("no firebase token")
		c.Set(models.FirebaseContextKey, "test")
		c.Next()
		return
	} else if strings.HasPrefix(authHeader, "Bearer") {
		token := strings.Split(authHeader, " ")[1]
		if firebaseToken, err := m.Client.VerifyIDToken(c, token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "unable to verify id token")
			return
		} else {
			c.Set(models.FirebaseContextKey, firebaseToken.UID)
			c.Next()
			return
		}
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}
