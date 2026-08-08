# List Resource Tests

## Acceptance Test Pattern

Follow the existing provider test conventions:

1. Define `TestAcc<Resource>List` with `resource.Test`.
2. Call the shared `testAccPreCheck` and a resource-specific precheck.
3. Configure `ProtoV6ProviderFactories` with `testAccProtoV6ProviderFactories`.
4. Add an initial empty configuration step, then a `list` block with `include_resource = true`, `Query: true`, and `GenerateConfig: true`.
5. In the precheck, clear and freeze the mock server, then register one expectation per API page. Set `times.remainingTimes` so each expectation is consumed once.
6. Use `querycheck.ExpectResourceKnownValues` and `queryfilter.ByResourceIdentity` to assert each result. Verify all identity components for compound identities and meaningful resource fields such as `name`, `key`, or nested rule values.

Example assertion shape:

```go
querycheck.ExpectResourceKnownValues(
    "hostinger_vps_public_key.test",
    queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
        "id": knownvalue.Int64Exact(325),
    }),
    []querycheck.KnownValueCheck{
        {tfjsonpath.New("name"), knownvalue.StringExact("My public key 1")},
    },
)
```

Force pagination by returning one item per page and setting `total` to the number of expected items. Include a final page whose metadata proves the loop stops without an unnecessary request.

## Unit/Error Tests

Add unit tests when the implementation introduces behavior that acceptance tests do not isolate. The minimum cases are:

- client method returns an error;
- response status is not `http.StatusOK`;
- successful response has a nil body;
- successful body has nil or empty data;
- pagination metadata is missing or indicates completion;
- `push` returns false and no further items/pages are emitted;
- nested response emits one result per nested object with complete identity values.

Use the narrowest test seam available in the provider. Do not weaken acceptance assertions merely to make generated configuration pass; identity and resource fields are the contract being tested.
