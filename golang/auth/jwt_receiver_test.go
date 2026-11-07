package auth

import (
	"net/http"
	"testing"
	"time"
)

func TestReceiver(t *testing.T) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	r, err := NewJWTReceiver(client)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	jwt := r.GetJWT()
	if jwt == "" {
		t.Errorf("jwt should be valid, got empty string")
	}

	sleepTime := 3 * time.Second
	go r.RefreshJWT(sleepTime)

	time.Sleep(1 * time.Second) // waiting for goroutine take JWT

	newJWT := r.GetJWT()

	if newJWT == "" {
		t.Errorf("new jwt should be valid, got empty string")
	}
	if newJWT == jwt {
		t.Errorf("jwt should be updated after sleeping, got same JWT: %s", newJWT)
	}

	time.Sleep(3 * time.Second)

	newJWT2 := r.GetJWT()

	if newJWT2 == "" {
		t.Errorf("new jwt should be valid, got empty string")
	}
	if newJWT2 == newJWT {
		t.Errorf("jwt should be updated after sleeping, got same JWT: %s", newJWT)
	}
}
