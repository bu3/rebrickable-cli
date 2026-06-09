package cmd

import (
	"github.com/spf13/cobra"
)

var colorID string

func init() {
	legoColorsCommands()
}

func legoColorsCommands() {
	legoCmd.AddCommand(legoColorsCmd)
	legoColorsCmd.AddCommand(getLegoColorsCmd)
	legoColorsCmd.AddCommand(getLegoColorCmd)

	getLegoColorCmd.Flags().StringVarP(&colorID, "id", "i", "", "Color id")
	_ = getLegoColorCmd.MarkFlagRequired("id")
}

var legoColorsCmd = &cobra.Command{
	Use:   "colors",
	Short: "LEGO catalog color actions",
}

var getLegoColorsCmd = &cobra.Command{
	Use:   "list",
	Short: "list all colors",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetLegoColors()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoColorCmd = &cobra.Command{
	Use:   "get",
	Short: "get a color by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetLegoColor(colorID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
