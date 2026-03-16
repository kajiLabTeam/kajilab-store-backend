// Package middleware provides middleware handlers for authentication and other HTTP middleware features.
package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ヘッダー取得
		apiKey := c.GetHeader("X-API-Key") // クライアントからAPIキーを取得

		fmt.Println("apikey", apiKey)

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
	fmt.Println("apikey", os.Getenv("KAJILABSTORE_API_KEY"))
	fmt.Println(os.Getenv("HASH_KEY"))
	if apiKey != os.Getenv("KAJILABSTORE_API_KEY") {
		err := errors.New("invalid API Key")
		log.Printf("%v", err)
		return err
	}
	return nil
}
