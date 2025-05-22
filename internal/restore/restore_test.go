package restore

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type FakeRGClient struct{}

func (c *FakeRGClient) Get(context.Context, string, *armresources.ResourceGroupsClientGetOptions) (armresources.ResourceGroupsClientGetResponse, error) {
	return armresources.ResourceGroupsClientGetResponse{
		ResourceGroup: armresources.ResourceGroup{
			Name:     to.Ptr("azure"),
			Location: to.Ptr("eastus"),
		},
	}, nil
}
func (c *FakeRGClient) CreateOrUpdate(context.Context, string, armresources.ResourceGroup, *armresources.ResourceGroupsClientCreateOrUpdateOptions) (armresources.ResourceGroupsClientCreateOrUpdateResponse, error) {
	return armresources.ResourceGroupsClientCreateOrUpdateResponse{
		ResourceGroup: armresources.ResourceGroup{
			Name:     to.Ptr("azure"),
			Location: to.Ptr("eastus"),
		},
	}, nil
}

type FakeDeploymentClient struct{}

func (c *FakeDeploymentClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters armresources.Deployment, options *armresources.DeploymentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[armresources.DeploymentsClientCreateOrUpdateResponse], error) {
	return runtime.NewPoller[armresources.DeploymentsClientCreateOrUpdateResponse](&http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(bytes.NewBufferString("{\"properties\":{\"provisioningState\":\"Succeeded\"}}")),
		Request: &http.Request{
			Method: http.MethodPost,
		},
	}, runtime.NewPipeline("", "", runtime.PipelineOptions{}, &policy.ClientOptions{}), &runtime.NewPollerOptions[armresources.DeploymentsClientCreateOrUpdateResponse]{})
}

func TestImportTemplate(t *testing.T) {

	dir := t.TempDir()
	tempFile, _ := os.Create(dir + "/azure.json")
	tempFile.WriteString("{\"$schema\": \"https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#\",\"contentVersion\": \"1.0.0.0\",\"parameters\": {},\"resources\": [],\"variables\": {}}")
	defer tempFile.Close()

	err := importOneTemplate(&FakeRGClient{}, &FakeDeploymentClient{}, dir, "azure.json", true, "eastus")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
