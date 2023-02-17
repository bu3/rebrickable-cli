package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

func main() {
	client := resty.New()
	apiKey, authToken := Login(client)

	//Get sets again
	GetUserSets(client, apiKey, authToken)

	//Add set
	StoreUserSet(client, apiKey, authToken)

	//Get sets again
	GetUserSets(client, apiKey, authToken)

	//Delete all sets
	DeleteUserSet(client, apiKey, authToken)

	//Get sets again
	GetUserSets(client, apiKey, authToken)
}

func Login(client *resty.Client) (string, *AuthToken) {
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
	return apiKey, authToken
}

func DeleteUserSet(client *resty.Client, apiKey string, authToken *AuthToken) {
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		Delete(fmt.Sprintf("https://rebrickable.com/api/v3/users/%s/sets/%s/", authToken.UserToken, "10276-1"))

	if resp.StatusCode() == 204 {
		fmt.Println("Deleted set: 10276-1")
	}
}

func StoreUserSet(client *resty.Client, apiKey string, authToken *AuthToken) {
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetBody(`{"set_num":  "10276-1","quantity": "1"}`).
		Post(fmt.Sprintf("https://rebrickable.com/api/v3/users/%s/sets/", authToken.UserToken))

	if resp.StatusCode() == 201 {
		fmt.Println("Sets saved: ", resp)
	}
}

func GetUserSets(client *resty.Client, apiKey string, authToken *AuthToken) *SetsResponse {
	setsResponse := &SetsResponse{}
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Get(fmt.Sprintf("https://rebrickable.com/api/v3/users/%s/sets", authToken.UserToken))

	if resp.StatusCode() == 200 {
		fmt.Println("Sets Found: ", resp)
	}
	return setsResponse
}

type AuthToken struct {
	UserToken string `json:"user_token"`
}

type SetsResponse struct {
	count   int `json:"count"`
	results []map[string]string
}
