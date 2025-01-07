package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Blue-Flamingo-Ventures/plateshare-backend/pkg/oidc"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	oidcClientID, ok := os.LookupEnv("OIDC_CLIENT_ID")
	if !ok {
		log.Panic("Missing OIDC_CLIENT_ID")
	}

	oidcClientSecret, ok := os.LookupEnv("OIDC_CLIENT_SECRET")
	if !ok {
		log.Panic("Missing OIDC_CLIENT_SECRET")
	}

	oidcBaseUrl, ok := os.LookupEnv("OIDC_BASE_URL")
	if !ok {
		log.Panic("Missing OIDC_BASE_URL")
	}

	oidcAudience, ok := os.LookupEnv("OIDC_AUDIENCE")
	if !ok {
		log.Panic("Missing OIDC_AUDIENCE")
	}

	router = gin.New()

	router.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello Word",
		})
	})
	// Login endpoint
	router.POST("/api/login", func(c *gin.Context) {
		var loginPayload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginPayload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}

		if loginPayload.Username == "" || loginPayload.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username and password are required",
			})
			return
		}

		payload := map[string]string{
			"grant_type":    "password",
			"username":      loginPayload.Username,
			"password":      loginPayload.Password,
			"client_id":     oidcClientID,
			"client_secret": oidcClientSecret,
			"scope":         "openid",
			"audience":      oidcAudience,
		}

		tokenResponse, err := oidc.FetchIDToken(fmt.Sprintf("%s/oauth/token", oidcBaseUrl), payload)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication failed",
				"details": err.Error(),
			})
			return
		}

		claims, err := oidc.ParseJWT(tokenResponse.IDToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to parse ID token",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"causality_key":   claims.CausalityKey,
			"causality_token": claims.CausalityToken,
			"nickname":        claims.Nickname,
		})
	})

}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
