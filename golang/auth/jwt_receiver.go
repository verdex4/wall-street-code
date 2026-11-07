package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type JWTReceiver struct {
	jwt       string
	client    *http.Client
	sleepTime time.Duration
}

// inializes JWTReceiver with valid JWT
func NewJWTReceiver(client *http.Client) (*JWTReceiver, error) {
	r := &JWTReceiver{
		jwt:    "",
		client: client,
	}
	t, err := r.newJWT()
	if err != nil {
		return r, err
	}
	r.jwt = t
	return r, nil
}

func (r *JWTReceiver) GetJWT() string {
	return r.jwt
}

func (r *JWTReceiver) newJWT() (string, error) {
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		return "", errors.New("Env key `API_TOKEN` not found")
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		return "", errors.New("Env key `BASE_URL` not found")
	}

	body := fmt.Appendf(nil, `{"secret": "%s"}`, apiToken)
	resp, err := r.client.Post(
		fmt.Sprintf("%s/v1/sessions", baseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var respStruct struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return "", err
	}

	return respStruct.Token, nil
}

// Infinite function to updating JWT, must be ran using new goroutine
func (r *JWTReceiver) RefreshJWT(sleepTime time.Duration) error {
	for {
		t, err := r.newJWT()
		if err != nil {
			return err
		}
		r.jwt = t
		time.Sleep(sleepTime)
	}
}
