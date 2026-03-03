package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bu3/rebrickable-cli/cli/cmd/api"
	"github.com/spf13/cobra"
)

var setNumber string
var setListName string
var setListID string
var quantity int

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

	saveSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	deleteSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	getSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	replaceSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	replaceSetCmd.Flags().IntVarP(&quantity, "quantity", "q", 1, "Quantity")
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

	getSetListSetsCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	getSetListSetCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	getSetListSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	saveSetListSetCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	saveSetListSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
	deleteSetListSetCmd.Flags().StringVarP(&setListID, "set_list_id", "l", "", "Set List id")
	deleteSetListSetCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set number")
}

func newAPIClient(cmd *cobra.Command) *api.Client {
	authToken := cmd.Context().Value(AuthToken).(string)
	apiKey := cmd.Context().Value(ApiKey).(string)
	return api.NewClient(apiKey, authToken)
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
		return client.StoreUserSetList(setListName)
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
		return client.DeleteUserSetList(setNumber)
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
		return client.DeleteUserSet(adjustedSetNumber())
	},
}

var saveSetsCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		return client.StoreUserSet(adjustedSetNumber())
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
		return client.UpdateUserSetList(setListID, setListName)
	},
}

var replaceSetListCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		return client.ReplaceUserSetList(setListID, setListName)
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
		return client.StoreUserSetListSet(setListID, adjustedSetNumber())
	},
}

var deleteSetListSetCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		return client.DeleteUserSetListSet(setListID, adjustedSetNumber())
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
		return client.ReplaceUserSet(adjustedSetNumber(), quantity)
	},
}

func adjustedSetNumber() string {
	if !strings.HasSuffix(setNumber, "-1") {
		return setNumber + "-1"
	}

	return setNumber
}
