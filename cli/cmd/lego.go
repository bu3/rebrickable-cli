package cmd

import (
	"os"

	rebrickable "github.com/bu3/rebrickable-go"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	rootCmd.AddCommand(legoCmd)
}

var legoCmd = &cobra.Command{
	Use:   "lego",
	Short: "LEGO catalog actions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("REBRICKABLE_API_KEY")
		client := rebrickable.NewClient(apiKey)
		cmd.SetContext(context.WithValue(cmd.Context(), rebrickableClient, client))
		return nil
	},
}

func newLegoAPIClient(cmd *cobra.Command) *rebrickable.Client {
	return cmd.Context().Value(rebrickableClient).(*rebrickable.Client)
}
