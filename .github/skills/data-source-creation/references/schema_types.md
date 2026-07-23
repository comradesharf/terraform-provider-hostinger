# Schema Type Mappings

Reference for mapping Go client types to `terraform-plugin-framework` schema types and model fields.

## Primitive types

| Go type | Schema attribute | Model field | Conversion |
|---------|-----------------|-------------|------------|
| `*string` | `schema.StringAttribute` | `types.String` | `types.StringPointerValue(v)` |
| `string` | `schema.StringAttribute` | `types.String` | `types.StringValue(v)` |
| `*int` / `*int64` | `schema.Int64Attribute` | `types.Int64` | `types.Int64Value(int64(*v))` |
| `*int32` | `schema.Int32Attribute` | `types.Int32` | `types.Int32Value(int32(*v))` |
| `*float64` | `schema.Float64Attribute` | `types.Float64` | `types.Float64Value(*v)` |
| `*bool` | `schema.BoolAttribute` | `types.Bool` | `types.BoolPointerValue(v)` |
| custom string type `T` | `schema.StringAttribute` | `types.String` | `types.StringPointerValue((*string)(v))` |

## Datetime and network custom types

| API/Go type | Schema attribute | Model field | Conversion |
|-------------|------------------|-------------|------------|
| `*time.Time` / RFC3339 datetime pointer | `schema.StringAttribute{CustomType: timetypes.RFC3339Type{}}` | `timetypes.RFC3339` | `timetypes.NewRFC3339TimePointerValue(v)` |
| `*string` IPv4 | `schema.StringAttribute{CustomType: iptypes.IPv4AddressType{}}` | `iptypes.IPv4Address` | `iptypes.NewIPv4AddressPointerValue(v)` |
| `*string` IPv6 | `schema.StringAttribute{CustomType: iptypes.IPv6AddressType{}}` | `iptypes.IPv6Address` | `iptypes.NewIPv6AddressPointerValue(v)` |

Imports:

```go
import (
    "github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
    "github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)
```

## Collection types

| Go type | Schema attribute | Model field | Notes |
|---------|-----------------|-------------|-------|
| `*[]Item` | `schema.ListNestedAttribute` | `[]ItemModel` | iterate with `for _, item := range *v` |
| `*map[string]interface{}` | `schema.MapAttribute{ElementType: types.StringType}` | `types.Map` | use `types.MapValueMust`; set `types.MapNull` when nil |
| `*[]string` | `schema.ListAttribute{ElementType: types.StringType}` | `types.List` | use `types.ListValueMust` |

## Map field mapping pattern

```go
if item.Metadata != nil {
    m := make(map[string]attr.Value, len(*item.Metadata))
    for k, v := range *item.Metadata {
        m[k] = types.StringValue(v.(string))
    }
    model.Metadata = types.MapValueMust(types.StringType, m)
} else {
    model.Metadata = types.MapNull(types.StringType)
}
```

## Validators (import `github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator`)

```go
// Enum restriction
Validators: []validator.String{
    stringvalidator.OneOf("VALUE_A", "VALUE_B"),
},
```

## Sensitive attributes

```go
"api_token": schema.StringAttribute{
    Sensitive: true,
    Optional:  true,
},
```

## MarkdownDescription

Add `MarkdownDescription: "..."` to any attribute that benefits from user-facing docs:

```go
"id": schema.StringAttribute{
    MarkdownDescription: "Unique identifier of the item",
    Computed:            true,
},
```
