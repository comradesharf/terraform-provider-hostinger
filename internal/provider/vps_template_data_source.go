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
	_ datasource.DataSource              = &VPSTemplateDataSource{}
	_ datasource.DataSourceWithConfigure = &VPSTemplateDataSource{}
)

func NewVPSTemplateDataSource() datasource.DataSource {
	return &VPSTemplateDataSource{}
}

type VPSTemplateDataSource struct {
	client *client.ClientWithResponses
}

type VPSTemplateDataSourceModel struct {
	VPSVirtualMachineTemplateModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (d *VPSTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_template"
}

func (d *VPSTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a VPS operating system template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:            true,
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
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *VPSTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VPSTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VPSTemplateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ID.IsUnknown() || config.ID.IsNull() || config.ID.ValueInt64() <= 0 {
		resp.Diagnostics.AddError("Invalid VPS Template ID", "ID is unknown, null, or not a positive integer, so the VPS template cannot be read.")
		return
	}

	readTimeout, diags := config.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := d.client.VPSGetTemplateDetailsV1WithResponse(ctx, client.TemplateId(config.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read VPS Template", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Read VPS Template", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Read VPS Template", "Response body is nil")
		return
	}

	config.Merge(*response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
