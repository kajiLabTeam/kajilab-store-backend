// Package middleware provides middleware handlers for authentication and other HTTP middleware features.
package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ヘッダー取得
		apiKey := c.GetHeader("X-API-Key") // クライアントからAPIキーを取得

		// トークンIDとAPIキーが両方空の場合は拒否する
		if apiKey == "" {
			log.Printf("Authorization and X-API-Key are empty")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if apiKey != "" {
			// APIキーの検証
			err := checkAPIKey(apiKey)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		c.Next()
	}
}

func checkAPIKey(apiKey string) error {
	if apiKey != os.Getenv("KAJILABSTORE_API_KEY") {
		err := errors.New("invalid API Key")
		log.Printf("%v", err)
		return err
	}
	return nil
}
