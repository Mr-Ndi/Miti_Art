package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type PaypackTokenResponse struct {
	Access string `json:"access"`
}

type PaymentPayload struct {
	Amount float64 `json:"amount"`
	Number string  `json:"number"`
}

func GetPaypackToken() (string, error) {
	body := map[string]string{
		"client_id":     os.Getenv("PAYPACK_CLIENT_ID"),
		"client_secret": os.Getenv("PAYPACK_CLIENT_SECRET"),
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "https://paypack.rw/api/auth/token", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("authentication failed")
	}

	var tokenRes PaypackTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return "", err
	}

	return tokenRes.Access, nil
}

func InitiatePayment(token string, payload PaymentPayload) error {
	jsonBody, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "https://paypack.rw/api/payments/collection", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("payment not accepted by Paypack")
	}

	return nil
}
