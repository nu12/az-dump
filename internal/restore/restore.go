package restore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/nu12/az-dump/internal/helpers"
)

type ResourceGroupsGetter interface {
	Get(context.Context, string, *armresources.ResourceGroupsClientGetOptions) (armresources.ResourceGroupsClientGetResponse, error)
	CreateOrUpdate(context.Context, string, armresources.ResourceGroup, *armresources.ResourceGroupsClientCreateOrUpdateOptions) (armresources.ResourceGroupsClientCreateOrUpdateResponse, error)
}

type DeploymentCreater interface {
	BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters armresources.Deployment, options *armresources.DeploymentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[armresources.DeploymentsClientCreateOrUpdateResponse], error)
}

func ImportAllTemplates(rgClient ResourceGroupsGetter, deploymentClient DeploymentCreater, inputPath string, allowCreate bool, location string) error {
	files, err := os.ReadDir(inputPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		err = importOneTemplate(rgClient, deploymentClient, inputPath, file.Name(), allowCreate, location)
		if err != nil {
			return err
		}

	}
	return nil
}

func ImportListOfTemplates(rgClient ResourceGroupsGetter, deploymentClient DeploymentCreater, commaList string, inputPath string, allowCreate bool, location string) error {
	files, err := os.ReadDir(inputPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !helpers.ComaListContains(commaList, helpers.GetResourceGroupNameFromFileName(file.Name())) {
			continue
		}
		err = importOneTemplate(rgClient, deploymentClient, inputPath, file.Name(), allowCreate, location)
		if err != nil {
			return err
		}

	}
	return nil

}

func importOneTemplate(rgClient ResourceGroupsGetter, deploymentClient DeploymentCreater, inputPath, fileName string, allowCreate bool, location string) error {
	b, err := os.ReadFile(inputPath + "/" + fileName)
	if err != nil {
		return err
	}
	// Create RG if it doesn't exist
	rgName := helpers.GetResourceGroupNameFromFileName(fileName)
	_, err = rgClient.Get(context.Background(), rgName, nil)
	if err != nil && allowCreate {
		fmt.Println(fmt.Sprintf("Creating resource group %s in %s", rgName, location))
		_, err = rgClient.CreateOrUpdate(context.Background(), rgName, armresources.ResourceGroup{
			Location: to.Ptr(location),
		}, nil)
		if err != nil {
			return err
		}
	}
	if err != nil && !allowCreate {
		fmt.Printf("Resource group %s doesn't exist and creation flag is not present. Skipping...\n", rgName)
		return nil
	}

	// Deploy template
	var template any
	err = json.Unmarshal(b, &template)
	if err != nil {
		return err
	}

	fmt.Printf("Importing template for %s\n", rgName)
	ctx := context.Background()
	poller, err := deploymentClient.BeginCreateOrUpdate(ctx, rgName, "az-dump-restore", armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
			Parameters: map[string]interface{}{},
			Template:   template,
		},
	}, nil)
	if err != nil {
		return err
	}

	if _, err := poller.PollUntilDone(ctx, nil); err != nil {
		return err
	}

	return nil
}
