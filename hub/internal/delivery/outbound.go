package delivery

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"
)

// payload is signed and delivers the JSON to a sub
func SendPayload(callbackURL, secret string, message []byte) error {
	//Client provided secret is used to generate the signature
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write(message)

	if err != nil {
		return err
	}

	// signature formatted as per documentation
	signature := "sha256=" + hex.EncodeToString(h.Sum(nil))
	req, err := http.NewRequest(http.MethodPost, callbackURL, bytes.NewBuffer(message))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature", signature)

	client := &http.Client{Timeout: time.Second * 15} //timeout set to 15 seconds
	_, err = client.Do(req)
	return err
}
