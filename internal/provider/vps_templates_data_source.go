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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &VPSTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &VPSTemplatesDataSource{}
)

func NewVPSTemplatesDataSource() datasource.DataSource {
	return &VPSTemplatesDataSource{}
}

type VPSTemplatesDataSource struct {
	client *client.ClientWithResponses
}

type VPSTemplatesDataSourceModel struct {
	Templates []VPSTemplateModel `tfsdk:"templates"`
	Timeouts  timeouts.Value     `tfsdk:"timeouts"`
}

func (d *VPSTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_templates"
}

func (d *VPSTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"templates": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of available VPS operating system templates.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Template ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Template name.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Template description.",
						},
						"documentation": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Link to the official operating system documentation.",
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *VPSTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VPSTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VPSTemplatesDataSourceModel
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

	response, err := d.client.VPSGetTemplatesV1WithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Templates",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Templates",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Read VPS Templates", "Response body is nil")
		return
	}

	for _, item := range *response.JSON200 {
		var template VPSTemplateModel
		template.Merge(item)
		config.Templates = append(config.Templates, template)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
