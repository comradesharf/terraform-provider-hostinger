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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &BillingCatalogsDataSource{}
	_ datasource.DataSourceWithConfigure = &BillingCatalogsDataSource{}
)

func NewBillingCatalogsDataSource() datasource.DataSource {
	return &BillingCatalogsDataSource{}
}

// BillingCatalogsDataSource defines the data source implementation.
type BillingCatalogsDataSource struct {
	client *client.ClientWithResponses
}

// BillingCatalogsDataSourceModel describes the data source data model.
type BillingCatalogsDataSourceModel struct {
	BillingCatalogs []BillingCatalogModel `tfsdk:"billing_catalogs"`
	Name            types.String          `tfsdk:"name"`
	Category        types.String          `tfsdk:"category"`
	Timeouts        timeouts.Value        `tfsdk:"timeouts"`
}

func (d *BillingCatalogsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_catalogs"
}

func (d *BillingCatalogsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"timeouts": timeouts.Attributes(ctx),
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
	var config BillingCatalogsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Name",
			"The 'name' attribute cannot be unknown",
		)
	}

	if config.Category.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Category",
			"The 'category' attribute cannot be unknown",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	params := client.BillingGetCatalogItemListV1Params{}

	if !config.Name.IsNull() && config.Name.ValueString() != "" {
		params.Name = config.Name.ValueStringPointer()
		ctx = tflog.SetField(ctx, "name", *params.Name)
	}

	if !config.Category.IsNull() && config.Category.ValueString() != "" {
		params.Category = (*client.BillingGetCatalogItemListV1ParamsCategory)(config.Category.ValueStringPointer())
		ctx = tflog.SetField(ctx, "category", &params.Category)
	}

	readTimeout, diags := config.Timeouts.Read(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

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
		var d BillingCatalogModel
		d.Merge(item)
		config.BillingCatalogs = append(config.BillingCatalogs, d)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
