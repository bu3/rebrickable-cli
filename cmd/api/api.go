package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
)

const apiBaseURI = "https://rebrickable.com/api/v3"

func GetURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return fmt.Sprintf(apiBaseURI + path)
}

func DeleteUserSet(client *resty.Client, apiKey string, authToken string, setNumber string) {
	url := GetURL(fmt.Sprintf("/users/%s/sets/%s/", authToken, setNumber))
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		Delete(url)

	if resp.StatusCode() == 204 {
		fmt.Println("Deleted set: 10276-1")
	}
}

func StoreUserSet(client *resty.Client, apiKey string, authToken string, setNumber string) {
	url := GetURL(fmt.Sprintf("/users/%s/sets/", authToken))
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetBody(fmt.Sprintf(`{"set_num":  "%s","quantity": "1"}`, setNumber)).
		Post(url)

	if resp.StatusCode() == 201 {
		fmt.Println("Sets saved: ", resp)
	}
}

func GetUserSets(client *resty.Client, apiKey string, authToken string) *SetsResponse {
	setsResponse := &SetsResponse{}
	url := GetURL(fmt.Sprintf("/users/%s/sets", authToken))
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Get(url)

	if resp.StatusCode() == 200 {
		fmt.Println("Sets Found: ", resp)
	}
	return setsResponse
}

type SetsResponse struct {
	count   int `json:"count"`
	results []map[string]string
}
