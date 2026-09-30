package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var HttpClient = &http.Client{Timeout: 5 * time.Second}

func GetEnv(key string) string {
	res := os.Getenv(key)
	if res == "" {
		panic(fmt.Sprintf("Env key `%s` not found", key))
	}
	return res
}