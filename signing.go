package coprime

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func generateSig(message, secret string) (string, error) {

	signature := hmac.New(sha256.New, []byte(secret))
	_, err := signature.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}
