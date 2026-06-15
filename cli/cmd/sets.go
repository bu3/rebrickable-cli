package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	rebrickable "github.com/bu3/rebrickable-go"
	"github.com/spf13/cobra"
)

var setNumber string
var setListName string
var setListID string
var quantity int
var includeSpares bool

func init() {
	setListsCommands()
	setCommands()
	setListSetsCommands()
}

func setCommands() {
	user.AddCommand(setsCmd)
	setsCmd.AddCommand(getSetsCmd)
	setsCmd.AddCommand(getSetCmd)
	setsCmd.AddCommand(saveSetsCmd)
	setsCmd.AddCommand(replaceSetCmd)
	setsCmd.AddCommand(deleteSetsCmd)
	setsCmd.AddCommand(syncSetsCmd)

	saveSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	deleteSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	getSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	replaceSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	replaceSetCmd.Flags().IntVarP(&quantity, "quantity", "q", 1, "Quantity")
	syncSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
}

func setListsCommands() {
	user.AddCommand(setListsCmd)
	setListsCmd.AddCommand(saveSetListCmd)
	setListsCmd.AddCommand(getSetListsCmd)
	setListsCmd.AddCommand(getSetListCmd)
	setListsCmd.AddCommand(updateSetListCmd)
	setListsCmd.AddCommand(replaceSetListCmd)
	setListsCmd.AddCommand(deleteSetListsCmd)

	saveSetListCmd.Flags().StringVarP(&setListName, "name", "n", "", "Set List name")
	deleteSetListsCmd.Flags().StringVarP(&setNumber, "set_list_num", "l", "", "Set List id")
	getSetListCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	updateSetListCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	updateSetListCmd.Flags().StringVarP(&setListName, "name", "n", "", "Set List name")
	replaceSetListCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	replaceSetListCmd.Flags().StringVarP(&setListName, "name", "n", "", "Set List name")
}

func setListSetsCommands() {
	user.AddCommand(setListSetsCmd)
	setListSetsCmd.AddCommand(getSetListSetsCmd)
	setListSetsCmd.AddCommand(getSetListSetCmd)
	setListSetsCmd.AddCommand(saveSetListSetCmd)
	setListSetsCmd.AddCommand(deleteSetListSetCmd)
	setListSetsCmd.AddCommand(updateSetListSetCmd)
	setListSetsCmd.AddCommand(replaceSetListSetCmd)

	getSetListSetsCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	getSetListSetCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	getSetListSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	saveSetListSetCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	saveSetListSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	deleteSetListSetCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	deleteSetListSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	updateSetListSetCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	updateSetListSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	updateSetListSetCmd.Flags().IntVarP(&quantity, "quantity", "q", 1, "Quantity")
	updateSetListSetCmd.Flags().BoolVar(&includeSpares, "include_spares", false, "Include spare parts")
	replaceSetListSetCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	replaceSetListSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	replaceSetListSetCmd.Flags().IntVarP(&quantity, "quantity", "q", 1, "Quantity")
	replaceSetListSetCmd.Flags().BoolVar(&includeSpares, "include_spares", false, "Include spare parts")
}

func newAPIClient(cmd *cobra.Command) *rebrickable.Client {
	return cmd.Context().Value(rebrickableClient).(*rebrickable.Client)
}

var setListsCmd = &cobra.Command{
	Use:   "setLists",
	Short: "setLists",
}

var saveSetListCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.StoreUserSetList(setListName); err != nil {
			return err
		}
		fmt.Println("SetList saved")
		return nil
	},
}

var getSetListsCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetUserSetLists()
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

var deleteSetListsCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.DeleteUserSetList(setNumber); err != nil {
			return err
		}
		fmt.Printf("Deleted set list: %s\n", setNumber)
		return nil
	},
}

var setsCmd = &cobra.Command{
	Use:   "sets",
	Short: "sets",
}

var getSetsCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		setsResponse, err := client.GetUserSets()
		if err != nil {
			return err
		}
		output, err := json.MarshalIndent(setsResponse, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		return nil
	},
}

var deleteSetsCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.DeleteUserSet(adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Printf("Deleted set: %s\n", adjustedSetNumber())
		return nil
	},
}

var saveSetsCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.StoreUserSet(adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Println("Set saved")
		return nil
	},
}

var getSetListCmd = &cobra.Command{
	Use:   "getOne",
	Short: "getOne",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetUserSetList(setListID)
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

var updateSetListCmd = &cobra.Command{
	Use:   "update",
	Short: "update",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.UpdateUserSetList(setListID, setListName); err != nil {
			return err
		}
		fmt.Printf("Updated set list: %s\n", setListID)
		return nil
	},
}

var replaceSetListCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.ReplaceUserSetList(setListID, setListName); err != nil {
			return err
		}
		fmt.Printf("Replaced set list: %s\n", setListID)
		return nil
	},
}

var setListSetsCmd = &cobra.Command{
	Use:   "setListSets",
	Short: "setListSets actions",
}

var getSetListSetsCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetUserSetListSets(setListID)
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

var getSetListSetCmd = &cobra.Command{
	Use:   "getOne",
	Short: "getOne",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetUserSetListSet(setListID, adjustedSetNumber())
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

var saveSetListSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.StoreUserSetListSet(setListID, adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Println("Set added to set list")
		return nil
	},
}

var deleteSetListSetCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.DeleteUserSetListSet(setListID, adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Printf("Deleted %s from set list %s\n", adjustedSetNumber(), setListID)
		return nil
	},
}

var getSetCmd = &cobra.Command{
	Use:   "getOne",
	Short: "getOne",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		result, err := client.GetUserSet(adjustedSetNumber())
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

var replaceSetCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.ReplaceUserSet(adjustedSetNumber(), quantity); err != nil {
			return err
		}
		fmt.Printf("Updated set: %s\n", adjustedSetNumber())
		return nil
	},
}

func adjustedSetNumber() string {
	if !strings.HasSuffix(setNumber, "-1") {
		return setNumber + "-1"
	}

	return setNumber
}

var syncSetsCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.SyncUserSet(adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Printf("Synced set: %s\n", adjustedSetNumber())
		return nil
	},
}

var updateSetListSetCmd = &cobra.Command{
	Use:   "update",
	Short: "update",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.UpdateUserSetListSet(setListID, adjustedSetNumber(), quantity, includeSpares); err != nil {
			return err
		}
		fmt.Println("Updated set in set list")
		return nil
	},
}

var replaceSetListSetCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.ReplaceUserSetListSet(setListID, adjustedSetNumber(), quantity, includeSpares); err != nil {
			return err
		}
		fmt.Println("Replaced set in set list")
		return nil
	},
}
