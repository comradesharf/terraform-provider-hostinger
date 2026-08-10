---
name: create-terraform-action
description: Create a Hostinger Terraform Plugin Framework action, its provider registration, acceptance test, and generated documentation. Use when implementing an API-backed action that invokes an operation immediately and may wait for an asynchronous VPS action to finish.
---

# Create Terraform Action

Implement actions in `internal/provider` using the pattern established by [`vps_firewall_activate_action.go`](../../../internal/provider/vps_firewall_activate_action.go) and [`vps_firewall_activate_action_test.go`](../../../internal/provider/vps_firewall_activate_action_test.go).

## Workflow

1. Identify the Terraform action name, required inputs, generated client operation, success status, response type, and whether the API returns an asynchronous action ID.
2. Inspect the generated client in `internal/client/client.gen.go` and a nearby action or resource that uses the same endpoint family. Confirm the action-detail operation and generated state constants before writing code.
3. Create `internal/provider/<name>_action.go`. Define the action, constructor, model, client configuration, metadata, schema, and invoke implementation. Register the constructor in the provider's `Actions` method in `internal/provider/provider.go`.
4. Create `internal/provider/<name>_action_test.go` as a MockServer-backed acceptance test. Gate it for Terraform 1.14 or newer, configure the action through Terraform's `action` and `action_trigger` blocks, and verify the exact request sequence.
5. Add or regenerate `docs/actions/<name>.md` using the repository documentation generation command. Do not hand-edit generated action documentation.
6. Run `gofmt`, focused provider tests, and the build. Run acceptance tests only when the required MockServer/API test environment is available.

## Implementation Requirements

- Use `action.Action` and `action.ActionWithConfigure` compile-time assertions and return `action.Action` from `New<PascalName>Action`.
- Keep the configured client as `*client.ClientWithResponses`; accept only that type in `Configure` and report an unexpected type through diagnostics.
- Set `Metadata` to `req.ProviderTypeName + "_<name>"`.
- Model every schema field with a matching `tfsdk` field. Include `timeouts.Value` and `timeouts.Attributes(ctx)` when the action performs network work.
- Use `schema.Int64Attribute`, `schema.StringAttribute`, `schema.BoolAttribute`, or the appropriate framework type with clear descriptions. Required identifiers must be validated before the API call.
- Resolve the invoke timeout with `config.Timeouts.Invoke(ctx, 20*time.Minute)`, wrap the context with `context.WithTimeout`, and always defer cancellation.
- Call generated `*WithResponse` client methods. Check transport errors, the expected HTTP status, a non-nil success body, and a non-nil action ID before polling.
- Make `wait_for_action` optional and treat null or unknown as `true`. When waiting, poll the generated action-details operation until the generated success or error state. Check transport status, response body, and state before dereferencing.
- On an in-progress state, honor context cancellation, use the established bounded polling interval, and send an `action.InvokeProgressEvent` describing the wait. Return diagnostics with the current action's name; do not copy unrelated error text from another action.
- Use `tflog.SetField` for non-sensitive identifiers. Do not log credentials or sensitive inputs.
- Keep implementation inline and follow the repository instruction not to create helper functions.

## Test Requirements

- Use `resource.Test`, `testAccProtoV6ProviderFactories`, `testAccPreCheck`, and a resource-specific MockServer precheck.
- Skip Terraform versions below `1.14.0`, because action blocks and `action_trigger` require that feature set.
- Clear the MockServer and freeze its clock in the precheck. Seed expectations using the generated OpenAPI `operationId`, not an ad hoc endpoint name.
- Configure one mutation response containing a stable action ID, then provide one action-details response per poll. Verify the exact POST/GET sequence with `VerifySequence` in `PostApplyFunc`.
- Cover the default wait behavior or explicit `wait_for_action = true`; add a separate no-wait case when the action exposes meaningful behavior for `false`.
- Prefer focused unit tests for deterministic validation and diagnostic branches. Acceptance tests should exercise Terraform configuration, provider wiring, generated client calls, and polling behavior.

## Validation

```bash
gofmt -w internal/provider/<name>_action.go internal/provider/<name>_action_test.go
go test ./internal/provider
make build
make generate
```

Read [action-implementation-reference.md](references/action-implementation-reference.md) when implementing invoke behavior and [action-test-reference.md](references/action-test-reference.md) when creating MockServer expectations.
