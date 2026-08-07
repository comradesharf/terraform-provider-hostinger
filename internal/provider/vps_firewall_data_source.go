// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &VPSFirewallDataSource{}
	_ datasource.DataSourceWithConfigure = &VPSFirewallDataSource{}
)

func NewVPSFirewallDataSource() datasource.DataSource {
	return &VPSFirewallDataSource{}
}

// VPSFirewallDataSource defines the data source implementation.
type VPSFirewallDataSource struct {
	client *client.ClientWithResponses
}

// VPSFirewallDataSourceModel describes the data source data model.
type VPSFirewallDataSourceModel struct {
	VPSFirewallModel
	Rules    []VPSFirewallRuleModel `tfsdk:"rules"`
	Timeouts timeouts.Value         `tfsdk:"timeouts"`
}

func (d *VPSFirewallDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_firewall"
}

func (d *VPSFirewallDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Firewall ID",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Firewall name",
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
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *VPSFirewallDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = c
}

func (d *VPSFirewallDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VPSFirewallDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown ID",
			"ID is unknown, unable to read VPS virtual machine.",
		)
		return
	}

	if data.ID.IsNull() || data.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError(
			"Null ID",
			"ID is null or zero, unable to read VPS virtual machine.",
		)
		return
	}

	readTimeout, diags := data.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := d.client.VPSGetFirewallDetailsV1WithResponse(ctx, client.FirewallId(data.ID.ValueInt64()))
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

	data.Merge(*response.JSON200)

	if response.JSON200.Rules != nil {
		data.Rules = make([]VPSFirewallRuleModel, len(*response.JSON200.Rules))
		for i, rule := range *response.JSON200.Rules {
			data.Rules[i].Merge(rule)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
