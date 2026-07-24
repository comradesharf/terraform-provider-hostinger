// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DataSourceVPSPublicKeys{}
	_ datasource.DataSourceWithConfigure = &DataSourceVPSPublicKeys{}
)

func NewDataSourceVPSPublicKeys() datasource.DataSource {
	return &DataSourceVPSPublicKeys{}
}

// DataSourceVPSPublicKeys defines the data source implementation.
type DataSourceVPSPublicKeys struct {
	client *client.ClientWithResponses
}

// VPSPublicKeysItemModel maps a single public key from the API response.
type VPSPublicKeysItemModel struct {
	ID        types.Int64       `tfsdk:"id"`
	Name      types.String      `tfsdk:"name"`
	Key       types.String      `tfsdk:"key"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
}

// DataSourceVPSPublicKeysModel describes the data source data model.
type DataSourceVPSPublicKeysModel struct {
	Name       types.String             `tfsdk:"name"`
	PublicKeys []VPSPublicKeysItemModel `tfsdk:"public_keys"`
}

func (d *DataSourceVPSPublicKeys) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_public_keys"
}

func (d *DataSourceVPSPublicKeys) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"public_keys": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"key": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed:   true,
							CustomType: timetypes.RFC3339Type{},
						},
					},
				},
			},
		},
	}
}

func (d *DataSourceVPSPublicKeys) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceVPSPublicKeys) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceVPSPublicKeysModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Name",
			"The 'name' attribute cannot be unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Name.IsNull() && data.Name.ValueString() != "" {
		ctx = tflog.SetField(ctx, "name", data.Name.ValueString())
	}

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

		publicKeys := response.JSON200.Data
		if publicKeys == nil || len(*publicKeys) == 0 {
			break
		}

		for _, item := range *publicKeys {
			m := VPSPublicKeysItemModel{
				ID:        int64Value(item.Id),
				Name:      types.StringPointerValue(item.Name),
				Key:       types.StringPointerValue(item.Key),
				CreatedAt: timetypes.NewRFC3339TimePointerValue(nil),
			}

			data.PublicKeys = append(data.PublicKeys, m)
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

	if !data.Name.IsNull() && data.Name.ValueString() != "" {
		filteredKeys := make([]VPSPublicKeysItemModel, 0, len(data.PublicKeys))
		for _, publicKey := range data.PublicKeys {
			if strings.Contains(publicKey.Name.ValueString(), data.Name.ValueString()) {
				filteredKeys = append(filteredKeys, publicKey)
			}
		}
		data.PublicKeys = filteredKeys
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
