package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ResourceVPSFirewall{}
	_ resource.ResourceWithImportState = &ResourceVPSFirewall{}
)

func NewResourceVPSFirewall() resource.Resource {
	return &ResourceVPSFirewall{}
}

// ResourceVPSFirewall defines the resource implementation.
type ResourceVPSFirewall struct {
	client *client.ClientWithResponses
}

// ResourceVPSFirewallModel describes the resource data model.
type ResourceVPSFirewallModel struct {
	VPSFirewallModel
}

func (r *ResourceVPSFirewall) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_firewall"
}

func (r *ResourceVPSFirewall) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Firewall ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Firewall name",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_synced": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Is current firewall synced with VPS",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Timestamp when the firewall was created (RFC3339).",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "Timestamp when the firewall was updated (RFC3339).",
			},
		},
	}
}

func (r *ResourceVPSFirewall) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = c
}

func (r *ResourceVPSFirewall) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state ResourceVPSFirewallModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := client.VPSCreateNewFirewallV1JSONRequestBody{
		Name: state.Name.ValueString(),
	}

	response, err := r.client.VPSCreateNewFirewallV1WithResponse(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create VPS Firewall",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Create VPS Firewall",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Create VPS Firewall",
			"Response body is nil",
		)
		return
	}

	state.Merge(response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResourceVPSFirewall) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceVPSFirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown ID",
			"ID is unknown, unable to read VPS firewall.",
		)
		return
	}

	if state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Null ID",
			"ID is null or zero, unable to read VPS firewall.",
		)
		return
	}

	response, err := r.client.VPSGetFirewallDetailsV1WithResponse(ctx, client.FirewallId(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Firewalls",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Firewalls",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Firewalls",
			"Response body is nil",
		)
		return
	}

	state.Merge(response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResourceVPSFirewall) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceVPSFirewallModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResourceVPSFirewall) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceVPSFirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown ID",
			"ID is unknown, unable to read VPS firewall.",
		)
		return
	}

	if state.ID.IsNull() || state.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Null ID",
			"ID is null or zero, unable to read VPS firewall.",
		)
		return
	}

	response, err := r.client.VPSDeleteFirewallV1WithResponse(ctx, client.FirewallId(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete VPS Firewalls",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Delete VPS Firewalls",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
}

func (r *ResourceVPSFirewall) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStatePassthroughInt64ID(ctx, path.Root("id"), req, resp)
}
