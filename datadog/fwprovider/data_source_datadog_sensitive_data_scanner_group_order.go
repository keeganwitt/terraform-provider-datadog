package fwprovider

import (
	"context"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-datadog/datadog/internal/utils"
)

var _ datasource.DataSourceWithConfigure = &SensitiveDataScannerGroupOrderDatasource{}

func NewSensitiveDataScannerGroupOrderDatasource() datasource.DataSource {
	return &SensitiveDataScannerGroupOrderDatasource{}
}

type SensitiveDataScannerGroupOrderDatasourceModel struct {
	ID       types.String `tfsdk:"id"`
	GroupIDs types.List   `tfsdk:"group_ids"`
}

type SensitiveDataScannerGroupOrderDatasource struct {
	Api  *datadogV2.SensitiveDataScannerApi
	Auth context.Context
}

func (d *SensitiveDataScannerGroupOrderDatasource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	providerData, ok := request.ProviderData.(*FrameworkProvider)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"")
		return
	}

	d.Api = providerData.DatadogApiInstances.GetSensitiveDataScannerApiV2()
	d.Auth = providerData.Auth
}

func (d *SensitiveDataScannerGroupOrderDatasource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "sensitive_data_scanner_group_order"
}

func (d *SensitiveDataScannerGroupOrderDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Provides a Datadog Sensitive Data Scanner Group Order API data source. This can be used to retrieve the order of Datadog Sensitive Data Scanner Groups.",
		Attributes: map[string]schema.Attribute{
			"group_ids": schema.ListAttribute{
				Description: "The list of Sensitive Data Scanner group IDs, in order. Logs are tested against the query filter of each index one by one following the order of the list.",
				ElementType: types.StringType,
				Computed:    true,
			},
			// Resource ID
			"id": utils.ResourceIDAttribute(),
		},
	}
}

func (d *SensitiveDataScannerGroupOrderDatasource) Read(ctx context.Context, _ datasource.ReadRequest, response *datasource.ReadResponse) {
	var state SensitiveDataScannerGroupOrderDatasourceModel
	if response.Diagnostics.HasError() {
		return
	}

	resp, httpResponse, err := d.Api.ListScanningGroups(d.Auth)
	if err != nil {
		response.Diagnostics.Append(utils.FrameworkErrorDiag(err, fmt.Sprintf("error reading SDS groups. http response: %v", httpResponse)))
		return
	}
	var groups []datadogV2.SensitiveDataScannerGroupItem
	var groupID string
	if respData, ok := resp.GetDataOk(); ok {
		if respRelationships, ok := respData.GetRelationshipsOk(); ok {
			if respGroups, ok := respRelationships.GetGroupsOk(); ok {
				groups = respGroups.GetData()
			}
		}
		groupID = respData.GetId()
	}
	tfList := make([]string, len(groups))
	for i, ddGroup := range groups {
		tfList[i] = ddGroup.GetId()
	}

	state.GroupIDs, _ = types.ListValueFrom(ctx, types.StringType, tfList)
	state.ID = types.StringValue(groupID)
	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}
