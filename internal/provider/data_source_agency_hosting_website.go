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

// DataSourceAgencyHostingWebsiteModel describes the data source data model.
type DataSourceAgencyHostingWebsiteModel struct {
	WebsiteUID types.String `tfsdk:"website_uid"`
	UID        types.String `tfsdk:"uid"`
	Domain     types.String `tfsdk:"domain"`
	Status     types.String `tfsdk:"status"`
	PhpVersion types.String `tfsdk:"php_version"`
	Flavor     types.String `tfsdk:"flavor"`
	Datacenter types.String `tfsdk:"datacenter"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

func (d *DataSourceAgencyHostingWebsite) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agency_hosting_website"
}

func (d *DataSourceAgencyHostingWebsite) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads details for a single Agency Plan website.",
		Attributes: map[string]schema.Attribute{
			"website_uid": schema.StringAttribute{
				MarkdownDescription: "UID of the Agency Plan website.",
				Required:            true,
			},
			"uid": schema.StringAttribute{
				MarkdownDescription: "Website UID.",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Primary domain name of the website.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Website state (e.g. active, suspended).",
				Computed:            true,
			},
			"php_version": schema.StringAttribute{
				MarkdownDescription: "PHP version configured for the website.",
				Computed:            true,
			},
			"flavor": schema.StringAttribute{
				MarkdownDescription: "Setup flavor of the website.",
				Computed:            true,
			},
			"datacenter": schema.StringAttribute{
				MarkdownDescription: "Hostname of the server hosting the website.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "RFC3339 timestamp of when the website was created.",
				Computed:            true,
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

	ctx = tflog.SetField(ctx, "website_uid", data.WebsiteUID.ValueString())

	response, err := d.client.AgencyHostingGetAgencyPlanWebsiteDetailsV1WithResponse(ctx, data.WebsiteUID.ValueString())
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

	// Direct pointer value assignments.
	data.UID = types.StringPointerValue(item.Uid)
	data.Flavor = types.StringPointerValue(item.Flavor)

	// Conditional assignments.
	if item.State != nil {
		data.Status = types.StringValue(string(*item.State))
	} else {
		data.Status = types.StringNull()
	}

	if item.CreatedAt != nil {
		data.CreatedAt = types.StringValue(item.CreatedAt.Format(time.RFC3339))
	} else {
		data.CreatedAt = types.StringNull()
	}

	// Extract domain from the first entry in the Domains collection.
	if item.Domains != nil && len(*item.Domains) > 0 {
		data.Domain = types.StringPointerValue((*item.Domains)[0].Fqdn)
	} else {
		data.Domain = types.StringNull()
	}

	// Extract PHP version from the Settings.Php union type.
	if item.Settings != nil && item.Settings.Php != nil {
		phpSettings, err := item.Settings.Php.AsAgencyHostingV1WebsitesWebsitePhpSettingsResource()
		if err == nil {
			data.PhpVersion = types.StringPointerValue(phpSettings.Version)
		} else {
			data.PhpVersion = types.StringNull()
		}
	} else {
		data.PhpVersion = types.StringNull()
	}

	// Extract datacenter from the Server hostname.
	if item.Server != nil {
		data.Datacenter = types.StringPointerValue(item.Server.Hostname)
	} else {
		data.Datacenter = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
