package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/verdex4/wall-street-code/utils"
)

var jwt string = ""

func init() {
	go refreshJWT()
}

func GetJWT() string {
	for jwt == "" {
		// ждем, пока горутина из init() не заполнит jwt
		time.Sleep(50 * time.Millisecond)
	}
	return jwt
}

// Обновляет JWT каждые 14,5 минут (действие одного - 15 минут согласно API)
func refreshJWT() {
	wg := sync.WaitGroup{}
	wg.Go(
		func() {
		defer wg.Done()
		for {
			jwt = newJWT()
			time.Sleep(14 * time.Minute + 30 * time.Second)
		}
	},
	)
	wg.Wait()
}

func newJWT() string {
	body := fmt.Appendf(nil, `{"secret": "%s"}`, utils.GetEnv("API_TOKEN"))

	resp, err := utils.HttpClient.Post(
		fmt.Sprintf("%s/v1/sessions", utils.GetEnv("BASE_URL")),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var respStruct struct {
		Token string `json:"token"`
	}

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		panic(err)
	}

	return respStruct.Token
}


