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
	_ datasource.DataSource              = &DataSourceAgencyHostingDomains{}
	_ datasource.DataSourceWithConfigure = &DataSourceAgencyHostingDomains{}
)

func NewDataSourceAgencyHostingDomains() datasource.DataSource {
	return &DataSourceAgencyHostingDomains{}
}

// DataSourceAgencyHostingDomains defines the data source implementation.
type DataSourceAgencyHostingDomains struct {
	client *client.ClientWithResponses
}

// AgencyHostingDomainsItemModel maps a single domain from the API response.
type AgencyHostingDomainsItemModel struct {
	FQDN       types.String      `tfsdk:"fqdn"`
	WebsiteUID types.String      `tfsdk:"website_uid"`
	CreatedAt  timetypes.RFC3339 `tfsdk:"created_at"`
}

// DataSourceAgencyHostingDomainsModel describes the data source data model.
type DataSourceAgencyHostingDomainsModel struct {
	WebsiteUIDs []types.String                  `tfsdk:"website_uids"`
	Domains     []AgencyHostingDomainsItemModel `tfsdk:"domains"`
}

func (d *DataSourceAgencyHostingDomains) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agency_hosting_domains"
}

func (d *DataSourceAgencyHostingDomains) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists domains linked to Agency Plan websites.",
		Attributes: map[string]schema.Attribute{
			"website_uids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Filter domains to specific website UIDs.",
			},
			"domains": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of domains linked to Agency Plan websites.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fqdn": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Fully qualified domain name.",
						},
						"website_uid": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "UID of the owning website.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "RFC3339 timestamp of when the domain was created.",
							CustomType:          timetypes.RFC3339Type{},
						},
					},
				},
			},
		},
	}
}

func (d *DataSourceAgencyHostingDomains) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceAgencyHostingDomains) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceAgencyHostingDomainsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := client.AgencyHostingListAgencyPlanDomainsV1Params{}

	if len(data.WebsiteUIDs) > 0 {
		uids := make([]client.WebsiteUid, len(data.WebsiteUIDs))
		for i, uid := range data.WebsiteUIDs {
			uids[i] = uid.ValueString()
		}
		params.WebsiteUuids = &uids
		ctx = tflog.SetField(ctx, "website_uids", &params.WebsiteUuids)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	page := 0
	for {
		params.Page = &page

		response, err := d.client.AgencyHostingListAgencyPlanDomainsV1WithResponse(ctx, &params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Agency Hosting Domains",
				fmt.Sprintf("Got error: %s", err),
			)
			return
		}
		if response.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unable to Read Agency Hosting Domains",
				fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
			)
			return
		}
		if response.JSON200 == nil {
			resp.Diagnostics.AddError(
				"Unable to Read Agency Hosting Domains",
				"Response body is nil",
			)
			return
		}

		domains := response.JSON200.Data
		if domains == nil || len(*domains) == 0 {
			break
		}

		for _, item := range *domains {
			var m AgencyHostingDomainsItemModel
			m.FQDN = types.StringPointerValue(item.Fqdn)
			m.WebsiteUID = types.StringPointerValue(item.WebsiteUid)
			m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)

			data.Domains = append(data.Domains, m)
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
