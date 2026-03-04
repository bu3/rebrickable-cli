package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	legoSetCommands()
}

func legoSetCommands() {
	legoCmd.AddCommand(legoSetsCmd)
	legoSetsCmd.AddCommand(getLegoSetsCmd)
	legoSetsCmd.AddCommand(getLegoSetCmd)
	legoSetsCmd.AddCommand(getLegoSetAlternatesCmd)
	legoSetsCmd.AddCommand(getLegoSetMinifigsCmd)
	legoSetsCmd.AddCommand(getLegoSetPartsCmd)
	legoSetsCmd.AddCommand(getLegoSetSetsCmd)

	getLegoSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	getLegoSetAlternatesCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	getLegoSetMinifigsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	getLegoSetPartsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	getLegoSetSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
}

var legoSetsCmd = &cobra.Command{
	Use:   "sets",
	Short: "LEGO catalog set actions",
}

var getLegoSetsCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoSets()
		if err != nil {
			return err
		}
		output, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	},
}

var getLegoSetCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoSet(adjustedSetNumber())
		if err != nil {
			return err
		}
		output, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	},
}

var getLegoSetAlternatesCmd = &cobra.Command{
	Use:   "alternates",
	Short: "alternates",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoSetAlternates(adjustedSetNumber())
		if err != nil {
			return err
		}
		output, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	},
}

var getLegoSetMinifigsCmd = &cobra.Command{
	Use:   "minifigs",
	Short: "minifigs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoSetMinifigs(adjustedSetNumber())
		if err != nil {
			return err
		}
		output, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	},
}

var getLegoSetPartsCmd = &cobra.Command{
	Use:   "parts",
	Short: "parts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoSetParts(adjustedSetNumber())
		if err != nil {
			return err
		}
		output, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	},
}

var getLegoSetSetsCmd = &cobra.Command{
	Use:   "sets",
	Short: "sets",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoSetSets(adjustedSetNumber())
		if err != nil {
			return err
		}
		output, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	},
}
