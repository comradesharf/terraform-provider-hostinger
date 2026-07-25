// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

		for _, item := range *response.JSON200.Data {
			var d VPSFirewallModel
			d.ID = int64Value(item.Id)
			d.Name = types.StringPointerValue(item.Name)
			d.IsSynced = types.BoolPointerValue(item.IsSynced)
			d.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
			d.UpdatedAt = timetypes.NewRFC3339TimePointerValue(item.UpdatedAt)

			if item.Rules != nil {
				for _, rule := range *item.Rules {
					var p VPSFirewallRuleModel
					p.ID = int64Value(rule.Id)
					p.Action = types.StringPointerValue((*string)(rule.Action))
					p.Protocol = types.StringPointerValue((*string)(rule.Protocol))
					p.Port = types.StringPointerValue(rule.Port)
					p.Source = types.StringPointerValue(rule.Source)
					p.SourceDetail = types.StringPointerValue(rule.SourceDetail)

					d.Rules = append(d.Rules, p)
				}
			}

			data.Firewalls = append(data.Firewalls, d)
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
