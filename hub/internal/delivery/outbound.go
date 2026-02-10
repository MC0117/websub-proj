package delivery

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

// payload is signed and delivers the JSON to a sub
func SendPayload(callbackURL, secret string, message []byte) error {
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write(message)

	if err != nil {
		return err
	}
	signature := "sha256=" + hex.EncodeToString(h.Sum(nil))

	req, _ := http.NewRequest(http.MethodPost, callbackURL, bytes.NewBuffer(message))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature", signature)

	client := &http.Client{}
	_, err = client.Do(req)
	return err
}
