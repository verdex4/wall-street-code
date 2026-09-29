package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/verdex4/wall-street-code/auth"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Not found .env file")
	}
	res := auth.Hello()
	fmt.Println(res)
}