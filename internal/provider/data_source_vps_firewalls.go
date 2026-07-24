// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DataSourceVPSFirewalls{}
	_ datasource.DataSourceWithConfigure = &DataSourceVPSFirewalls{}
)

func NewDataSourceVPSFirewalls() datasource.DataSource {
	return &DataSourceVPSFirewalls{}
}

// DataSourceVPSFirewalls defines the data source implementation.
type DataSourceVPSFirewalls struct {
	client *client.ClientWithResponses
}

// VPSFirewallRuleModel maps a single firewall rule from the API response.
type VPSFirewallRuleModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Action        types.String `tfsdk:"action"`
	Protocol      types.String `tfsdk:"protocol"`
	Port          types.String `tfsdk:"port"`
	Source        types.String `tfsdk:"source"`
	SourceDetail types.String `tfsdk:"source_detail"`
}

// VPSFirewallModel maps a single firewall from the API response.
type VPSFirewallModel struct {
	ID        types.Int64            `tfsdk:"id"`
	Name      types.String           `tfsdk:"name"`
	IsSynced  types.Bool             `tfsdk:"is_synced"`
	Rules     []VPSFirewallRuleModel `tfsdk:"rules"`
	CreatedAt timetypes.RFC3339      `tfsdk:"created_at"`
	UpdatedAt timetypes.RFC3339      `tfsdk:"updated_at"`
}

// DataSourceVPSFirewallsModel describes the data source data model.
type DataSourceVPSFirewallsModel struct {
	Firewalls []VPSFirewallModel `tfsdk:"firewalls"`
}

func (d *DataSourceVPSFirewalls) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_firewalls"
}

func (d *DataSourceVPSFirewalls) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"firewalls": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of firewall groups.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Firewall ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Firewall name.",
						},
						"is_synced": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the firewall is synced with the VPS.",
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
							Computed:            true,
							MarkdownDescription: "Firewall rules. Populated only when `firewall_id` is set.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed:            true,
										MarkdownDescription: "Firewall rule ID.",
									},
									"action": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Firewall rule action. Can be `accept` or `drop`.",
									},
									"protocol": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Firewall rule protocol (e.g. TCP, UDP, ICMP).",
									},
									"port": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Destination port or port range.",
									},
									"source": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Source of the rule. Can be `any` or `custom`.",
									},
									"source_detail": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Source detail of the rule. Populated when `source` is `custom`.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *DataSourceVPSFirewalls) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceVPSFirewalls) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceVPSFirewallsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &client.VPSGetFirewallListV1Params{}

	page := 1
	for {
		params.Page = &page

		response, err := d.client.VPSGetFirewallListV1WithResponse(ctx, params)
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

		if response.JSON200.Data != nil {
			for _, item := range *response.JSON200.Data {
				m := VPSFirewallModel{
					ID:        int64Value(item.Id),
					Name:      types.StringPointerValue(item.Name),
					IsSynced:  types.BoolPointerValue(item.IsSynced),
					CreatedAt: timetypes.NewRFC3339TimePointerValue(item.CreatedAt),
					Rules:     []VPSFirewallRuleModel{},
				}
				data.Firewalls = append(data.Firewalls, m)
			}

			meta := response.JSON200.Meta
			if meta == nil || meta.CurrentPage == nil || meta.PerPage == nil || meta.Total == nil {
				break
			}
			fetched := (*meta.CurrentPage) * (*meta.PerPage)
			if fetched >= *meta.Total {
				break
			}
			page++
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
