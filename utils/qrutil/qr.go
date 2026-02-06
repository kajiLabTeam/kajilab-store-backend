package qrutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func HMACSHA256Short(message string) string {
	secret := os.Getenv("HASH_KEY")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	full := hex.EncodeToString(mac.Sum(nil))
	return full[:16]
}
