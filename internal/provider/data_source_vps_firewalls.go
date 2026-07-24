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
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	ID       types.Int64  `tfsdk:"id"`
	Protocol types.String `tfsdk:"protocol"`
	Port     types.String `tfsdk:"port"`
	Source   types.String `tfsdk:"source"`
	Action   types.String `tfsdk:"action"`
}

// VPSFirewallModel maps a single firewall from the API response.
type VPSFirewallModel struct {
	ID        types.Int64       `tfsdk:"id"`
	Name      types.String      `tfsdk:"name"`
	IsSynced  types.Bool        `tfsdk:"is_synced"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
	Rules     []VPSFirewallRuleModel `tfsdk:"rules"`
}

// DataSourceVPSFirewallsModel describes the data source data model.
type DataSourceVPSFirewallsModel struct {
	FirewallID types.Int64        `tfsdk:"firewall_id"`
	Firewalls  []VPSFirewallModel `tfsdk:"firewalls"`
}

func (d *DataSourceVPSFirewalls) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_firewalls"
}

func (d *DataSourceVPSFirewalls) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"firewall_id": schema.Int64Attribute{
				Optional: true,
			},
			"firewalls": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"is_synced": schema.BoolAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed:   true,
							CustomType: timetypes.RFC3339Type{},
						},
						"rules": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed: true,
									},
									"protocol": schema.StringAttribute{
										Computed: true,
									},
									"port": schema.StringAttribute{
										Computed: true,
									},
									"source": schema.StringAttribute{
										Computed: true,
									},
									"action": schema.StringAttribute{
										Computed: true,
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

	if data.FirewallID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Firewall ID",
			"The 'firewall_id' attribute cannot be unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.FirewallID.IsNull() {
		firewallID := int(data.FirewallID.ValueInt64())
		ctx = tflog.SetField(ctx, "firewall_id", firewallID)

		response, err := d.client.VPSGetFirewallDetailsV1WithResponse(ctx, firewallID)
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
				fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
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

		item := response.JSON200
		m := VPSFirewallModel{
			ID:        int64Value(item.Id),
			Name:      types.StringPointerValue(item.Name),
			IsSynced:  types.BoolPointerValue(item.IsSynced),
			CreatedAt: timetypes.NewRFC3339TimePointerValue(item.CreatedAt),
		}

		if item.Rules != nil {
			for _, rule := range *item.Rules {
				r := VPSFirewallRuleModel{
					ID:       int64Value(rule.Id),
					Protocol: types.StringPointerValue((*string)(rule.Protocol)),
					Port:     types.StringPointerValue(rule.Port),
					Source:   types.StringPointerValue(rule.Source),
					Action:   types.StringPointerValue((*string)(rule.Action)),
				}
				m.Rules = append(m.Rules, r)
			}
		}

		if m.Rules == nil {
			m.Rules = []VPSFirewallRuleModel{}
		}

		data.Firewalls = append(data.Firewalls, m)
	} else {
		response, err := d.client.VPSGetFirewallListV1WithResponse(ctx, &client.VPSGetFirewallListV1Params{})
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
				fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
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
		}
	}

	if data.Firewalls == nil {
		data.Firewalls = []VPSFirewallModel{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
