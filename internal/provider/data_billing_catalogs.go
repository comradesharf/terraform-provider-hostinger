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
var _ datasource.DataSource = &BillingCatalogsDataSource{}

func NewBillingCatalogsDataSource() datasource.DataSource {
	return &BillingCatalogsDataSource{}
}

// BillingCatalogsDataSource defines the data source implementation.
type BillingCatalogsDataSource struct {
	client *client.ClientWithResponses
}

type BillingCatalogsModel struct {
	ID       types.String `tfsdk:"id"`
	Category types.String `tfsdk:"category"`
	Name     types.String `tfsdk:"name"`
	//Metadata types.Map    `tfsdk:"metadata"`
}

// billingCatalogsDataSourceModel describes the data source data model.
type billingCatalogsDataSourceModel struct {
	BillingCatalogs []BillingCatalogsModel `tfsdk:"billing_catalogs"`
}

func (d *BillingCatalogsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_catalogs"
}

func (d *BillingCatalogsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"billing_catalogs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Catalog item ID",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						//"metadata": schema.MapAttribute{
						//	ElementType: types.StringType,
						//	Computed:    true,
						//},
						//"prices": schema.ListNestedAttribute{
						//	NestedObject: schema.NestedAttributeObject{
						//		Attributes: map[string]schema.Attribute{
						//			"currency": schema.StringAttribute{
						//				Computed: true,
						//			},
						//			"first_period_price": schema.NumberAttribute{
						//				Computed: true,
						//			},
						//			"id": schema.StringAttribute{
						//				Computed: true,
						//			},
						//			"name": schema.StringAttribute{
						//				Computed: true,
						//			},
						//			"period": schema.NumberAttribute{
						//				Computed: true,
						//			},
						//			"period_unit": schema.StringAttribute{
						//				Computed: true,
						//			},
						//			"price": schema.NumberAttribute{
						//				Computed: true,
						//			},
						//		},
						//	},
						//	Computed: true,
						//},
					},
				},
			},
		},
	}
}

func (d *BillingCatalogsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BillingCatalogsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var params client.BillingGetCatalogItemListV1Params
	response, err := d.client.BillingGetCatalogItemListV1WithResponse(ctx, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Billing Catalogs",
			fmt.Sprintf("Got error: %s", err),
		)
		return
	}
	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unable to Read Billing Catalogs",
			fmt.Sprintf("Unexpected status code: %d", response.StatusCode()),
		)
		return
	}
	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Unable to Read Billing Catalogs",
			"Response body is nil",
		)
		return
	}

	var data billingCatalogsDataSourceModel
	for _, item := range *response.JSON200 {
		var d BillingCatalogsModel
		d.ID = types.StringPointerValue(item.Id)
		d.Category = types.StringPointerValue(item.Category)
		d.Name = types.StringPointerValue(item.Name)
		//
		//metadata := make(map[string]attr.Value, len(*item.Metadata))
		//for key, value := range *item.Metadata {
		//	metadata[key] = types.StringValue(value.(string))
		//}
		//
		//d.Metadata = types.MapValueMust(types.StringType, metadata)
		data.BillingCatalogs = append(data.BillingCatalogs, d)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
