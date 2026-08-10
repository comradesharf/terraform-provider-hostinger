# Action Acceptance-Test Reference

The canonical example is [`internal/provider/vps_firewall_activate_action_test.go`](../../../../internal/provider/vps_firewall_activate_action_test.go).

## Test shape

- The test uses `resource.Test` and `testAccProtoV6ProviderFactories`.
- `TerraformVersionChecks` skips versions below `tfversion.Version1_14_0`.
- `PreCheck` calls the common provider precheck and an action-specific precheck.
- The action-specific precheck clears MockServer, freezes the clock, and installs JSON expectations through the MockServer REST endpoint.
- The config creates a `terraform_data` resource with an `action_trigger`, then declares the provider action in an `action` block.
- `PostApplyFunc` verifies the mutation POST followed by the expected number of action-details GET requests.

## Expectations

Use the generated OpenAPI `operationId` for each expectation. The mutation response must contain the same action ID used in the request path for every poll. Set `remainingTimes` to the number of calls the implementation will make, and provide response states that demonstrate the polling transition, for example `created`, `created`, then `success`.

Add cases for non-success mutation responses, missing action IDs, action-detail failures, action error states, and `wait_for_action = false` when those branches are part of the action's behavior. Keep expected paths and IDs deterministic so `VerifySequence` catches incorrect argument wiring.
