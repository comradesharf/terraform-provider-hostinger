// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &VPSPostInstallScriptsDataSource{}
	_ datasource.DataSourceWithConfigure = &VPSPostInstallScriptsDataSource{}
)

func NewVPSPostInstallScriptsDataSource() datasource.DataSource {
	return &VPSPostInstallScriptsDataSource{}
}

type VPSPostInstallScriptsDataSource struct {
	client *client.ClientWithResponses
}

type VPSPostInstallScriptsDataSourceModel struct {
	Scripts  []VPSPostInstallScriptModel `tfsdk:"scripts"`
	Timeouts timeouts.Value              `tfsdk:"timeouts"`
}

func (d *VPSPostInstallScriptsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_post_install_scripts"
}

func (d *VPSPostInstallScriptsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists VPS post-install scripts.",
		Attributes: map[string]schema.Attribute{
			"scripts": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The VPS post-install scripts in the account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Post-install script ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Post-install script name.",
						},
						"content": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Post-install script content.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
							MarkdownDescription: "Timestamp when the post-install script was created (RFC3339).",
						},
						"updated_at": schema.StringAttribute{
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
							MarkdownDescription: "Timestamp when the post-install script was updated (RFC3339).",
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *VPSPostInstallScriptsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.ClientWithResponses)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}
	d.client = c
}

func (d *VPSPostInstallScriptsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VPSPostInstallScriptsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := config.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	params := client.VPSGetPostInstallScriptsV1Params{}

	page := 1
	for {
		params.Page = &page

		response, err := d.client.VPSGetPostInstallScriptsV1WithResponse(ctx, &params)
		if err != nil {
			resp.Diagnostics.AddError("Unable to Read VPS Post-install Scripts", fmt.Sprintf("Got error: %s", err))
			return
		}
		if response.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Post-install Scripts",
				fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
			)
			return
		}
		if response.JSON200 == nil {
			resp.Diagnostics.AddError("Unable to Read VPS Post-install Scripts", "Response body is nil")
			return
		}
		if response.JSON200.Data == nil || len(*response.JSON200.Data) == 0 {
			break
		}
		for _, item := range *response.JSON200.Data {
			var script VPSPostInstallScriptModel
			script.Merge(item)
			config.Scripts = append(config.Scripts, script)
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
