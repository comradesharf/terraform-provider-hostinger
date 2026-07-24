// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DataSourceBillingCatalogs{}
	_ datasource.DataSourceWithConfigure = &DataSourceBillingCatalogs{}
)

func NewDataSourceBillingCatalogs() datasource.DataSource {
	return &DataSourceBillingCatalogs{}
}

// DataSourceBillingCatalogs defines the data source implementation.
type DataSourceBillingCatalogs struct {
	client *client.ClientWithResponses
}

type BillingCatalogsPricesModel struct {
	Currency         types.String `tfsdk:"currency"`
	FirstPeriodPrice types.Int32  `tfsdk:"first_period_price"`
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Period           types.Int32  `tfsdk:"period"`
	PeriodUnit       types.String `tfsdk:"period_unit"`
	Price            types.Int32  `tfsdk:"price"`
}

type BillingCatalogsModel struct {
	ID       types.String                 `tfsdk:"id"`
	Category types.String                 `tfsdk:"category"`
	Name     types.String                 `tfsdk:"name"`
	Metadata map[string]types.String      `tfsdk:"metadata"`
	Prices   []BillingCatalogsPricesModel `tfsdk:"prices"`
}

// DataSourceBillingCatalogsModel describes the data source data model.
type DataSourceBillingCatalogsModel struct {
	BillingCatalogs []BillingCatalogsModel `tfsdk:"billing_catalogs"`
	Name            types.String           `tfsdk:"name"`
	Category        types.String           `tfsdk:"category"`
}

func (d *DataSourceBillingCatalogs) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_catalogs"
}

func (d *DataSourceBillingCatalogs) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
			},
			"category": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("DOMAIN", "VPS"),
				},
			},
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
						"metadata": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"prices": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"currency": schema.StringAttribute{
										Computed: true,
									},
									"first_period_price": schema.Int32Attribute{
										Computed: true,
									},
									"id": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
									"period": schema.Int32Attribute{
										Computed: true,
									},
									"period_unit": schema.StringAttribute{
										Computed: true,
									},
									"price": schema.Int32Attribute{
										Computed: true,
									},
								},
							},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *DataSourceBillingCatalogs) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceBillingCatalogs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceBillingCatalogsModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Name",
			"The 'name' attribute cannot be unknown",
		)
	}

	if data.Category.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Category",
			"The 'category' attribute cannot be unknown",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	params := client.BillingGetCatalogItemListV1Params{}

	if !data.Name.IsNull() && data.Name.ValueString() != "" {
		params.Name = data.Name.ValueStringPointer()
		ctx = tflog.SetField(ctx, "name", *params.Name)
	}

	if !data.Category.IsNull() && data.Category.ValueString() != "" {
		params.Category = (*client.BillingGetCatalogItemListV1ParamsCategory)(data.Category.ValueStringPointer())
		ctx = tflog.SetField(ctx, "category", &params.Category)
	}

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
			fmt.Sprintf("Unexpected status code: %d, response: %s", response.StatusCode(), string(response.Body)),
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

	for _, item := range *response.JSON200 {
		var d BillingCatalogsModel
		d.ID = types.StringPointerValue(item.Id)
		d.Category = types.StringPointerValue(item.Category)
		d.Name = types.StringPointerValue(item.Name)

		if item.Metadata != nil {
			d.Metadata = make(map[string]types.String, len(*item.Metadata))
			for k, v := range *item.Metadata {
				if s, ok := v.(string); ok {
					d.Metadata[k] = types.StringValue(s)
				}
			}
		}

		if item.Prices != nil {
			for _, price := range *item.Prices {
				var p BillingCatalogsPricesModel
				p.Currency = types.StringPointerValue(price.Currency)
				p.FirstPeriodPrice = int32Value(price.FirstPeriodPrice)
				p.ID = types.StringPointerValue(price.Id)
				p.Name = types.StringPointerValue(price.Name)
				p.Period = int32Value(price.Period)
				p.PeriodUnit = types.StringPointerValue((*string)(price.PeriodUnit))
				p.Price = int32Value(price.Price)

				d.Prices = append(d.Prices, p)
			}
		}

		data.BillingCatalogs = append(data.BillingCatalogs, d)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
