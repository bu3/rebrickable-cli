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
		fmt.Println("Login successful!!")
	}

	// Find existing sets
	setsResponse := &SetsResponse{}
	resp, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Get(fmt.Sprintf("https://rebrickable.com/api/v3/users/%s/sets", authToken.UserToken))

	if resp.StatusCode() == 200 {
		fmt.Println("Sets Found: ", resp)
		fmt.Println(resp)
	}

	//Add set
	resp, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetBody(`{"set_num":  "10276-1","quantity": "1"}`).
		Post(fmt.Sprintf("https://rebrickable.com/api/v3/users/%s/sets/", authToken.UserToken))

	if resp.StatusCode() != 201 {
		fmt.Println(resp.StatusCode())
		panic("Save failed")
	}

	//Get sets again
	setsResponse = &SetsResponse{}
	resp, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Get(fmt.Sprintf("https://rebrickable.com/api/v3/users/%s/sets", authToken.UserToken))

	if resp.StatusCode() == 200 {
		fmt.Println("Sets Found: ", resp)
		fmt.Println(resp)
	}

	//Delete all sets
	resp, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Delete(fmt.Sprintf("https://rebrickable.com/api/v3/users/%s/sets/%s/", authToken.UserToken, "10276-1"))

	if resp.StatusCode() == 200 {
		fmt.Println("Deleted set: 10276-1")
	}

	//Get sets again
	setsResponse = &SetsResponse{}
	resp, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Get(fmt.Sprintf("https://rebrickable.com/api/v3/users/%s/sets", authToken.UserToken))

	if resp.StatusCode() == 200 {
		fmt.Println("Sets Found: ", resp)
		fmt.Println(resp)
	}

}

type AuthToken struct {
	UserToken string `json:"user_token"`
}

type SetsResponse struct {
	count   int `json:"count"`
	results []map[string]string
}
