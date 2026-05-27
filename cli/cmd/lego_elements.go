package cmd

import (
	"github.com/spf13/cobra"
)

var elementID string

func init() {
	legoElementsCommands()
}

func legoElementsCommands() {
	legoCmd.AddCommand(legoElementsCmd)
	legoElementsCmd.AddCommand(getLegoElementCmd)

	getLegoElementCmd.Flags().StringVarP(&elementID, "id", "i", "", "Element id")
	_ = getLegoElementCmd.MarkFlagRequired("id")
}

var legoElementsCmd = &cobra.Command{
	Use:   "elements",
	Short: "LEGO catalog element actions",
}

var getLegoElementCmd = &cobra.Command{
	Use:   "get",
	Short: "get an element by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoElement(elementID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
