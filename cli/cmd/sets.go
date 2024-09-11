package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/bu3/rebrickable-cli/cli/cmd/api"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"strings"
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

var setListsCmd = &cobra.Command{
	Use:   "setLists",
	Short: "setLists",
}

var saveSetListCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		api.StoreUserSetList(client, apiKey, authToken, setListName)
		return nil
	},
}

var getSetListsCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		api.GetUserSetLists(client, apiKey, authToken)
		return nil
	},
}

var deleteSetListsCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		api.DeleteUserSetList(client, apiKey, authToken, setNumber)
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
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		//TODO: Add error handling
		//TODO: Move Json output to a dedicated class/function/whatever
		setsResponse, _ := api.GetUserSets(client, apiKey, authToken)
		output, _ := json.MarshalIndent(setsResponse, "", "\t")
		fmt.Println(string(output))
		return nil
	},
}

var deleteSetsCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		api.DeleteUserSet(client, apiKey, authToken, adjustedSetNumber())
		return nil
	},
}

var saveSetsCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := resty.New()
		authToken := cmd.Context().Value(AuthToken).(string)
		apiKey := cmd.Context().Value(ApiKey).(string)
		api.StoreUserSet(client, apiKey, authToken, adjustedSetNumber())
		return nil
	},
}

func adjustedSetNumber() string {
	if !strings.HasSuffix(setNumber, "-1") {
		return setNumber + "-1"
	}

	return setNumber
}
