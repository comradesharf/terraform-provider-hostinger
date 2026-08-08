---
name: create-terraform-list-resource
description: Create or extend a Hostinger Terraform Plugin Framework list resource and its tests in this provider. Use when implementing a new paginated list resource, adding list test coverage for an existing resource, or reviewing registration, identity, pagination, diagnostics, and include_resource behavior.
---

# Create Terraform List Resource

## Overview

Implement list resources by following the established `internal/provider/*_list.go` pattern. A list resource enumerates API objects, exposes each object's identity and display name, and optionally returns the corresponding resource model when `include_resource` is enabled.

## Workflow

1. Inspect the singular resource, model, identity type, generated client methods, provider registration, and a neighboring list implementation. Prefer the closest resource by API response shape and identity fields.
2. Confirm the generated endpoint's pagination parameters and response metadata. Use the generated `*Params` type and `*WithResponse` method; do not hand-build HTTP requests.
3. Create `internal/provider/<resource>_list.go` with compile-time interfaces, a constructor, client configuration, matching metadata, an empty config schema unless filters are needed, and a lazy `List` stream.
4. In `List`, read config diagnostics, iterate pages, log the page, check transport/status/nil-body/nil-data failures, merge each item into the singular resource model, set display name and identity, and set the resource only when `req.IncludeResource` is true.
5. Register the constructor in `hostingerProvider.ListResources` in `internal/provider/provider.go`.
6. Create `internal/provider/<resource>_list_test.go` alongside the implementation. Cover at least two pages, generated configuration, identity filtering, and selected resource attributes with `include_resource = true`. Match mock-server expectations to the generated operation ID.
7. Add focused error-path tests when the list method has custom branching, especially non-200 responses, transport errors, nil response bodies, incomplete metadata, and nested items.
8. Run `gofmt` and focused provider tests; run the broader test target when practical.

## Implementation Rules

- Stream lazily through `stream.Results` and stop immediately when `push` returns false.
- Set the page field in the logging context before each API call.
- Check `StatusCode() == http.StatusOK`, `JSON200 != nil`, and `JSON200.Data != nil` before dereferencing.
- Treat incomplete pagination metadata as a stop condition, matching existing provider behavior.
- Reuse the singular resource model's `Merge` method and identity type. Do not duplicate API-to-Terraform mapping.
- Before `result.Resource.Set`, assign the resource timeout object to null attributes for create/read/update/delete, as in existing list implementations.
- For nested objects such as firewall rules, emit one result per nested object, populate every identity field, and derive a compound display name.
- Tests should verify every emitted identity field and at least one resource attribute, not only the number of results.
- Use acceptance tests for Terraform query/config generation and API pagination; use unit tests for deterministic diagnostics and edge cases that do not require Terraform protocol execution.

Read [list-resource-reference.md](references/list-resource-reference.md) when writing or reviewing the concrete `List` method, and [list-resource-tests.md](references/list-resource-tests.md) when creating tests.
