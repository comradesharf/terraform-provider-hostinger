---
name: data-source-creation
description: 'Create a new Terraform data source for the Hostinger provider. Use when adding a new data source, implementing a read-only data source, scaffolding a data source, or wiring up a new API endpoint as a data source.'
argument-hint: '<data_source_name> e.g. billing_subscriptions'
---

# Data Source Creation

Creates a new Terraform data source for the Hostinger provider following the established `terraform-plugin-framework` patterns.

## When to Use

- Adding a new read-only data source backed by the Hostinger API
- Scaffolding a data source implementation + test + docs

## Prerequisites

- Identify the **data source name** (snake_case, e.g. `billing_subscriptions`)
- Identify the **client method** in `internal/client/client.gen.go` (e.g. `BillingGetSubscriptionListV1WithResponse`)
- Know the **response struct** and its fields to map to the Terraform schema

## Procedure

### 1. Explore the client

Search `internal/client/client.gen.go` for the relevant API method and its response/params types to understand:
- Method signature (`WithResponse` variant)
- `Params` struct (optional query filters → become `Optional` schema attributes)
- Response struct (`JSON200`) and its nested types → become `Computed` schema attributes

### 2. Create the implementation file

Create `internal/provider/data_source_<name>.go` following the [implementation template](./assets/data_source.go.tmpl).

Key rules:
- Copyright header: `// Copyright IBM Corp. 2021, 2025\n// SPDX-License-Identifier: MPL-2.0`
- Package: `package provider`
- Struct name pattern: `DataSource<PascalName>` (e.g. `DataSourceBillingCatalogs`)
- Constructor: `NewDataSource<PascalName>() datasource.DataSource`
- `Metadata`: set `resp.TypeName = req.ProviderTypeName + "_<name>"`
- `Configure`: extract `*client.ClientWithResponses` from `req.ProviderData`
- `Schema`: filter params → `Optional`, response fields → `Computed`
- `Read`: call the `WithResponse` method, check `StatusCode() != http.StatusOK`, map `JSON200` fields to model structs
- Use `types.StringPointerValue`, `types.Int32Value`, `types.MapValueMust`, etc. for type conversions
- Use `tflog.SetField` to attach filter values to context before the API call

### 3. Register in the provider

Add the constructor to the `DataSources` slice in `internal/provider/provider.go`:

```go
func (p *hostingerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        // ... existing entries ...
        NewDataSource<PascalName>,
    }
}
```

### 4. Create the acceptance test

Create `internal/provider/data_source_<name>_test.go` following the [test template](./assets/data_source_test.go.tmpl).

Key rules:
- Use `resource.Test` with `TF_ACC=1`
- Reference the data source as `data.hostinger_<name>.test`
- Use `statecheck.ExpectKnownValue` with typed `knownvalue.*` checkers
- Add a `const testAcc<PascalName>Config` HCL string with the minimal filter config

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
| `internal/provider/data_source_<name>.go` | Implementation |
| `internal/provider/data_source_<name>_test.go` | Acceptance tests |
| `docs/data-sources/<name>.md` | Generated docs (via `make generate`) |

## Reference

- [Implementation template](./assets/data_source.go.tmpl)
- [Test template](./assets/data_source_test.go.tmpl)
- [Existing example](./references/billing_catalogs_example.md)
- [Schema type mappings](./references/schema_types.md)
