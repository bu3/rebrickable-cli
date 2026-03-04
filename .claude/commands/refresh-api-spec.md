Download the latest Rebrickable OpenAPI spec and update the requirements doc.

## Steps

1. Download the spec from the live endpoint and overwrite `openapi.spec.json`:
   ```bash
   curl -s "https://rebrickable.com/api/v3/swagger/?format=openapi" -o openapi.spec.json
   ```

2. Pretty-print it in place:
   ```bash
   python3 -m json.tool openapi.spec.json > /tmp/openapi_pretty.json && mv /tmp/openapi_pretty.json openapi.spec.json
   ```

3. Compare the updated spec against the currently implemented methods in `cli/cmd/api/api.go` and rewrite `requirements.md` to reflect:
   - All paths and HTTP methods in the new spec
   - Which are already implemented (mark as done or exclude)
   - Which are not yet implemented, grouped by resource
   - Any implementation notes relevant to the gaps (auth requirements, composite keys, destructive operations, txtar suitability)
