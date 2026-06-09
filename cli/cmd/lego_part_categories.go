package cmd

import (
	"github.com/spf13/cobra"
)

var partCategoryID string

func init() {
	legoPartCategoriesCommands()
}

func legoPartCategoriesCommands() {
	legoCmd.AddCommand(legoPartCategoriesCmd)
	legoPartCategoriesCmd.AddCommand(getLegoPartCategoriesCmd)
	legoPartCategoriesCmd.AddCommand(getLegoPartCategoryCmd)

	getLegoPartCategoryCmd.Flags().StringVarP(&partCategoryID, "id", "i", "", "Part category id")
	_ = getLegoPartCategoryCmd.MarkFlagRequired("id")
}

var legoPartCategoriesCmd = &cobra.Command{
	Use:   "part-categories",
	Short: "LEGO catalog part category actions",
}

var getLegoPartCategoriesCmd = &cobra.Command{
	Use:   "list",
	Short: "list all part categories",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetLegoPartCategories()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoPartCategoryCmd = &cobra.Command{
	Use:   "get",
	Short: "get a part category by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetLegoPartCategory(partCategoryID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
