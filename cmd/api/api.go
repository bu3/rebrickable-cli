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

func StoreUserSetList(client *resty.Client, apiKey string, authToken string, name string) {
	url := GetURL(fmt.Sprintf("/users/%s/setlists/", authToken))
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetBody(fmt.Sprintf(`{"name":  "%s"}`, name)).
		Post(url)

	if resp.StatusCode() == 201 {
		fmt.Println("SetList saved: ", resp)
	}
}

func GetUserSetLists(client *resty.Client, apiKey string, authToken string) *SetsResponse {
	setsResponse := &SetsResponse{}
	url := GetURL(fmt.Sprintf("/users/%s/setlists", authToken))
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Get(url)

	if resp.StatusCode() == 200 {
		fmt.Println("SetLists Found: ", resp)
	}
	return setsResponse
}

func DeleteUserSetList(client *resty.Client, apiKey string, authToken string, id string) {
	url := GetURL(fmt.Sprintf("/users/%s/setlists/%s/", authToken, id))
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		Delete(url)

	if resp.StatusCode() == 204 {
		fmt.Println(fmt.Sprintf("Deleted set: %s", id))
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

func GetUserSets(client *resty.Client, apiKey string, authToken string) (*SetsResponse, error) {
	setsResponse := &SetsResponse{}
	url := GetURL(fmt.Sprintf("/users/%s/sets", authToken))
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Get(url)

	if err != nil || resp.StatusCode() != 200 {
		fmt.Println(fmt.Sprintf("Something went wrong: Status code: [%d] - Error: %v", resp.StatusCode(), err))
		return nil, err
	}
	return setsResponse, nil
}

func DeleteUserSet(client *resty.Client, apiKey string, authToken string, setNumber string) {
	url := GetURL(fmt.Sprintf("/users/%s/sets/%s/", authToken, setNumber))
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		Delete(url)

	if resp.StatusCode() == 204 {
		fmt.Println(fmt.Sprintf("Deleted set: %s", setNumber))
	}
}

type SetsResponse struct {
	Count   int              `json:"count"`
	Results []map[string]any `json:"results"`
}
