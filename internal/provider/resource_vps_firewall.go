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
			"rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Firewall rule ID",
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"action": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Firewall rule action",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"accept",
									"drop",
								),
							},
						},
						"protocol": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Firewall rule protocol",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"TCP",
									"UDP",
									"ICMP",
									"GRE",
									"any",
									"ESP",
									"AH",
									"ICMPv6",
									"SSH",
									"HTTP",
									"HTTPS",
									"MySQL",
									"PostgreSQL",
								),
							},
						},
						"port": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Firewall rule destination port: single or port range",
						},
						"source": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Firewall rule source. Can be `any` or `custom`",
						},
						"source_detail": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Firewall rule source detail. Can be `any` or IP address, CIDR or range",
						},
					},
				},
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
	var data ResourceVPSFirewallModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := client.VPSCreateNewFirewallV1JSONRequestBody{}
	p.Name = data.Name.ValueString()

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

	data.Merge(response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceVPSFirewall) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ResourceVPSFirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown ID",
			"ID is unknown, unable to read VPS firewall.",
		)
		return
	}

	if data.ID.IsNull() || data.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Null ID",
			"ID is null or zero, unable to read VPS firewall.",
		)
		return
	}

	response, err := r.client.VPSGetFirewallDetailsV1WithResponse(ctx, client.FirewallId(data.ID.ValueInt64()))
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

	data.Merge(response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceVPSFirewall) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ResourceVPSFirewallModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceVPSFirewall) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceVPSFirewallModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown ID",
			"ID is unknown, unable to read VPS firewall.",
		)
		return
	}

	if data.ID.IsNull() || data.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Null ID",
			"ID is null or zero, unable to read VPS firewall.",
		)
		return
	}

	response, err := r.client.VPSDeleteFirewallV1WithResponse(ctx, client.FirewallId(data.ID.ValueInt64()))
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
