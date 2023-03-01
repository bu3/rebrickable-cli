package cmd

import (
	"fmt"
	"github.com/bu3/rebrickable-cli/cmd/api"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var setNumber string

func init() {
	user.AddCommand(setsCmd)
	setsCmd.AddCommand(getSetsCmd)
	setsCmd.AddCommand(saveSetsCmd)
	setsCmd.AddCommand(deleteSetsCmd)

	saveSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set numbers")
	deleteSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set numbers")
}

var setsCmd = &cobra.Command{
	Use:   "sets",
	Short: "sets",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var getSetsCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		GetUserSets(client, apiKey, authToken)
		return nil
	},
}

var deleteSetsCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		DeleteUserSet(client, apiKey, authToken)
		return nil
	},
}

var saveSetsCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		StoreUserSet(client, apiKey, authToken, setNumber)
		return nil
	},
}

func DeleteUserSet(client *resty.Client, apiKey string, authToken string) {
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		Delete(api.GetURL(fmt.Sprintf("/users/%s/sets/%s/", authToken, setNumber)))

	if resp.StatusCode() == 204 {
		fmt.Println("Deleted set: 10276-1")
	}
}

func StoreUserSet(client *resty.Client, apiKey string, authToken string, setNumber string) {
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetBody(fmt.Sprintf(`{"set_num":  "%s","quantity": "1"}`, setNumber)).
		Post(api.GetURL(fmt.Sprintf("/users/%s/sets/", authToken)))

	if resp.StatusCode() == 201 {
		fmt.Println("Sets saved: ", resp)
	}
}

func GetUserSets(client *resty.Client, apiKey string, authToken string) *SetsResponse {
	setsResponse := &SetsResponse{}
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetResult(setsResponse).
		Get(api.GetURL(fmt.Sprintf("/users/%s/sets", authToken)))

	if resp.StatusCode() == 200 {
		fmt.Println("Sets Found: ", resp)
	}
	return setsResponse
}

type SetsResponse struct {
	count   int `json:"count"`
	results []map[string]string
}
