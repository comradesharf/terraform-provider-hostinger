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

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &VPSPublicKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &VPSPublicKeysDataSource{}
)

func NewVPSPublicKeysDataSource() datasource.DataSource {
	return &VPSPublicKeysDataSource{}
}

// VPSPublicKeysDataSource defines the data source implementation.
type VPSPublicKeysDataSource struct {
	client *client.ClientWithResponses
}

// VPSPublicKeysDataSourceModel describes the data source data model.
type VPSPublicKeysDataSourceModel struct {
	PublicKeys []VPSPublicKeyModel `tfsdk:"public_keys"`
	Timeouts   timeouts.Value      `tfsdk:"timeouts"`
}

func (d *VPSPublicKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_public_keys"
}

func (d *VPSPublicKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"public_keys": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of SSH public keys in the account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Public key ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Public key name.",
						},
						"key": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Public key content.",
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *VPSPublicKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VPSPublicKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VPSPublicKeysDataSourceModel
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

	params := client.VPSGetPublicKeysV1Params{}

	page := 1
	for {
		params.Page = &page

		response, err := d.client.VPSGetPublicKeysV1WithResponse(ctx, &params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Public Keys",
				fmt.Sprintf("Got error: %s", err),
			)
			return
		}
		if response.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Public Keys",
				fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
			)
			return
		}
		if response.JSON200 == nil {
			resp.Diagnostics.AddError(
				"Unable to Read VPS Public Keys",
				"Response body is nil",
			)
			return
		}

		if response.JSON200.Data == nil || len(*response.JSON200.Data) == 0 {
			break
		}

		for _, item := range *response.JSON200.Data {
			var d VPSPublicKeyModel
			d.Merge(item)
			config.PublicKeys = append(config.PublicKeys, d)
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
