// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DataSourceAgencyHostingWebsite{}
	_ datasource.DataSourceWithConfigure = &DataSourceAgencyHostingWebsite{}
)

func NewDataSourceAgencyHostingWebsite() datasource.DataSource {
	return &DataSourceAgencyHostingWebsite{}
}

// DataSourceAgencyHostingWebsite defines the data source implementation.
type DataSourceAgencyHostingWebsite struct {
	client *client.ClientWithResponses
}

type AgencyHostingWebsiteSSLCertModelItem struct {
	Names     []types.String    `tfsdk:"names"`
	ExpiresAt timetypes.RFC3339 `tfsdk:"expires_at"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
}

type AgencyHostingWebsiteCustomSSLCertModelItem struct {
	IsExpired types.Bool        `tfsdk:"is_expired"`
	ExpiresAt timetypes.RFC3339 `tfsdk:"expires_at"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
}

type AgencyHostingWebsitePreviewDomainModelItem struct {
	FQDN      types.String      `tfsdk:"fqdn"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
}

type AgencyHostingWebsiteDomainModelItem struct {
	FQDN          types.String                                `tfsdk:"fqdn"`
	ParentFQDN    types.String                                `tfsdk:"parent_fqdn"`
	IPv6          iptypes.IPv6Address                         `tfsdk:"ipv6"`
	CreatedAt     timetypes.RFC3339                           `tfsdk:"created_at"`
	Nameservers   []types.String                              `tfsdk:"nameservers"`
	SSLCert       *AgencyHostingWebsiteSSLCertModelItem       `tfsdk:"ssl_cert"`
	CustomSSLCert *AgencyHostingWebsiteCustomSSLCertModelItem `tfsdk:"custom_ssl_cert"`
}

type AgencyHostingWebsiteSettingsPHPModelItem struct {
	Version types.String `tfsdk:"version"`
	Workers types.Int32  `tfsdk:"workers"`
}

type AgencyHostingWebsiteSettingsModelItem struct {
	PHP AgencyHostingWebsiteSettingsPHPModelItem `tfsdk:"php"`
}

type AgencyHostingWebsiteWordpressModelItem struct {
	Domain         types.String      `tfsdk:"domain"`
	Title          types.String      `tfsdk:"title"`
	Language       types.String      `tfsdk:"language"`
	IsConfigLocked types.Bool        `tfsdk:"is_config_locked"`
	CreatedAt      timetypes.RFC3339 `tfsdk:"created_at"`
}

type AgencyHostingWebsiteRemoteAccessSSHModelItem struct {
	Username          types.String `tfsdk:"username"`
	Host              types.String `tfsdk:"host"`
	Port              types.Int32  `tfsdk:"port"`
	IsEnabled         types.Bool   `tfsdk:"is_enabled"`
	IsPasswordEnabled types.Bool   `tfsdk:"is_password_enabled"`
}

type AgencyHostingWebsiteRemoteAccessSFTPModelItem struct {
	Username  types.String `tfsdk:"username"`
	Host      types.String `tfsdk:"host"`
	Port      types.Int32  `tfsdk:"port"`
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
}

type AgencyHostingWebsiteRemoteAccessModelItem struct {
	Mode types.String                                  `tfsdk:"mode"`
	SSH  AgencyHostingWebsiteRemoteAccessSSHModelItem  `tfsdk:"ssh"`
	SFTP AgencyHostingWebsiteRemoteAccessSFTPModelItem `tfsdk:"sftp"`
}

type AgencyHostingWebsiteServerModelItem struct {
	Hostname    types.String `tfsdk:"hostname"`
	CountryCode types.String `tfsdk:"country_code"`
}

type AgencyHostingWebsiteOrderPlanParametersModelItem struct {
	DiskQuotaBytes         types.Int64 `tfsdk:"disk_quota_bytes"`
	InodeQuota             types.Int64 `tfsdk:"inode_quota"`
	CPUCores               types.Int32 `tfsdk:"cpu_cores"`
	MemoryQuotaBytes       types.Int64 `tfsdk:"memory_quota_bytes"`
	DiskIOPSQuota          types.Int64 `tfsdk:"disk_iops_quota"`
	ProcessQuota           types.Int32 `tfsdk:"process_quota"`
	WebsiteQuota           types.Int32 `tfsdk:"website_quota"`
	MaxDatabasesPerWebsite types.Int32 `tfsdk:"max_databases_per_website"`
	IsCDNAvailable         types.Bool  `tfsdk:"is_cdn_available"`
}

type AgencyHostingWebsiteOrderPlanModelItem struct {
	Name       types.String                                     `tfsdk:"name"`
	Parameters AgencyHostingWebsiteOrderPlanParametersModelItem `tfsdk:"parameters"`
}

type AgencyHostingWebsiteOrderModelItem struct {
	ID        types.Int64                            `tfsdk:"id"`
	Status    types.String                           `tfsdk:"status"`
	CreatedAt timetypes.RFC3339                      `tfsdk:"created_at"`
	Plan      AgencyHostingWebsiteOrderPlanModelItem `tfsdk:"plan"`
}

type AgencyHostingWebsiteUserModelItem struct {
	Username types.String `tfsdk:"username"`
	State    types.String `tfsdk:"state"`
}

type AgencyHostingWebsiteStagingRootModelItem struct {
	UID types.String `tfsdk:"uid"`
}

// DataSourceAgencyHostingWebsiteModel describes the data source data model.
type DataSourceAgencyHostingWebsiteModel struct {
	UID           types.String                                `tfsdk:"uid"`
	IPv4          iptypes.IPv4Address                         `tfsdk:"ipv4"`
	Flavor        types.String                                `tfsdk:"flavor"`
	Type          types.String                                `tfsdk:"type"`
	Description   types.String                                `tfsdk:"description"`
	State         types.String                                `tfsdk:"state"`
	CreatedAt     timetypes.RFC3339                           `tfsdk:"created_at"`
	Domains       []AgencyHostingWebsiteDomainModelItem       `tfsdk:"domains"`
	PreviewDomain *AgencyHostingWebsitePreviewDomainModelItem `tfsdk:"preview_domain"`
	Settings      *AgencyHostingWebsiteSettingsModelItem      `tfsdk:"settings"`
	Wordpress     *AgencyHostingWebsiteWordpressModelItem     `tfsdk:"wordpress"`
	RemoteAccess  *AgencyHostingWebsiteRemoteAccessModelItem  `tfsdk:"remote_access"`
	Server        *AgencyHostingWebsiteServerModelItem        `tfsdk:"server"`
	Order         *AgencyHostingWebsiteOrderModelItem         `tfsdk:"order"`
	User          *AgencyHostingWebsiteUserModelItem          `tfsdk:"user"`
	StagingRoot   *AgencyHostingWebsiteStagingRootModelItem   `tfsdk:"staging_root"`
}

func (d *DataSourceAgencyHostingWebsite) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agency_hosting_website"
}

func (d *DataSourceAgencyHostingWebsite) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
		},
	}
}

func (d *DataSourceAgencyHostingWebsite) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceAgencyHostingWebsite) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceAgencyHostingWebsiteModel
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
			fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
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

	data.Flavor = types.StringPointerValue(item.Flavor)
	data.IPv4 = iptypes.NewIPv4AddressPointerValue(item.Ipv4)
	data.Flavor = types.StringPointerValue(item.Flavor)
	data.Type = types.StringPointerValue(item.Type)
	data.Description = types.StringPointerValue(item.Description)
	data.State = types.StringPointerValue((*string)(item.State))
	data.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)

	if item.Domains != nil {
		data.Domains = make([]AgencyHostingWebsiteDomainModelItem, len(*item.Domains))
		for i, domain := range *item.Domains {
			var d AgencyHostingWebsiteDomainModelItem
			d.FQDN = types.StringPointerValue(domain.Fqdn)
			d.ParentFQDN = types.StringPointerValue(domain.ParentFqdn)
			d.IPv6 = iptypes.NewIPv6AddressPointerValue(domain.Ipv6)
			d.CreatedAt = timetypes.NewRFC3339TimePointerValue(domain.CreatedAt)

			if domain.Nameservers != nil {
				d.Nameservers = make([]types.String, len(*domain.Nameservers))
				for j, nameserver := range *domain.Nameservers {
					d.Nameservers[j] = types.StringValue(nameserver)
				}
			}

			if domain.SslCert != nil {
				v, err := domain.SslCert.AsAgencyHostingV1WebsitesSslCertResource()
				if err == nil {
					var s AgencyHostingWebsiteSSLCertModelItem
					s.CreatedAt = timetypes.NewRFC3339TimePointerValue(v.CreatedAt)
					s.ExpiresAt = timetypes.NewRFC3339TimePointerValue(v.ExpiresAt)

					if v.Names != nil {
						s.Names = make([]types.String, len(*v.Names))
						for j, name := range *v.Names {
							s.Names[j] = types.StringValue(name)
						}
					}
					d.SSLCert = &s
				}
			}

			if domain.CustomSslCert != nil {
				v, err := domain.CustomSslCert.AsAgencyHostingV1WebsitesCustomSslCertResource()
				if err == nil {
					var s AgencyHostingWebsiteCustomSSLCertModelItem
					s.IsExpired = types.BoolPointerValue(v.IsExpired)
					s.CreatedAt = timetypes.NewRFC3339TimePointerValue(v.CreatedAt)
					s.ExpiresAt = timetypes.NewRFC3339TimePointerValue(v.ExpiresAt)
					d.CustomSSLCert = &s
				}
			}

			data.Domains[i] = d
		}
	}

	if item.PreviewDomain != nil {
		var d AgencyHostingWebsitePreviewDomainModelItem
		v, err := item.PreviewDomain.AsAgencyHostingV1WebsitesWebsitePreviewDomainResource()
		if err == nil {
			d.FQDN = types.StringPointerValue(v.Fqdn)
			d.CreatedAt = timetypes.NewRFC3339TimePointerValue(v.CreatedAt)
		}
		data.PreviewDomain = &d
	}

	if item.Settings != nil {
		var d AgencyHostingWebsiteSettingsModelItem
		v, err := item.Settings.Php.AsAgencyHostingV1WebsitesWebsitePhpSettingsResource()
		if err == nil {
			var p AgencyHostingWebsiteSettingsPHPModelItem
			p.Version = types.StringPointerValue(v.Version)
			p.Workers = int32Value(v.Workers)
			d.PHP = p
		}
		data.Settings = &d
	}

	if item.Wordpress != nil {
		var d AgencyHostingWebsiteWordpressModelItem
		v, err := item.Wordpress.AsAgencyHostingV1WebsitesWordPressInstallResource()
		if err == nil {
			d.Domain = types.StringPointerValue(v.Domain)
			d.Title = types.StringPointerValue(v.Title)
			d.Language = types.StringPointerValue(v.Language)
			d.IsConfigLocked = types.BoolPointerValue(v.IsConfigLocked)
			d.CreatedAt = timetypes.NewRFC3339TimePointerValue(v.CreatedAt)
		}
		data.Wordpress = &d
	}

	if item.RemoteAccess != nil {
		var d AgencyHostingWebsiteRemoteAccessModelItem
		d.Mode = types.StringPointerValue(item.RemoteAccess.Mode)

		if item.RemoteAccess.Ssh != nil {
			var ssh AgencyHostingWebsiteRemoteAccessSSHModelItem
			ssh.Username = types.StringPointerValue(item.RemoteAccess.Ssh.Username)
			ssh.Host = types.StringPointerValue(item.RemoteAccess.Ssh.Host)
			ssh.Port = int32Value(item.RemoteAccess.Ssh.Port)
			ssh.IsEnabled = types.BoolPointerValue(item.RemoteAccess.Ssh.IsEnabled)
			ssh.IsPasswordEnabled = types.BoolPointerValue(item.RemoteAccess.Ssh.IsPasswordEnabled)

			d.SSH = ssh
		}

		if item.RemoteAccess.Sftp != nil {
			var sftp AgencyHostingWebsiteRemoteAccessSFTPModelItem
			sftp.Username = types.StringPointerValue(item.RemoteAccess.Sftp.Username)
			sftp.Host = types.StringPointerValue(item.RemoteAccess.Sftp.Host)
			sftp.Port = int32Value(item.RemoteAccess.Sftp.Port)
			sftp.IsEnabled = types.BoolPointerValue(item.RemoteAccess.Sftp.IsEnabled)

			d.SFTP = sftp
		}

		if item.Server != nil {
			var s AgencyHostingWebsiteServerModelItem
			s.Hostname = types.StringPointerValue(item.Server.Hostname)
			s.CountryCode = types.StringPointerValue(item.Server.CountryCode)

			data.Server = &s
		}

		if item.Order != nil {
			var o AgencyHostingWebsiteOrderModelItem
			o.ID = int64Value(item.Order.Id)
			o.Status = types.StringPointerValue(item.Order.Status)
			o.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.Order.CreatedAt)

			if item.Order.Plan != nil {
				var p AgencyHostingWebsiteOrderPlanModelItem
				p.Name = types.StringPointerValue(item.Order.Plan.Name)

				if item.Order.Plan.Parameters != nil {
					var pp AgencyHostingWebsiteOrderPlanParametersModelItem
					pp.DiskQuotaBytes = int64Value(item.Order.Plan.Parameters.DiskQuotaBytes)
					pp.InodeQuota = int64Value(item.Order.Plan.Parameters.InodeQuota)
					pp.CPUCores = int32Value(item.Order.Plan.Parameters.CpuCores)
					pp.MemoryQuotaBytes = int64Value(item.Order.Plan.Parameters.MemoryQuotaBytes)
					pp.DiskIOPSQuota = int64Value(item.Order.Plan.Parameters.DiskIopsQuota)
					pp.ProcessQuota = int32Value(item.Order.Plan.Parameters.ProcessQuota)
					pp.WebsiteQuota = int32Value(item.Order.Plan.Parameters.WebsiteQuota)
					pp.MaxDatabasesPerWebsite = int32Value(item.Order.Plan.Parameters.MaxDatabasesPerWebsite)
					pp.IsCDNAvailable = types.BoolPointerValue(item.Order.Plan.Parameters.IsCdnAvailable)

					p.Parameters = pp
				}

				o.Plan = p
			}

			data.Order = &o
		}

		if item.User != nil {
			var u AgencyHostingWebsiteUserModelItem
			u.Username = types.StringPointerValue(item.User.Username)
			u.State = types.StringPointerValue(item.User.State)

			data.User = &u
		}

		if item.StagingRoot != nil {
			var s AgencyHostingWebsiteStagingRootModelItem
			v, err := item.StagingRoot.AsAgencyHostingV1WebsitesWebsiteStagingRootResource()
			if err == nil {
				s.UID = types.StringPointerValue(v.Uid)
			}
			data.StagingRoot = &s
		}

		data.RemoteAccess = &d
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
