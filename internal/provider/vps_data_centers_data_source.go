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
	_ datasource.DataSource              = &VPSDataCentersDataSource{}
	_ datasource.DataSourceWithConfigure = &VPSDataCentersDataSource{}
)

func NewVPSDataCentersDataSource() datasource.DataSource {
	return &VPSDataCentersDataSource{}
}

type VPSDataCentersDataSource struct {
	client *client.ClientWithResponses
}

type VPSDataCentersDataSourceModel struct {
	DataCenters []VPSDataCenterModel `tfsdk:"data_centers"`
	Timeouts    timeouts.Value       `tfsdk:"timeouts"`
}

func (d *VPSDataCentersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_data_centers"
}

func (d *VPSDataCentersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_centers": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of VPS data centers.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Data center ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Data center name.",
						},
						"city": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Data center location city.",
						},
						"continent": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Data center location continent.",
						},
						"location": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Data center location country (two-letter code).",
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *VPSDataCentersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VPSDataCentersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VPSDataCentersDataSourceModel
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

	response, err := d.client.VPSGetDataCenterListV1WithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Data Centers",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Data Centers",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Read VPS Data Centers",
			"Response body is nil",
		)
		return
	}

	for _, item := range *response.JSON200 {
		var dataCenter VPSDataCenterModel
		dataCenter.Merge(item)
		config.DataCenters = append(config.DataCenters, dataCenter)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
