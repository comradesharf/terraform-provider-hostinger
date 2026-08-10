# Repository Instructions

- Do not create helper functions. Inline implementation at the call site unless an existing shared function is being reused.

## Testing

- Run acceptance tests with `TF_ACC=1`, `HOSTINGER_HOST`, and `HOSTINGER_API_TOKEN`; plain `go test` does not run the acceptance tests correctly.
- To run a focused acceptance test, use `TF_ACC=1 HOSTINGER_HOST=http://localhost:1234 HOSTINGER_API_TOKEN=123 go test -v -cover -timeout 120m ./... -count 1 -run hostinger_vps_firewall_deactivate`.
