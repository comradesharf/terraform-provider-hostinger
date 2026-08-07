---
name: resource-creation
description: 'Create a new Terraform resource for the Hostinger provider. Use when adding a managed resource, implementing CRUD against a generated Hostinger API client, scaffolding resource tests and docs, or wiring a resource into the provider.'
---

# Resource Creation

Creates a managed Terraform resource in this provider using the existing `terraform-plugin-framework` and generated-client patterns.

## When to Use

- Adding a new API-backed Terraform resource
- Scaffolding a resource implementation, acceptance test, and generated docs
- Extending an existing resource with lifecycle behavior or import support

## Prerequisites

- Resource name in snake_case and its Terraform type name
- Generated client operations and response types in `internal/client/client.gen.go`
- Required, optional, computed, and replacement semantics for each API field
- Import ID format, if the API needs more than the Terraform `id`
- A reusable `internal/provider/<name>_model.go` model, or a decision to create one with the [model-creation skill](../model-creation/SKILL.md)

## Procedure

### 1. Inspect the API and neighboring code

Search `internal/client/client.gen.go` for the create, read, update, and delete `WithResponse` methods. Record their parameter types, request bodies, success status, not-found behavior, and `JSON200` response type. Inspect the closest resource and model in `internal/provider` and reuse established conversions and helpers.

For a resource related to VPS firewalls, use [`vps_firewall_rule_resource.go`](../../../internal/provider/vps_firewall_rule_resource.go) as the lifecycle reference. Its rule model is embedded in the resource model, while `firewall_id` is kept as resource-specific state.

### 2. Define the model

Create or reuse `internal/provider/<name>_model.go` for API response fields and their `Merge` method. Keep `timeouts.Value`, lookup/parent IDs, and other lifecycle-only state in the resource model. Use framework value types and preserve API nulls; do not put schema declarations or API calls in the model file.

Use the [model-creation skill](../model-creation/SKILL.md) for nested response mapping, generated enums, IP values, timestamps, and collections.

### 3. Implement the resource

Create `internal/provider/<name>_resource.go` with:

- `resource.Resource`, `resource.ResourceWithConfigure`, and `resource.ResourceWithImportState` assertions when applicable
- `New<PascalName>Resource() resource.Resource`
- A client field of type `*client.ClientWithResponses`
- A resource model embedding or containing the reusable API model and a `Timeouts timeouts.Value` field tagged `tfsdk:"timeouts"`
- `Metadata` setting `req.ProviderTypeName + "_<name>"`
- `Configure` accepting only `*client.ClientWithResponses` and reporting an unexpected type through diagnostics
- `Schema` with descriptions, framework types, validators, and plan modifiers
- `"timeouts": timeouts.AttributesAll(ctx)` and a default operation timeout of `20*time.Minute`

Schema rules:

- API-generated IDs are `Computed` and normally use `int64planmodifier.UseStateForUnknown()`
- Parent IDs or fields that cannot be updated use `RequiresReplace()`
- Create/update inputs are `Required` or `Optional` according to the API contract
- API response-only fields are `Computed`
- Enum strings use `stringvalidator.OneOf(...)`; length and other constraints should use framework validators
- Schema custom types must match the model types, such as RFC3339 or network types

Lifecycle rules:

- `Create` reads `req.Config`, resolves `state.Timeouts.Create(ctx, 20*time.Minute)`, builds the generated request body, calls the client, checks transport errors, expected status, and nil `JSON200`, merges the response, and writes state
- `Read` reads state, validates IDs before calling the API, resolves the read timeout, removes the resource on a documented not-found response, merges the matching API object, and writes state
- `Update` reads `req.Plan`, validates IDs, resolves the update timeout, sends only mutable fields, merges the response, and writes state. Fields marked `RequiresReplace` must not be updated
- `Delete` reads state, validates IDs, resolves the delete timeout, calls the delete endpoint, and accepts the documented success status
- Use `context.WithTimeout` and `defer cancel()` for every API operation
- Use `tflog.SetField` for useful non-sensitive identifiers; never log credentials or sensitive values

For collection responses, locate the item by its API ID before merging. If the item or parent is absent, call `resp.State.RemoveResource(ctx)` rather than persisting stale state.

### 4. Implement import

Add `ImportState` when the resource requires an import workflow. For a single numeric ID, use the repository helper or framework passthrough pattern. For a composite ID, parse and validate every component before setting state attributes. The firewall-rule canonical format is `firewall_id/rule_id`; reject empty, malformed, or non-numeric parts with a diagnostic and stop before setting partial state.

### 5. Register the resource

Add `New<PascalName>Resource` to the `Resources` slice in `internal/provider/provider.go`:

```go
func (p *hostingerProvider) Resources(ctx context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        // existing resources
        New<PascalName>Resource,
    }
}
```

### 6. Add tests

Create `internal/provider/<name>_resource_test.go` following the existing acceptance-test style:

- Use `resource.Test` and `testAccProtoV6ProviderFactories`
- Gate setup with `testAccPreCheck` and a resource-specific mock/API precheck
- Cover create, read/refresh, update, replacement behavior, import, and destroy where supported
- Assert IDs, parent IDs, inputs, and computed response fields with typed `knownvalue.*` checks
- For mock-server tests, match generated operation IDs and provide enough response repetitions for every refresh/read
- Add invalid import IDs and validation tests when parsing or schema validation has meaningful branches

Start with `go test ./internal/provider` for unit/package tests. Run `TF_ACC=1 ... go test` only when acceptance infrastructure and credentials are available.

### 7. Generate documentation and verify

Run:

```bash
gofmt -w internal/provider/<name>_model.go internal/provider/<name>_resource.go internal/provider/<name>_resource_test.go
make generate
go test ./internal/provider
make build
```

Confirm `docs/resources/<name>.md` reflects the final schema. Do not hand-edit generated provider documentation.

## File Map

| File | Purpose |
|---|---|
| `internal/provider/<name>_model.go` | Reusable API response model and `Merge` mapping |
| `internal/provider/<name>_resource.go` | Schema, configuration, CRUD, timeouts, and import |
| `internal/provider/<name>_resource_test.go` | Acceptance and resource behavior tests |
| `internal/provider/provider.go` | Resource registration |
| `docs/resources/<name>.md` | Generated resource documentation |

## References

- Canonical resource: [`vps_firewall_rule_resource.go`](../../../internal/provider/vps_firewall_rule_resource.go)
- Canonical resource model: [`vps_firewall_rule_model.go`](../../../internal/provider/vps_firewall_rule_model.go)
- [Model creation skill](../model-creation/SKILL.md)
