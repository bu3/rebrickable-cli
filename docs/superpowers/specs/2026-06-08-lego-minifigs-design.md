# LEGO Minifigs API — Design

**Date:** 2026-06-08

## Summary

Implement the four LEGO minifigs catalog endpoints:

- `GET /api/v3/lego/minifigs/` — list all minifigs
- `GET /api/v3/lego/minifigs/{set_num}/` — get a single minifig by ID
- `GET /api/v3/lego/minifigs/{set_num}/parts/` — list parts of a minifig
- `GET /api/v3/lego/minifigs/{set_num}/sets/` — list sets containing a minifig

These are read-only, unauthenticated endpoints (API key only, no user token). Minifig IDs follow the pattern `fig-000001`.

---

## Data Model

**File:** `cli/cmd/api/api.go`

Add a new `Minifig` struct with only the fields minifigs actually have (distinct from the existing `Set` struct, which has `Year` and `ThemeID` that minifigs don't):

```go
type Minifig struct {
	SetNum         string `json:"set_num"`
	Name           string `json:"name"`
	NumParts       int    `json:"num_parts"`
	SetImgURL      string `json:"set_img_url"`
	SetURL         string `json:"set_url"`
	LastModifiedDt string `json:"last_modified_dt"`
}
```

Add a `MinifigsResponse` type for paginated list results:

```go
type MinifigsResponse struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Minifig `json:"results"`
}
```

The `/parts/` and `/sets/` sub-endpoints reuse existing types: `SetPartsResponse` and `LegoSetsResponse` respectively — no new types needed.

---

## API Client

**File:** `cli/cmd/api/lego.go`

Four new methods on `*Client`:

| Method | Path | Return type |
|--------|------|-------------|
| `GetLegoMinifigs() (*MinifigsResponse, error)` | `/lego/minifigs/` | `fetchAllPages[Minifig]` |
| `GetLegoMinifig(figNum string) (*Minifig, error)` | `/lego/minifigs/{figNum}/` | direct GET |
| `GetLegoMinifigParts(figNum string) (*SetPartsResponse, error)` | `/lego/minifigs/{figNum}/parts/` | `fetchAllPages[SetPart]` |
| `GetLegoMinifigSets(figNum string) (*LegoSetsResponse, error)` | `/lego/minifigs/{figNum}/sets/` | `fetchAllPages[Set]` |

---

## CLI Commands

**File:** `cli/cmd/lego_minifigs.go` (new file)

New `lego minifigs` subcommand registered on `legoCmd`, with four sub-subcommands:

| Command | Flags | Action |
|---------|-------|--------|
| `lego minifigs list` | none | calls `GetLegoMinifigs()`, prints JSON |
| `lego minifigs get` | `--fig_num` / `-f` (required) | calls `GetLegoMinifig(figNum)`, prints JSON |
| `lego minifigs parts` | `--fig_num` / `-f` (required) | calls `GetLegoMinifigParts(figNum)`, prints JSON |
| `lego minifigs sets` | `--fig_num` / `-f` (required) | calls `GetLegoMinifigSets(figNum)`, prints JSON |

Registration follows the same `legoMinifigsCommands()` + `init()` pattern as `lego_colors.go`, `lego_elements.go`, etc. Uses `newLegoAPIClient(cmd)` and `printJSON()`.

---

## Tests

### Unit tests

**File:** `cli/cmd/api/api_test.go`

Four table-driven tests using `httptest.NewServer` + `newClientWithBaseURL`:

- `TestGetLegoMinifigs` — cases: returns minifigs (200), server error (500)
- `TestGetLegoMinifig` — cases: returns minifig (200), not found (404), server error (500)
- `TestGetLegoMinifigParts` — cases: returns parts (200), server error (500)
- `TestGetLegoMinifigSets` — cases: returns sets (200), server error (500)

### Integration test (txtar)

**File:** `testdata/lego_minifigs.txtar` (new file)

```
# List all minifigs
exec cli lego minifigs list
stdout '"count":'
! stderr .

# Get a specific minifig by stable ID
exec cli lego minifigs get --fig_num fig-000001
stdout '"set_num": "fig-000001"'
! stderr .
```

`fig-000001` is a stable, well-known identifier in the Rebrickable dataset.

---

## Files Changed

| File | Change |
|------|--------|
| `cli/cmd/api/api.go` | Add `Minifig`, `MinifigsResponse` |
| `cli/cmd/api/lego.go` | Add 4 client methods |
| `cli/cmd/api/api_test.go` | Add 4 unit tests |
| `cli/cmd/lego_minifigs.go` | New file — Cobra commands |
| `cli/cmd/BUILD.bazel` | Add `lego_minifigs.go` to srcs |
| `testdata/lego_minifigs.txtar` | New integration test |
| `requirements.md` | Mark Minifigs endpoints as implemented |
