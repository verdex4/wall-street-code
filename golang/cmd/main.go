package main

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/verdex4/wall-street-code/auth"
)

func testJWT() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			apiToken := auth.GetJWT()
			curTime := time.Now().Format(time.TimeOnly)
			fmt.Printf("API token: %s, time: %s\n", apiToken, curTime)
			time.Sleep(time.Minute)
		}
	}()
	wg.Wait()
}

func main() {
	fmt.Println("Hello, World!")
}