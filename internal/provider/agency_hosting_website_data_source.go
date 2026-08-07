// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &AgencyHostingWebsiteDataSource{}
	_ datasource.DataSourceWithConfigure = &AgencyHostingWebsiteDataSource{}
)

func NewAgencyHostingWebsiteDataSource() datasource.DataSource {
	return &AgencyHostingWebsiteDataSource{}
}

// AgencyHostingWebsiteDataSource defines the data source implementation.
type AgencyHostingWebsiteDataSource struct {
	client *client.ClientWithResponses
}

// AgencyHostingWebsiteDataSourceModel describes the data source data model.
type AgencyHostingWebsiteDataSourceModel struct {
	AgencyHostingWebsiteModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (d *AgencyHostingWebsiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agency_hosting_website"
}

func (d *AgencyHostingWebsiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads details for a single Agency Plan website.",
		Attributes: map[string]schema.Attribute{
			"uid": schema.StringAttribute{
				MarkdownDescription: "Website UID.",
				Required:            true,
			},
			"ipv4": schema.StringAttribute{
				MarkdownDescription: "Website IPv4 address.",
				Computed:            true,
				CustomType:          iptypes.IPv4AddressType{},
			},
			"flavor": schema.StringAttribute{
				MarkdownDescription: "Website flavor.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Website type.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Website description.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "Website state.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Website creation timestamp.",
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
			},
			"domains": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fqdn": schema.StringAttribute{
							MarkdownDescription: "Fully qualified domain name.",
							Computed:            true,
						},
						"parent_fqdn": schema.StringAttribute{
							MarkdownDescription: "Parent fully qualified domain name.",
							Computed:            true,
						},
						"ipv6": schema.StringAttribute{
							MarkdownDescription: "IPv6 address.",
							Computed:            true,
							CustomType:          iptypes.IPv6AddressType{},
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "Domain creation timestamp.",
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
						},
						"nameservers": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of nameservers for the domain.",
							Computed:            true,
						},
						"ssl_cert": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"names": schema.ListAttribute{
									ElementType:         types.StringType,
									MarkdownDescription: "List of names covered by the SSL certificate.",
									Computed:            true,
								},
								"expires_at": schema.StringAttribute{
									MarkdownDescription: "SSL certificate expiration timestamp.",
									Computed:            true,
									CustomType:          timetypes.RFC3339Type{},
								},
								"created_at": schema.StringAttribute{
									MarkdownDescription: "SSL certificate creation timestamp.",
									Computed:            true,
									CustomType:          timetypes.RFC3339Type{},
								},
							},
						},
						"custom_ssl_cert": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"is_expired": schema.BoolAttribute{
									MarkdownDescription: "Indicates if the custom SSL certificate is expired.",
									Computed:            true,
								},
								"expires_at": schema.StringAttribute{
									MarkdownDescription: "Custom SSL certificate expiration timestamp.",
									Computed:            true,
									CustomType:          timetypes.RFC3339Type{},
								},
								"created_at": schema.StringAttribute{
									MarkdownDescription: "Custom SSL certificate creation timestamp.",
									Computed:            true,
									CustomType:          timetypes.RFC3339Type{},
								},
							},
						},
					},
				},
			},
			"preview_domain": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"fqdn": schema.StringAttribute{
						MarkdownDescription: "Fully qualified domain name for the preview domain.",
						Computed:            true,
					},
					"created_at": schema.StringAttribute{
						MarkdownDescription: "Preview domain creation timestamp.",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
				},
			},
			"settings": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"php": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"version": schema.StringAttribute{
								MarkdownDescription: "PHP version.",
								Computed:            true,
							},
							"workers": schema.Int64Attribute{
								MarkdownDescription: "Number of PHP workers.",
								Computed:            true,
							},
						},
					},
				},
			},
			"wordpress": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"domain": schema.StringAttribute{
						MarkdownDescription: "WordPress domain.",
						Computed:            true,
					},
					"title": schema.StringAttribute{
						MarkdownDescription: "WordPress site title.",
						Computed:            true,
					},
					"language": schema.StringAttribute{
						MarkdownDescription: "WordPress site language.",
						Computed:            true,
					},
					"is_config_locked": schema.BoolAttribute{
						MarkdownDescription: "Indicates if WordPress configuration is locked.",
						Computed:            true,
					},
					"created_at": schema.StringAttribute{
						MarkdownDescription: "WordPress creation timestamp.",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
				},
			},
			"remote_access": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						MarkdownDescription: "Remote access mode.",
						Computed:            true,
					},
					"ssh": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								MarkdownDescription: "SSH username.",
								Computed:            true,
							},
							"host": schema.StringAttribute{
								MarkdownDescription: "SSH host.",
								Computed:            true,
							},
							"port": schema.Int64Attribute{
								MarkdownDescription: "SSH port.",
								Computed:            true,
							},
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Indicates if SSH is enabled.",
								Computed:            true,
							},
							"is_password_enabled": schema.BoolAttribute{
								MarkdownDescription: "Indicates if SSH password authentication is enabled.",
								Computed:            true,
							},
						},
					},
					"sftp": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								MarkdownDescription: "SFTP username.",
								Computed:            true,
							},
							"host": schema.StringAttribute{
								MarkdownDescription: "SFTP host.",
								Computed:            true,
							},
							"port": schema.Int64Attribute{
								MarkdownDescription: "SFTP port.",
								Computed:            true,
							},
							"is_enabled": schema.BoolAttribute{
								MarkdownDescription: "Indicates if SFTP is enabled.",
								Computed:            true,
							},
						},
					},
				},
			},
			"server": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"hostname": schema.StringAttribute{
						MarkdownDescription: "Server hostname.",
						Computed:            true,
					},
					"country_code": schema.StringAttribute{
						MarkdownDescription: "Server country code.",
						Computed:            true,
					},
				},
			},
			"order": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						MarkdownDescription: "Order ID.",
						Computed:            true,
					},
					"status": schema.StringAttribute{
						MarkdownDescription: "Order status.",
						Computed:            true,
					},
					"created_at": schema.StringAttribute{
						MarkdownDescription: "Order creation timestamp.",
						Computed:            true,
						CustomType:          timetypes.RFC3339Type{},
					},
					"plan": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								MarkdownDescription: "Plan name.",
								Computed:            true,
							},
							"parameters": schema.SingleNestedAttribute{
								Computed: true,
								Attributes: map[string]schema.Attribute{
									"disk_quota_bytes": schema.Int64Attribute{
										MarkdownDescription: "Disk quota in bytes.",
										Computed:            true,
									},
									"inode_quota": schema.Int64Attribute{
										MarkdownDescription: "Inode quota.",
										Computed:            true,
									},
									"cpu_cores": schema.Int64Attribute{
										MarkdownDescription: "CPU cores.",
										Computed:            true,
									},
									"memory_quota_bytes": schema.Int64Attribute{
										MarkdownDescription: "Memory quota in bytes.",
										Computed:            true,
									},
									"disk_iops_quota": schema.Int64Attribute{
										MarkdownDescription: "Disk IOPS quota.",
										Computed:            true,
									},
									"process_quota": schema.Int64Attribute{
										MarkdownDescription: "Process quota.",
										Computed:            true,
									},
									"website_quota": schema.Int64Attribute{
										MarkdownDescription: "Website quota.",
										Computed:            true,
									},
									"max_databases_per_website": schema.Int64Attribute{
										MarkdownDescription: "Maximum databases per website.",
										Computed:            true,
									},
									"is_cdn_available": schema.BoolAttribute{
										MarkdownDescription: "Indicates if CDN is available.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
			"user": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						MarkdownDescription: "User username.",
						Computed:            true,
					},
					"state": schema.StringAttribute{
						MarkdownDescription: "User state.",
						Computed:            true,
					},
				},
			},
			"staging_root": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"uid": schema.StringAttribute{
						MarkdownDescription: "Staging root UID.",
						Computed:            true,
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *AgencyHostingWebsiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AgencyHostingWebsiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AgencyHostingWebsiteDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.UID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Website UID",
			"The 'uid' attribute cannot be unknown.",
		)
		return
	}

	if data.UID.IsNull() || data.UID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Website UID",
			"The 'uid' attribute must be a non-empty string.",
		)
		return
	}

	ctx = tflog.SetField(ctx, "website_uid", data.UID.ValueString())

	readTimeout, diags := data.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := d.client.AgencyHostingGetAgencyPlanWebsiteDetailsV1WithResponse(ctx, data.UID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Agency Hosting Website",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read Agency Hosting Website",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Read Agency Hosting Website",
			"Response body is nil",
		)
		return
	}

	item := response.JSON200
	data.Merge(*item)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
