package cmd

import (
	"encoding/json"
	"fmt"

	rebrickable "github.com/bu3/rebrickable-go"
	"github.com/spf13/cobra"
)

var (
	partNumber  string
	partColorID string
	partsFilter rebrickable.PartsFilter
)

func init() {
	legoPartsCommands()
}

func legoPartsCommands() {
	legoCmd.AddCommand(legoPartsCmd)
	legoPartsCmd.AddCommand(getLegoPartsCmd)
	legoPartsCmd.AddCommand(getLegoPartCmd)
	legoPartsCmd.AddCommand(getLegoPartColorsCmd)
	legoPartsCmd.AddCommand(getLegoPartColorCmd)
	legoPartsCmd.AddCommand(getLegoPartColorSetsCmd)

	getLegoPartsCmd.Flags().StringVar(&partsFilter.PartNum, "part_num", "", "filter by part number")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.PartNums, "part_nums", "", "comma-separated part numbers")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.PartCatID, "part_cat_id", "", "filter by part category id")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.ColorID, "color_id", "", "filter by color id")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.BricklinkID, "bricklink_id", "", "filter by BrickLink id")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.BrickowlID, "brickowl_id", "", "filter by BrickOwl id")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.LegoID, "lego_id", "", "filter by LEGO id")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.LdrawID, "ldraw_id", "", "filter by LDraw id")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.Ordering, "ordering", "", "ordering field")
	getLegoPartsCmd.Flags().StringVar(&partsFilter.Search, "search", "", "search term")

	getLegoPartCmd.Flags().StringVarP(&partNumber, "part_num", "n", "", "Part number")
	_ = getLegoPartCmd.MarkFlagRequired("part_num")

	getLegoPartColorsCmd.Flags().StringVarP(&partNumber, "part_num", "n", "", "Part number")
	_ = getLegoPartColorsCmd.MarkFlagRequired("part_num")

	getLegoPartColorCmd.Flags().StringVarP(&partNumber, "part_num", "n", "", "Part number")
	getLegoPartColorCmd.Flags().StringVarP(&partColorID, "color_id", "c", "", "Color id")
	_ = getLegoPartColorCmd.MarkFlagRequired("part_num")
	_ = getLegoPartColorCmd.MarkFlagRequired("color_id")

	getLegoPartColorSetsCmd.Flags().StringVarP(&partNumber, "part_num", "n", "", "Part number")
	getLegoPartColorSetsCmd.Flags().StringVarP(&partColorID, "color_id", "c", "", "Color id")
	_ = getLegoPartColorSetsCmd.MarkFlagRequired("part_num")
	_ = getLegoPartColorSetsCmd.MarkFlagRequired("color_id")
}

var legoPartsCmd = &cobra.Command{
	Use:   "parts",
	Short: "LEGO catalog part actions",
}

var getLegoPartsCmd = &cobra.Command{
	Use:   "list",
	Short: "list parts (supports filters)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoParts(partsFilter)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoPartCmd = &cobra.Command{
	Use:   "get",
	Short: "get a part by part_num",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPart(partNumber)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoPartColorsCmd = &cobra.Command{
	Use:   "colors",
	Short: "list colors a part has appeared in",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPartColors(partNumber)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoPartColorCmd = &cobra.Command{
	Use:   "colorDetail",
	Short: "get a specific part/color combination",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPartColor(partNumber, partColorID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoPartColorSetsCmd = &cobra.Command{
	Use:   "colorSets",
	Short: "list sets containing a part/color combination",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPartColorSets(partNumber, partColorID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

func printJSON(v any) error {
	output, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
