/*
Copyright © 2025 nu12
*/
package cmd

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/nu12/az-dump/internal/restore"
	"github.com/spf13/cobra"
)

var inputPath string
var allowCreate bool
var location string

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a backup from ARM templates",
	Long: `Restore a backup from ARM templates
	
Restore all resource groups from a subscription:
az-dump restore -s <subscription-id> -i <input-path>

Restore a specific resource group from a subscription:
az-dump restore -s <subscription-id> -g <resource-group-name> -i <input-path>

Restore a list of resource groups from a subscription:
az-dump restore -s <subscription-id> -g <resource-group-name>,<resource-group-name>[,<resource-group-name>] -i <input-path>

Restore all resource groups allowing the creation of the missing ones:
az-dump restore -s <subscription-id> -g <resource-group-name> -i <input-path> --create --location <location>`,

	Run: func(cmd *cobra.Command, args []string) {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			panic(err)
		}
		client, err := armresources.NewDeploymentsClient(subscriptionID, cred, nil)
		if err != nil {
			panic(err)
		}
		rgClient, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
		if err != nil {
			panic(err)
		}

		if rgName == "" {
			if err := restore.ImportAllTemplates(rgClient, client, inputPath, allowCreate, location); err != nil {
				panic(err)
			}
			return
		}

		if err := restore.ImportListOfTemplates(rgClient, client, rgName, inputPath, allowCreate, location); err != nil {
			panic(err)
		}
	},
}

func init() {
	restoreCmd.Flags().StringVarP(&subscriptionID, "subscription", "s", "", "Subscription ID to use (required)")
	restoreCmd.Flags().StringVarP(&rgName, "rg", "g", "", "Comma separated list of resource group names to restore (empty for all resource groups)")
	restoreCmd.Flags().StringVarP(&inputPath, "input", "i", "", "Path where the templates are stored (required)")
	restoreCmd.Flags().BoolVarP(&allowCreate, "create", "c", false, "Create the resource group if it doesn't exist")
	restoreCmd.Flags().StringVarP(&location, "location", "l", "", "Location of the resource group to be created (required if --create is set)")
	err := createCmd.MarkFlagRequired("subscription")
	if err != nil {
		panic(err)
	}
	err = restoreCmd.MarkFlagRequired("input")
	if err != nil {
		panic(err)
	}
	restoreCmd.MarkFlagsRequiredTogether("create", "location")
}
