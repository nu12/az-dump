package create

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

type FakeClient struct{}

func (c *FakeClient) BeginExportTemplate(context.Context, string, armresources.ExportTemplateRequest, *armresources.ResourceGroupsClientBeginExportTemplateOptions) (*runtime.Poller[armresources.ResourceGroupsClientExportTemplateResponse], error) {

	return runtime.NewPoller[armresources.ResourceGroupsClientExportTemplateResponse](&http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       io.NopCloser(bytes.NewBufferString("{\"template\":{\"$schema\":\"https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#\",\"contentVersion\":\"1.0.0.0\",\"parameters\":{},\"resources\":[],\"variables\":{}}}")),
		Request: &http.Request{
			Method: http.MethodPost,
		},
	}, runtime.NewPipeline("", "", runtime.PipelineOptions{}, &policy.ClientOptions{}), &runtime.NewPollerOptions[armresources.ResourceGroupsClientExportTemplateResponse]{})
}

func (w *FakeClient) NewListPager(options *armresources.ResourceGroupsClientListOptions) *runtime.Pager[armresources.ResourceGroupsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[armresources.ResourceGroupsClientListResponse]{
		More: func(armresources.ResourceGroupsClientListResponse) bool {
			return false
		},
		Fetcher: func(context.Context, *armresources.ResourceGroupsClientListResponse) (armresources.ResourceGroupsClientListResponse, error) {
			return armresources.ResourceGroupsClientListResponse{
				ResourceGroupListResult: armresources.ResourceGroupListResult{
					Value: []*armresources.ResourceGroup{
						{
							Name: to.Ptr("azure"),
						},
					},
					NextLink: nil,
				},
			}, nil
		},
	})
}

func TestExportTemplate(t *testing.T) {
	client := &FakeClient{}

	output := t.TempDir()
	err := exportOneTemplate(client, "azure", output)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if _, err := os.Stat(output + "/azure.json"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("expected file to exist, got %v", err)
		}
	}
}
