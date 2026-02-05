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

func init() {
	setListsCommands()
	setCommands()
}

func setCommands() {
	user.AddCommand(setsCmd)
	setsCmd.AddCommand(getSetsCmd)
	setsCmd.AddCommand(saveSetsCmd)
	setsCmd.AddCommand(deleteSetsCmd)

	saveSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set numbers")
	deleteSetsCmd.Flags().StringVarP(&setNumber, "set_num", "n", "", "Set numbers")
}

func setListsCommands() {
	user.AddCommand(setListsCmd)
	setListsCmd.AddCommand(saveSetListCmd)
	setListsCmd.AddCommand(getSetListsCmd)
	setListsCmd.AddCommand(deleteSetListsCmd)

	saveSetListCmd.Flags().StringVarP(&setListName, "name", "n", "", "Set List name")
	deleteSetListsCmd.Flags().StringVarP(&setNumber, "set_list_num", "l", "", "Set List id")
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

func adjustedSetNumber() string {
	if !strings.HasSuffix(setNumber, "-1") {
		return setNumber + "-1"
	}

	return setNumber
}
