package cmd

import (
	"fmt"
	"github.com/bu3/rebrickable-cli/cmd/api"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
)

func init() {
	rootCmd.AddCommand(user)
}

const (
	ApiKey    = "api_key"
	AuthToken = "auth_token"
)

var user = &cobra.Command{
	Use:   "user",
	Short: "user actions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken, err := login(client)

		if err != nil {
			return err
		}

		ctx := context.WithValue(cmd.Context(), AuthToken, authToken.UserToken)
		ctx = context.WithValue(ctx, ApiKey, authToken.ApiKey)
		cmd.SetContext(ctx)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func login(client *resty.Client) (*authToken, error) {
	username := os.Getenv("REBRICKABLE_USERNAME")
	password := os.Getenv("REBRICKABLE_PASSWORD")
	apiKey := os.Getenv("REBRICKABLE_API_KEY")
	authToken := &authToken{
		ApiKey: apiKey,
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
		SetFormData(map[string]string{
			"username": username,
			"password": password,
		}).
		SetResult(authToken).
		Post(api.GetURL("/users/_token/"))

	if resp.StatusCode() != 200 || err != nil {
		fmt.Println("Login was not successful")
	}
	return authToken, nil
}

type authToken struct {
	UserToken string `json:"user_token"`
	ApiKey    string
}
