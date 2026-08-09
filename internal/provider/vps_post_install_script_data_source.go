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
	_ datasource.DataSource              = &VPSPostInstallScriptDataSource{}
	_ datasource.DataSourceWithConfigure = &VPSPostInstallScriptDataSource{}
)

func NewVPSPostInstallScriptDataSource() datasource.DataSource {
	return &VPSPostInstallScriptDataSource{}
}

type VPSPostInstallScriptDataSource struct{ client *client.ClientWithResponses }

type VPSPostInstallScriptDataSourceModel struct {
	VPSPostInstallScriptModel
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (d *VPSPostInstallScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_post_install_script"
}

func (d *VPSPostInstallScriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a VPS post-install script.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:            true,
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
			"timeouts": timeouts.Attributes(ctx),
		}}
}

func (d *VPSPostInstallScriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VPSPostInstallScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VPSPostInstallScriptDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ID.IsUnknown() || config.ID.IsNull() || config.ID.ValueInt64() == 0 {
		resp.Diagnostics.AddError("Invalid Post-install Script ID", "ID is unknown, null, or zero, so the post-install script cannot be read.")
		return
	}

	readTimeout, diags := config.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := d.client.VPSGetPostInstallScriptV1WithResponse(ctx, client.PostInstallScriptId(config.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read VPS Post-install Script", fmt.Sprintf("Got error: %s", err))
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Unable to Read VPS Post-install Script", fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)))
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError("Unable to Read VPS Post-install Script", "Response body is nil")
		return
	}

	config.Merge(*response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
