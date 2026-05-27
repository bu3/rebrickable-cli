# LEGO Colors API — Design

**Date:** 2026-05-27

## Summary

Implement the two LEGO colors catalog endpoints:

- `GET /api/v3/lego/colors/` — list all colors
- `GET /api/v3/lego/colors/{id}/` — get a single color by ID

These are read-only, unauthenticated endpoints (API key only, no user token), consistent with all other LEGO catalog endpoints.

---

## Data Model

**File:** `cli/cmd/api/api.go`

Expand the existing `PartColor` struct with `RGB` and `IsTrans` fields:

```go
type PartColor struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    RGB     string `json:"rgb"`
    IsTrans bool   `json:"is_trans"`
}
```

Add a `ColorsResponse` type for paginated list results:

```go
type ColorsResponse struct {
    Count    int         `json:"count"`
    Next     string      `json:"next"`
    Previous string      `json:"previous"`
    Results  []PartColor `json:"results"`
}
```

`PartColor` continues to serve as the embedded color in `SetPart`. The new fields will be zero-valued when the API response omits them (e.g., in set parts responses), which is safe.

---

## API Client

**File:** `cli/cmd/api/lego.go`

Two new methods on `*Client`:

- `GetLegoColors() (*ColorsResponse, error)` — uses `fetchAllPages[PartColor]` against `/lego/colors/`
- `GetLegoColor(id string) (*PartColor, error)` — direct GET to `/lego/colors/{id}/`, returns error on non-200

---

## CLI Commands

**File:** `cli/cmd/lego_colors.go` (new file)

New `lego colors` subcommand registered on `legoCmd`, with two sub-subcommands:

| Command | Flags | Action |
|---------|-------|--------|
| `lego colors list` | none | calls `GetLegoColors()`, prints JSON |
| `lego colors get` | `--id` (required) | calls `GetLegoColor(id)`, prints JSON |

Registration follows the same `legoColorsCommands()` + `init()` pattern as `lego_parts.go` and `lego_sets.go`. Uses `newLegoAPIClient(cmd)` and `printJSON()`.

---

## Tests

### Unit tests

**File:** `cli/cmd/api/api_test.go`

Two table-driven tests using `httptest.NewServer` + `newClientWithBaseURL`:

- `TestGetLegoColors` — cases: returns colors (200), server error (500)
- `TestGetLegoColor` — cases: returns color (200), not found (404), server error (500)

### Integration test (txtar)

**File:** `testdata/lego_colors.txtar` (new file)

```
# List all colors
exec cli lego colors list
stdout '"count":'
! stderr .

# Get a specific color by stable ID
exec cli lego colors get --id 0
stdout '"name": "Black"'
! stderr .
```

Color ID 0 (Black) is a stable, well-known identifier in the Rebrickable dataset.

---

## Files Changed

| File | Change |
|------|--------|
| `cli/cmd/api/api.go` | Expand `PartColor`; add `ColorsResponse` |
| `cli/cmd/api/lego.go` | Add `GetLegoColors`, `GetLegoColor` |
| `cli/cmd/api/api_test.go` | Add `TestGetLegoColors`, `TestGetLegoColor` |
| `cli/cmd/lego_colors.go` | New file — Cobra commands |
| `testdata/lego_colors.txtar` | New integration test |
| `requirements.md` | Mark Colors endpoints as implemented |
