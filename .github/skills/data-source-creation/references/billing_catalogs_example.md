# Billing Catalogs — Full Example

This is the complete, working implementation of the `billing_catalogs` data source. Use it as a reference when implementing a new data source.

## Implementation: `internal/provider/data_source_billing_catalogs.go`

```go
// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &DataSourceBillingCatalogs{}
	_ datasource.DataSourceWithConfigure = &DataSourceBillingCatalogs{}
)

func NewDataSourceBillingCatalogs() datasource.DataSource {
	return &DataSourceBillingCatalogs{}
}

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
	Metadata types.Map                    `tfsdk:"metadata"`
	Prices   []BillingCatalogsPricesModel `tfsdk:"prices"`
}

type DataSourceBillingCatalogsModel struct {
	BillingCatalogs []BillingCatalogsModel `tfsdk:"billing_catalogs"`
	Name            types.String           `tfsdk:"name"`
	Category        types.String           `tfsdk:"category"`
}
```

## Key patterns demonstrated

### Optional enum filter with validator
```go
"category": schema.StringAttribute{
    Optional: true,
    Validators: []validator.String{
        stringvalidator.OneOf("DOMAIN", "VPS"),
    },
},
```

### Nested list attribute
```go
"prices": schema.ListNestedAttribute{
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{ ... },
    },
    Computed: true,
},
```

### Map attribute (string → string)
```go
"metadata": schema.MapAttribute{
    ElementType: types.StringType,
    Computed:    true,
},
```

### Mapping a map field in Read
```go
if item.Metadata != nil {
    metadataMap := make(map[string]attr.Value, len(*item.Metadata))
    for k, v := range *item.Metadata {
        metadataMap[k] = types.StringValue(v.(string))
    }
    d.Metadata = types.MapValueMust(types.StringType, metadataMap)
} else {
    d.Metadata = types.MapNull(types.StringType)
}
```

### Casting a custom string type
```go
p.PeriodUnit = types.StringPointerValue((*string)(price.PeriodUnit))
```

### Casting custom int type
```go
p.Price = types.Int32Value(int32(*price.Price))
```
