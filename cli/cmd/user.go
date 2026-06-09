package cmd

import (
	"os"

	rebrickable "github.com/bu3/rebrickable-go"
	"github.com/spf13/cobra"
	"context"
)

type contextKey string

const rebrickableClient contextKey = "rebrickable_client"

func init() {
	rootCmd.AddCommand(user)
}

var user = &cobra.Command{
	Use:   "user",
	Short: "user actions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("REBRICKABLE_API_KEY")
		username := os.Getenv("REBRICKABLE_USERNAME")
		password := os.Getenv("REBRICKABLE_PASSWORD")
		client, err := rebrickable.NewAuthenticatedClient(apiKey, username, password)
		if err != nil {
			return err
		}
		cmd.SetContext(context.WithValue(cmd.Context(), rebrickableClient, client))
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
