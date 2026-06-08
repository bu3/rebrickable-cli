package cmd

import (
	"github.com/spf13/cobra"
)

var themeID string

func init() {
	legoThemesCommands()
}

func legoThemesCommands() {
	legoCmd.AddCommand(legoThemesCmd)
	legoThemesCmd.AddCommand(getLegoThemesCmd)
	legoThemesCmd.AddCommand(getLegoThemeCmd)

	getLegoThemeCmd.Flags().StringVarP(&themeID, "id", "i", "", "Theme id")
	_ = getLegoThemeCmd.MarkFlagRequired("id")
}

var legoThemesCmd = &cobra.Command{
	Use:   "themes",
	Short: "LEGO catalog theme actions",
}

var getLegoThemesCmd = &cobra.Command{
	Use:   "list",
	Short: "list all themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoThemes()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoThemeCmd = &cobra.Command{
	Use:   "get",
	Short: "get a theme by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoTheme(themeID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
