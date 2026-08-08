# List Resource Reference

Canonical local examples:

- `internal/provider/vps_public_key_list.go`
- `internal/provider/vps_firewall_list.go`
- `internal/provider/vps_firewall_rule_list.go`
- `internal/provider/vps_public_key_list_test.go`

Use this shape for a normal paginated endpoint, substituting generated types and resource-specific names:

```go
func (l *ResourceList) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
    var data ResourceListModel
    diags := req.Config.Get(ctx, &data)
    if diags.HasError() {
        stream.Results = list.ListResultsStreamDiagnostics(diags)
        return
    }

    stream.Results = func(push func(list.ListResult) bool) {
        params := &client.GetItemsParams{}
        page := 1
        for {
            params.Page = &page
            ctx = tflog.SetField(ctx, "page", params.Page)
            response, err := l.client.GetItemsWithResponse(ctx, params)
            if err != nil {
                result := req.NewListResult(ctx)
                result.Diagnostics.AddError("Unable to Read Resources", fmt.Sprintf("Got error: %s", err))
                push(result)
                return
            }
            if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
                // Emit separate status and nil-body diagnostics using repository wording.
                return
            }
            if response.JSON200.Data == nil || len(*response.JSON200.Data) == 0 {
                break
            }

            for _, item := range *response.JSON200.Data {
                var model ResourceResourceModel
                model.Merge(item)
                result := req.NewListResult(ctx)
                result.DisplayName = model.Name.ValueString()
                identity := ResourceIdentity{ID: model.ID}
                result.Diagnostics.Append(result.Identity.Set(ctx, &identity)...)
                if req.IncludeResource {
                    model.Timeouts = resourcetimeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
                        "create": types.StringType, "read": types.StringType,
                        "update": types.StringType, "delete": types.StringType,
                    })}
                    result.Diagnostics.Append(result.Resource.Set(ctx, &model)...)
                }
                if !push(result) {
                    return
                }
            }

            meta := response.JSON200.Meta
            if meta == nil || meta.CurrentPage == nil || meta.PerPage == nil || meta.Total == nil {
                break
            }
            if (*meta.CurrentPage)*(*meta.PerPage) >= *meta.Total {
                break
            }
            page++
        }
    }
}
```

For nested objects, keep the pagination loop around the top-level response and emit results inside the nested loop. Populate every identity field and use a stable compound display name, as `VPSFirewallRuleList` does.

## Test Pairing

Every new list implementation should have a matching `*_list_test.go` file. The acceptance test should configure the list with `include_resource = true`, mock multiple pages with `meta.total`, `meta.current_page`, and `meta.per_page`, then use `querycheck.ExpectResourceKnownValues` with `queryfilter.ByResourceIdentity` to verify identities and resource attributes. See `internal/provider/vps_public_key_list_test.go` and `internal/provider/vps_firewall_list_test.go`.

For a nested list, assert the complete compound identity and display-relevant resource attributes for at least two nested items. Ensure the mock response contains the exact generated `operationId` and enough nested data to exercise the loop.

When adding unit tests for error paths, verify that each failure produces a diagnostic and terminates streaming: API call errors, unexpected status codes, nil `JSON200`, nil data, and incomplete pagination metadata. Keep those tests focused on behavior rather than reproducing the full acceptance setup.
