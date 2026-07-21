// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
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
	Domain     types.String `tfsdk:"domain"`
	WebsiteUID types.String `tfsdk:"website_uid"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

// DataSourceAgencyHostingDomainsModel describes the data source data model.
type DataSourceAgencyHostingDomainsModel struct {
	WebsiteUIDs types.List                      `tfsdk:"website_uids"`
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
						"domain": schema.StringAttribute{
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

	if !data.WebsiteUIDs.IsNull() && !data.WebsiteUIDs.IsUnknown() {
		var uids []string
		resp.Diagnostics.Append(data.WebsiteUIDs.ElementsAs(ctx, &uids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(uids) > 0 {
			websiteUuids := client.WebsiteUuids(uids)
			params.WebsiteUuids = &websiteUuids
			ctx = tflog.SetField(ctx, "website_uids", uids)
		}
	}

	// Iterate all pages.
	page := 1
	for {
		p := client.Page(page)
		params.Page = &p

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

		if response.JSON200.Data == nil || len(*response.JSON200.Data) == 0 {
			break
		}

		for _, item := range *response.JSON200.Data {
			var m AgencyHostingDomainsItemModel
			m.Domain = types.StringPointerValue(item.Fqdn)
			m.WebsiteUID = types.StringPointerValue(item.WebsiteUid)
			if item.CreatedAt != nil {
				m.CreatedAt = types.StringValue(item.CreatedAt.Format(time.RFC3339))
			} else {
				m.CreatedAt = types.StringNull()
			}
			data.Domains = append(data.Domains, m)
		}

		// Check if there are more pages.
		if response.JSON200.Meta == nil || response.JSON200.Meta.CurrentPage == nil || response.JSON200.Meta.Total == nil || response.JSON200.Meta.PerPage == nil {
			break
		}
		currentPage := *response.JSON200.Meta.CurrentPage
		total := *response.JSON200.Meta.Total
		perPage := *response.JSON200.Meta.PerPage
		if perPage <= 0 || currentPage*perPage >= total {
			break
		}
		page++
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
