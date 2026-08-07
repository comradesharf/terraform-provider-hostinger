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
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ReachSegmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &ReachSegmentsDataSource{}
)

func NewReachSegmentsDataSource() datasource.DataSource {
	return &ReachSegmentsDataSource{}
}

// ReachSegmentsDataSource defines the data source implementation.
type ReachSegmentsDataSource struct {
	client *client.ClientWithResponses
}

// ReachSegmentsDataSourceModel describes the data source data model.
type ReachSegmentsDataSourceModel struct {
	Segments []ReachSegmentsSegmentModel `tfsdk:"segments"`
	Timeouts timeouts.Value              `tfsdk:"timeouts"`
}

func (d *ReachSegmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reach_segments"
}

func (d *ReachSegmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists contact segments from the Hostinger Reach API.",
		Attributes: map[string]schema.Attribute{
			"segments": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of contact segments.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique identifier of the segment.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the segment.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "RFC3339 timestamp of when the segment was created.",
							CustomType:          timetypes.RFC3339Type{},
						},
						"updated_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "RFC3339 timestamp of when the segment was last updated.",
							CustomType:          timetypes.RFC3339Type{},
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *ReachSegmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ReachSegmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ReachSegmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := data.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	response, err := d.client.ReachListSegmentsV1WithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Reach Segments",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read Reach Segments",
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Read Reach Segments",
			"Response body is nil",
		)
		return
	}

	for _, item := range *response.JSON200 {
		var m ReachSegmentsSegmentModel
		m.Merge(item)
		data.Segments = append(data.Segments, m)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
