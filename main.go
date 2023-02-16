package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

func main() {
	client := resty.New()
	username := os.Getenv("REBRICKABLE_USERNAME")
	password := os.Getenv("REBRICKABLE_PASSWORD")
	apiKey := os.Getenv("REBRICKABLE_API_KEY")
	authToken := &AuthToken{}

	resp, _ := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetFormData(map[string]string{
			"username": username,
			"password": password,
		}).
		SetResult(authToken).
		Post("https://rebrickable.com/api/v3/users/_token/")

	if resp.StatusCode() == 200 {
		fmt.Println("Login successful!!!")
	}
}

type AuthToken struct {
	UserToken string `json:"user_token"`
}
