package create

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Exporter interface {
	BeginExportTemplate(context.Context, string, armresources.ExportTemplateRequest, *armresources.ResourceGroupsClientBeginExportTemplateOptions) (*runtime.Poller[armresources.ResourceGroupsClientExportTemplateResponse], error)
	NewListPager(options *armresources.ResourceGroupsClientListOptions) *runtime.Pager[armresources.ResourceGroupsClientListResponse]
}

func ExportAllTemplates(client Exporter, outputPath string) error {
	ctx := context.Background()
	pager := client.NewListPager(nil)
	for {
		if !pager.More() {
			break
		}
		page, _ := pager.NextPage(ctx)
		for _, rg := range page.Value {
			if err := exportOneTemplate(client, *rg.Name, outputPath); err != nil {
				return err
			}
		}

	}
	return nil
}
func ExportListOfTemplates(client Exporter, commaList string, outputPath string) error {
	for _, rg := range strings.Split(commaList, ",") {
		if err := exportOneTemplate(client, rg, outputPath); err != nil {
			return err
		}
	}
	return nil
}

func exportOneTemplate(client Exporter, rg string, outputPath string) error {
	fmt.Printf("Exporting template for resource group %s\n", rg)
	ctx := context.Background()
	poller, err := client.BeginExportTemplate(ctx, rg, armresources.ExportTemplateRequest{
		Options: to.Ptr("SkipAllParameterization"),
		Resources: []*string{
			to.Ptr("*")},
	}, nil)
	if err != nil {
		return err
	}

	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(res.Template, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath+"/"+rg+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}
