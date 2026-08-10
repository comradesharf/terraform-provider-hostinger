# Action Implementation Reference

The canonical example is [`internal/provider/vps_firewall_activate_action.go`](../../../../internal/provider/vps_firewall_activate_action.go).

## Structure

The file contains, in order:

1. Imports for context, formatting, HTTP status constants, timeouts, the generated client, framework action/schema/types, and logging.
2. Compile-time assertions for `action.Action` and `action.ActionWithConfigure`.
3. A constructor, action struct with the configured client, and an action model with `timeouts`, API identifiers, and `wait_for_action`.
4. `Configure`, `Schema`, `Metadata`, and `Invoke` methods.

## Invoke sequence

Use this order:

1. Decode `req.Config` into the model and stop on diagnostics.
2. Reject null or unknown required IDs before resolving a timeout or calling the API.
3. Resolve `Invoke` timeout with the action default of `20*time.Minute`; wrap the context and defer cancellation.
4. Log the relevant identifier, call the generated mutation `WithResponse` method, and check transport error, expected status, success body, and action ID.
5. If waiting is enabled, call the generated action-details method with the virtual machine ID and captured action ID. Handle transport errors, non-200 responses, nil bodies, and missing states.
6. Exit on the generated success state, emit an action-specific diagnostic on the generated error state, or continue after a bounded delay while sending progress.

The sample's null/unknown handling means the default is to wait. Preserve this behavior for asynchronous actions unless the API contract explicitly differs. Prefer generated enum constants over raw state strings.
