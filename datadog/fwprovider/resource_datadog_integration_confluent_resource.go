package fwprovider

import (
	"context"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-datadog/datadog/internal/utils"
	"github.com/terraform-providers/terraform-provider-datadog/datadog/internal/validators"
)

var (
	_ resource.ResourceWithConfigure   = &IntegrationConfluentResourceResource{}
	_ resource.ResourceWithImportState = &IntegrationConfluentResourceResource{}
)

type IntegrationConfluentResourceResource struct {
	Api  *datadogV2.ConfluentCloudApi
	Auth context.Context
}

type IntegrationConfluentResourceModel struct {
	ID           types.String `tfsdk:"id"`
	AccountId    types.String `tfsdk:"account_id"`
	ResourceId   types.String `tfsdk:"resource_id"`
	ResourceType types.String `tfsdk:"resource_type"`
	Tags         types.Set    `tfsdk:"tags"`
}

func NewIntegrationConfluentResourceResource() resource.Resource {
	return &IntegrationConfluentResourceResource{}
}

func (r *IntegrationConfluentResourceResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	providerData := request.ProviderData.(*FrameworkProvider)
	r.Api = providerData.DatadogApiInstances.GetConfluentCloudApiV2()
	r.Auth = providerData.Auth
}

func (r *IntegrationConfluentResourceResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "integration_confluent_resource"
}

func (r *IntegrationConfluentResourceResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Provides a Datadog IntegrationConfluentResource resource. This can be used to create and manage Datadog integration_confluent_resource.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Required:    true,
				Description: "Confluent Account ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_id": schema.StringAttribute{
				Description: "The ID associated with a Confluent resource.",
				Required:    true,
			},
			"resource_type": schema.StringAttribute{
				Optional:    true,
				Description: "The resource type of the Resource. Can be `kafka`, `connector`, `ksql`, or `schema_registry`.",
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				Description: "A list of strings representing tags. Can be a single key, or key-value pairs separated by a colon.",
				ElementType: types.StringType,
				Validators:  []validator.Set{validators.TagsSetIsNormalized()},
			},
			"id": utils.ResourceIDAttribute(),
		},
	}
}

func (r *IntegrationConfluentResourceResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	accountID, resourceID, err := utils.AccountIDAndResourceIDFromID(request.ID)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), "")
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("account_id"), accountID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("resource_id"), resourceID)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), request.ID)...)
}

func (r *IntegrationConfluentResourceResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state IntegrationConfluentResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	accountID, resourceID, err := utils.AccountIDAndResourceIDFromID(state.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), "")
		return
	}

	resp, httpResp, err := r.Api.GetConfluentResource(r.Auth, accountID, resourceID)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.Append(utils.FrameworkErrorDiag(err, "error retrieving API Key"))
		return
	}
	if err := utils.CheckForUnparsed(resp); err != nil {
		response.Diagnostics.AddError("response contains unparsedObject", err.Error())
		return
	}

	r.updateState(ctx, &state, &resp)

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *IntegrationConfluentResourceResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var state IntegrationConfluentResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	accountId := state.AccountId.ValueString()

	body, diags := r.buildIntegrationConfluentResourceRequestBody(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, _, err := r.Api.CreateConfluentResource(r.Auth, accountId, *body)
	if err != nil {
		response.Diagnostics.Append(utils.FrameworkErrorDiag(err, "error retrieving IntegrationConfluentResource"))
		return
	}
	if err := utils.CheckForUnparsed(resp); err != nil {
		response.Diagnostics.AddError("response contains unparsedObject", err.Error())
		return
	}
	r.updateState(ctx, &state, &resp)

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *IntegrationConfluentResourceResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state IntegrationConfluentResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	accountID, resourceID, err := utils.AccountIDAndResourceIDFromID(state.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), "")
		return
	}

	body, diags := r.buildIntegrationConfluentResourceRequestBody(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resp, _, err := r.Api.UpdateConfluentResource(r.Auth, accountID, resourceID, *body)
	if err != nil {
		response.Diagnostics.Append(utils.FrameworkErrorDiag(err, "error retrieving IntegrationConfluentResource"))
		return
	}
	if err := utils.CheckForUnparsed(resp); err != nil {
		response.Diagnostics.AddError("response contains unparsedObject", err.Error())
		return
	}
	r.updateState(ctx, &state, &resp)

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *IntegrationConfluentResourceResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state IntegrationConfluentResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	accountID, resourceID, err := utils.AccountIDAndResourceIDFromID(state.ID.ValueString())
	if err != nil {
		response.Diagnostics.AddError(err.Error(), "")
		return
	}

	httpResp, err := r.Api.DeleteConfluentResource(r.Auth, accountID, resourceID)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}
		response.Diagnostics.Append(utils.FrameworkErrorDiag(err, "error deleting integration_confluent_resource"))
		return
	}
}

func (r *IntegrationConfluentResourceResource) updateState(ctx context.Context, state *IntegrationConfluentResourceModel, resp *datadogV2.ConfluentResourceResponse) {
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", state.AccountId.ValueString(), resp.Data.Id))

	data := resp.GetData()
	attributes := data.GetAttributes()

	if resourceType, ok := attributes.GetResourceTypeOk(); ok {
		state.ResourceType = types.StringValue(*resourceType)
	}

	if tags, ok := attributes.GetTagsOk(); ok && len(*tags) > 0 {
		state.Tags, _ = types.SetValueFrom(ctx, types.StringType, *tags)
	}
}

func (r *IntegrationConfluentResourceResource) buildIntegrationConfluentResourceRequestBody(ctx context.Context, state *IntegrationConfluentResourceModel) (*datadogV2.ConfluentResourceRequest, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	attributes := datadogV2.NewConfluentResourceRequestAttributesWithDefaults()

	if !state.ResourceType.IsNull() {
		attributes.SetResourceType(state.ResourceType.ValueString())
	}

	if !state.Tags.IsNull() {
		var tags []string
		diags.Append(state.Tags.ElementsAs(ctx, &tags, false)...)
		attributes.SetTags(tags)
	}

	req := datadogV2.NewConfluentResourceRequestWithDefaults()
	req.Data = *datadogV2.NewConfluentResourceRequestDataWithDefaults()
	req.Data.SetId(state.ResourceId.ValueString())
	req.Data.SetAttributes(*attributes)

	return req, diags
}
