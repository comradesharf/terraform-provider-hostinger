// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DataSourceReachSegments{}
	_ datasource.DataSourceWithConfigure = &DataSourceReachSegments{}
)

func NewDataSourceReachSegments() datasource.DataSource {
	return &DataSourceReachSegments{}
}

// DataSourceReachSegments defines the data source implementation.
type DataSourceReachSegments struct {
	client *client.ClientWithResponses
}

// ReachSegmentsItemModel maps a single segment from the API response.
type ReachSegmentsItemModel struct {
	Uuid      types.String `tfsdk:"uuid"`
	Name      types.String `tfsdk:"name"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// DataSourceReachSegmentsModel describes the data source data model.
type DataSourceReachSegmentsModel struct {
	Segments []ReachSegmentsItemModel `tfsdk:"segments"`
}

func (d *DataSourceReachSegments) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reach_segments"
}

func (d *DataSourceReachSegments) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						},
						"updated_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "RFC3339 timestamp of when the segment was last updated.",
						},
					},
				},
			},
		},
	}
}

func (d *DataSourceReachSegments) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceReachSegments) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceReachSegmentsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
			fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
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
		var m ReachSegmentsItemModel
		m.Uuid = types.StringPointerValue(item.Uuid)
		m.Name = types.StringPointerValue(item.Name)
		if item.CreatedAt != nil {
			m.CreatedAt = types.StringValue(item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
		} else {
			m.CreatedAt = types.StringNull()
		}
		if item.UpdatedAt != nil {
			m.UpdatedAt = types.StringValue(item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
		} else {
			m.UpdatedAt = types.StringNull()
		}

		data.Segments = append(data.Segments, m)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
