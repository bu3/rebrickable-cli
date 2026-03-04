package cmd

import (
	"os"

	"github.com/bu3/rebrickable-cli/cli/cmd/api"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

const LegoApiKey = "lego_api_key"

func init() {
	rootCmd.AddCommand(legoCmd)
}

var legoCmd = &cobra.Command{
	Use:   "lego",
	Short: "LEGO catalog actions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("REBRICKABLE_API_KEY")
		ctx := context.WithValue(cmd.Context(), LegoApiKey, apiKey)
		cmd.SetContext(ctx)
		return nil
	},
}

func newLegoAPIClient(cmd *cobra.Command) *api.Client {
	apiKey := cmd.Context().Value(LegoApiKey).(string)
	return api.NewLegoClient(apiKey)
}
