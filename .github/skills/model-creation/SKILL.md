---
name: model-creation
description: 'Create a reusable Terraform framework model in an internal/provider *_model.go file. Use when mapping a Hostinger API response to Terraform types, adding a Merge method, or extracting shared data source/resource response models.'
argument-hint: '<model_name> e.g. vps_virtual_machine'
---

# Model Creation

Creates a reusable API response model for the Hostinger provider. Model files define Terraform framework field types and API-to-model conversion; data source and resource files define schemas, lifecycle behavior, timeouts, and API calls.

## When to Use

- Adding a new `internal/provider/<name>_model.go` file
- Mapping a generated client response to `terraform-plugin-framework` values
- Sharing a response model between a data source and a resource
- Adding nested response models and their conversion methods

## Prerequisites

- Identify the generated client response type in `internal/client/client.gen.go`
- Inspect the complete response shape, including pointers, slices, maps, enums, and `oneOf` wrapper types
- Inspect neighboring model files for naming and conversion patterns
- Decide whether the model is a reusable item model or a provider-specific wrapper model
- Identify whether resources need a separate `<PascalName>Identity` model for Terraform resource identity import

## Procedure

### 1. Inspect the generated response

Find the response type returned by the API endpoint and record:

- Exact generated Go type used by `JSON200` or a nested response field
- Pointer versus value semantics
- Primitive field types and generated enum types
- Collection element types
- `oneOf` or union fields and their generated `As<Type>()` conversion methods
- Nested response resources that deserve their own model

Use the generated type in the `Merge` method signature. Do not use `interface{}`, JSON maps, or hand-written API structs as a substitute.

### 2. Create the model file

Create `internal/provider/<name>_model.go`.

Key rules:

- Use `package provider`
- Match the copyright header used by neighboring model files:

  ```go
  // Copyright (c) HashiCorp, Inc.
  // SPDX-License-Identifier: MPL-2.0
  ```

- Name the primary model `<PascalName>Model`, for example `VPSVirtualMachineModel`
- Name nested models with the primary model prefix, for example `VPSVirtualMachineTemplateModel`
- Keep identity structs beside the related model, with `types.*` fields and `tfsdk` tags matching the identity schema
- Add `tfsdk:"<snake_case>"` to every model field
- Import the generated `internal/client` package and only the framework type packages required by the fields
- Add `func (m *<PascalName>Model) Merge(item client.<GeneratedResponseType>)`
- Keep schema declarations, data source/resource structs, timeouts, and API calls out of `_model.go`

### 3. Select framework field types

Use framework values rather than native primitives for scalar fields:

| API value | Model field | Conversion |
|---|---|---|
| `*string` or enum pointer | `types.String` | `types.StringPointerValue(v)` or `types.StringPointerValue((*string)(v))` |
| `*int` | `types.Int32` or `types.Int64` | Use the shared `int32Value(v)` or `int64Value(v)` helper |
| `*bool` | `types.Bool` | `types.BoolPointerValue(v)` |
| `*float64` | `types.Float64` | `types.Float64Value(*v)` with nil handling |
| `*time.Time` | `timetypes.RFC3339` | `timetypes.NewRFC3339TimePointerValue(v)` |
| IP address string | `iptypes.IPAddress`, `iptypes.IPv4Address`, or `iptypes.IPv6Address` | Use the matching `New...PointerValue` constructor |

Use `iptypes.IPAddress` when the API value can contain either address family. Use the family-specific type only when the API contract guarantees it.

For collections:

- `*[]T` becomes `[]TModel` or `[]types.String`, with a nil check before iteration
- `*map[string]interface{}` becomes `map[string]types.String` when the API map values are strings; preserve nil as nil
- Nested objects become `*NestedModel` when the response is optional, or `NestedModel` when it is always present

### 4. Implement `Merge`

Map every response field explicitly. Preserve null values by using framework null values instead of zero-value native types.

Typical scalar mapping:

```go
func (m *VPSFirewallModel) Merge(item client.VPSV1FirewallFirewallResource) {
    m.ID = int64Value(item.Id)
    m.Name = types.StringPointerValue(item.Name)
    m.IsSynced = types.BoolPointerValue(item.IsSynced)
    m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
}
```

For nested collections, instantiate the child model and call its `Merge` method:

```go
if item.Prices != nil {
    for _, price := range *item.Prices {
        var p BillingCatalogPriceModel
        p.Merge(price)
        m.Prices = append(m.Prices, p)
    }
}
```

For generated `oneOf` fields, call the generated `As<Type>()` method, check the error, and merge only the successfully converted value. Do not access union fields as though they were the concrete response type.

### 5. Connect the model

- Data source and resource models should embed or contain the reusable model instead of duplicating its fields
- Resource models may embed the reusable response model while a sibling identity model carries import-only identifiers
- Their schemas must mirror the model’s `tfsdk` tags and framework custom types
- A list data source model should contain a slice such as `[]VPSVirtualMachineModel`
- A single-item data source model may embed the reusable model and add only lookup fields and `timeouts`
- Keep `timeouts.Value` on the data source/resource model, not on the reusable API response model
- Keep identity fields separate from the reusable API response model when they are only needed to drive import and lifecycle state

## Verification

Run formatting and focused tests after editing:

```bash
gofmt -w internal/provider/<name>_model.go
go test ./internal/provider
```

If the model is used by an acceptance test, verify nested lists, nested objects, null values, custom IP values, and RFC3339 values with typed state checks.

## File Map

| File | Purpose |
|---|---|
| `internal/provider/<name>_model.go` | Reusable framework model and `Merge` methods |
| `internal/provider/<name>_data_source.go` | Data source schema, read behavior, filters, and timeouts |
| `internal/provider/<name>_resource.go` | Resource schema and lifecycle behavior |

## Reference

- [Data source creation skill](../data-source-creation/SKILL.md)
- [Schema type mappings](../data-source-creation/references/schema_types.md)
