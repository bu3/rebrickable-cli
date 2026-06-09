package cmd

import (
	"github.com/spf13/cobra"
)

var figNum string

func init() {
	legoMinifigsCommands()
}

func legoMinifigsCommands() {
	legoCmd.AddCommand(legoMinifigsCmd)
	legoMinifigsCmd.AddCommand(getLegoMinifigsCmd)
	legoMinifigsCmd.AddCommand(getLegoMinifigCmd)
	legoMinifigsCmd.AddCommand(getLegoMinifigPartsCmd)
	legoMinifigsCmd.AddCommand(getLegoMinifigSetsCmd)

	getLegoMinifigCmd.Flags().StringVarP(&figNum, "fig_num", "f", "", "Minifig number")
	_ = getLegoMinifigCmd.MarkFlagRequired("fig_num")

	getLegoMinifigPartsCmd.Flags().StringVarP(&figNum, "fig_num", "f", "", "Minifig number")
	_ = getLegoMinifigPartsCmd.MarkFlagRequired("fig_num")

	getLegoMinifigSetsCmd.Flags().StringVarP(&figNum, "fig_num", "f", "", "Minifig number")
	_ = getLegoMinifigSetsCmd.MarkFlagRequired("fig_num")
}

var legoMinifigsCmd = &cobra.Command{
	Use:   "minifigs",
	Short: "LEGO catalog minifig actions",
}

var getLegoMinifigsCmd = &cobra.Command{
	Use:   "list",
	Short: "list all minifigs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetLegoMinifigs()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoMinifigCmd = &cobra.Command{
	Use:   "get",
	Short: "get a minifig by fig_num",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetLegoMinifig(figNum)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoMinifigPartsCmd = &cobra.Command{
	Use:   "parts",
	Short: "list parts of a minifig",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetLegoMinifigParts(figNum)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoMinifigSetsCmd = &cobra.Command{
	Use:   "sets",
	Short: "list sets containing a minifig",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetLegoMinifigSets(figNum)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
