package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(user)
}

var user = &cobra.Command{
	Use:   "user",
	Short: "user actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		//if err := someFunc(); err != nil {
		//	return err
		//}
		return nil
	},
}
