/*
Copyright © 2025 nu12
*/
package cmd

import (
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/nu12/az-dump/internal/create"
	"github.com/spf13/cobra"
)

var outputPath string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a backup by exporting the ARM templates",
	Long: `Create a backup by exporting the ARM templates.
	
Export all resource groups from a subscription:
az-dump create -s <subscription-id>

Export a specific resource group from a subscription:
az-dump create -s <subscription-id> -g <resource-group-name>

Export a list of resource groups from a subscription:
az-dump create -s <subscription-id> -g <resource-group-name>,<resource-group-name>[,<resource-group-name>]

Export all resource groups from a specific subscription and save the templates to a specific path:
az-dump create -s <subscription-id> -o <output-path>`,

	Run: func(cmd *cobra.Command, args []string) {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			panic(err)
		}
		client, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(outputPath, 0755)
		if err != nil {
			panic(err)
		}
		if rgName == "" {
			err := create.ExportAllTemplates(client, outputPath)
			if err != nil {
				panic(err)
			}
			return
		}
		err = create.ExportListOfTemplates(client, rgName, outputPath)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	createCmd.Flags().StringVarP(&subscriptionID, "subscription", "s", "", "Subscription ID to use (required)")
	createCmd.Flags().StringVarP(&rgName, "rg", "g", "", "Comma separated list of resource group names to dump (empty for all resource groups)")
	createCmd.Flags().StringVarP(&outputPath, "output", "o", time.Now().Format("20060102150405"), "Path where to save the templates")
	err := createCmd.MarkFlagRequired("subscription")
	if err != nil {
		panic(err)
	}
}
