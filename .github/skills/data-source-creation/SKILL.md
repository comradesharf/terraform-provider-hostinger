---
name: data-source-creation
description: 'Create a new Terraform data source for the Hostinger provider. Use when adding a new data source, implementing a read-only data source, scaffolding a data source, or wiring up a new API endpoint as a data source.'
argument-hint: '<data_source_name> e.g. billing_subscriptions'
---

# Data Source Creation

Creates a new Terraform data source for the Hostinger provider following the established `terraform-plugin-framework` patterns used by the current provider implementations.

## When to Use

- Adding a new read-only data source backed by the Hostinger API
- Scaffolding a data source implementation + test + docs

## Prerequisites

- Identify the **data source name** (snake_case, e.g. `billing_subscriptions`)
- Identify the **client method** in `internal/client/client.gen.go` (e.g. `BillingGetSubscriptionListV1WithResponse`)
- Know the **response struct** and its fields to map to the Terraform schema
- Determine whether the response needs a reusable `internal/provider/<name>_model.go` model; use the [model-creation skill](../model-creation/SKILL.md) for that mapping

## Procedure

### 1. Explore the client and existing models

Search `internal/client/client.gen.go` for the relevant API method and its response/params types to understand:
- Method signature (`WithResponse` variant)
- `Params` struct (optional query filters → become `Optional` schema attributes)
- Response struct (`JSON200`) and its nested types → become `Computed` schema attributes

Also inspect the closest existing data source and any shared response model in `internal/provider`. Reuse a shared model when one already represents the API response (for example, `VPSVirtualMachineModel`), rather than creating a second representation. If no suitable model exists, create it before completing the data source.

### 2. Create the implementation file

Create `internal/provider/<name>_data_source.go` following the [implementation template](./assets/data_source.go.tmpl), adapting the template to the current naming and timeout conventions below.

Key rules:
- Copyright header: `// Copyright IBM Corp. 2021, 2025\n// SPDX-License-Identifier: MPL-2.0`
- Package: `package provider`
- Struct name pattern: `<PascalName>DataSource` (e.g. `BillingCatalogsDataSource`)
- Constructor: `New<PascalName>DataSource() datasource.DataSource`
- Add compile-time interface assertions for `datasource.DataSource` and `datasource.DataSourceWithConfigure`
- Keep the API client on the data source as `client *client.ClientWithResponses`
- Define a `<PascalName>DataSourceModel` with `tfsdk` tags. Keep reusable API fields in `<name>_model.go` and use a slice of item models for list results; for a single result, embed or reuse the corresponding shared model when appropriate.
- `Metadata`: set `resp.TypeName = req.ProviderTypeName + "_<name>"`
- `Configure`: extract `*client.ClientWithResponses` from `req.ProviderData`
- `Schema`: filter params → `Optional`, response fields → `Computed`; represent API arrays as `schema.ListNestedAttribute`, nested objects as `schema.SingleNestedAttribute`, and nested arrays with their own `NestedAttributeObject`
- Include a `Timeouts timeouts.Value` model field tagged `tfsdk:"timeouts"` and a `"timeouts": timeouts.Attributes(ctx)` schema attribute for API reads. Resolve the read timeout with `data.Timeouts.Read(ctx, 20*time.Minute)` and apply it with `context.WithTimeout`.
- `Read`: load configuration into the model, reject unknown optional filter values before the API call, build any params, call the `WithResponse` method, check `StatusCode() != http.StatusOK`, reject a nil `JSON200`, map the response into model structs, and persist with `resp.State.Set`
- Use `types.StringPointerValue`, `types.Int32Value`, `types.MapValueMust`, etc. for type conversions
- Use `timetypes.RFC3339` for datetime model fields, `schema.StringAttribute{CustomType: timetypes.RFC3339Type{}}` in schema, and `timetypes.NewRFC3339TimePointerValue(...)` when mapping API response values
- Use `terraform-plugin-framework-nettypes` custom types for network values. For values that may be either IPv4 or IPv6, use `iptypes.IPAddress` with `iptypes.IPAddressType{}`; use `iptypes.IPv4Address` / `iptypes.IPv6Address` when the API guarantees the address family. Set the matching schema `CustomType` and map with the corresponding constructors.
- Use `tflog.SetField` to attach filter values to context before the API call

For a collection data source with no filters, the current shape is:

```go
type VPSVirtualMachinesDataSourceModel struct {
    VirtualMachines []VPSVirtualMachineModel `tfsdk:"virtual_machines"`
    Timeouts        timeouts.Value           `tfsdk:"timeouts"`
}
```

Its schema exposes `virtual_machines` as a computed `schema.ListNestedAttribute`, with nested computed attributes for primitive fields, nested lists, nested objects, custom IP types, and RFC3339 timestamps. The read method calls the list endpoint without params and appends each mapped item to the result slice.

For paginated collection endpoints, initialize the page parameter, append each page through the reusable item model, and stop when the response has no data or `current_page * per_page >= total`. Check for a nil `JSON200` before dereferencing it and preserve the API's null/empty behavior.

### 3. Register in the provider

Add the constructor to the `DataSources` slice in `internal/provider/provider.go`:

```go
func (p *hostingerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        // ... existing entries ...
        New<PascalName>DataSource,
    }
}
```

### 4. Create the acceptance test

Create `internal/provider/<name>_data_source_test.go` following the [test template](./assets/data_source_test.go.tmpl).

Key rules:
- Use `resource.Test` with the repository provider factory; enable the acceptance run externally with `TF_ACC=1`
- Reference the data source as `data.hostinger_<name>.test`
- Use `statecheck.ExpectKnownValue` with typed `knownvalue.*` checkers
- Add a `const testAcc<PascalName>DataSourceConfig` HCL string with the minimal filter config, or an empty data block when the endpoint has no filters
- For local MockServer-backed tests, clear/freeze the server in a resource-specific precheck, match generated `operationId` values, and provide one expectation for every paginated request

### 5. Generate documentation

Run `make generate` to regenerate provider docs, then verify `docs/data-sources/<name>.md` was created.

### 6. Build & test

```bash
make build          # verifies the code compiles
make testacc        # runs acceptance tests (requires HOSTINGER_HOST + HOSTINGER_API_TOKEN)
```

## File Map

| File | Purpose |
|------|---------|
| `internal/provider/<name>_model.go` | Reusable API response model and `Merge` mappings, when needed |
| `internal/provider/<name>_data_source.go` | Implementation |
| `internal/provider/<name>_data_source_test.go` | Acceptance tests |
| `docs/data-sources/<name>.md` | Generated docs (via `make generate`) |

## Reference

- [Implementation template](./assets/data_source.go.tmpl)
- [Test template](./assets/data_source_test.go.tmpl)
- [Existing example](./references/billing_catalogs_example.md)
- [Schema type mappings](./references/schema_types.md)
