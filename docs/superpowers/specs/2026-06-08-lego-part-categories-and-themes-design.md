# LEGO Part Categories and Themes API — Design

**Date:** 2026-06-08

## Summary

Implement the four remaining LEGO catalog endpoints:

- `GET /api/v3/lego/part_categories/` — list all part categories
- `GET /api/v3/lego/part_categories/{id}/` — get a single part category by ID
- `GET /api/v3/lego/themes/` — list all themes
- `GET /api/v3/lego/themes/{id}/` — get a single theme by ID

Both are read-only, unauthenticated endpoints (API key only, no user token). IDs are integers. This completes the full LEGO catalog (20 of 20 endpoints).

---

## Data Model

**File:** `cli/cmd/api/api.go`

```go
type PartCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	PartCount int    `json:"part_count"`
}

type PartCategoriesResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []PartCategory `json:"results"`
}

type Theme struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

type ThemesResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous string  `json:"previous"`
	Results  []Theme `json:"results"`
}
```

`Theme.ParentID` is a pointer because the field is nullable — top-level themes have `parent_id: null`.

---

## API Client

**File:** `cli/cmd/api/lego.go`

Four new methods on `*Client`:

| Method | Path | Return type |
|--------|------|-------------|
| `GetLegoPartCategories() (*PartCategoriesResponse, error)` | `/lego/part_categories/` | `fetchAllPages[PartCategory]` |
| `GetLegoPartCategory(id int) (*PartCategory, error)` | `/lego/part_categories/{id}/` | direct GET |
| `GetLegoThemes() (*ThemesResponse, error)` | `/lego/themes/` | `fetchAllPages[Theme]` |
| `GetLegoTheme(id int) (*Theme, error)` | `/lego/themes/{id}/` | direct GET |

Single-resource methods return error on non-200, following the same pattern as `GetLegoColor` and `GetLegoElement`.

---

## CLI Commands

**File:** `cli/cmd/lego_part_categories.go` (new)

| Command | Flags | Action |
|---------|-------|--------|
| `lego part-categories list` | none | calls `GetLegoPartCategories()`, prints JSON |
| `lego part-categories get` | `--id` / `-i` (required) | calls `GetLegoPartCategory(id)`, prints JSON |

**File:** `cli/cmd/lego_themes.go` (new)

| Command | Flags | Action |
|---------|-------|--------|
| `lego themes list` | none | calls `GetLegoThemes()`, prints JSON |
| `lego themes get` | `--id` / `-i` (required) | calls `GetLegoTheme(id)`, prints JSON |

Both follow the `legoXxxCommands()` + `init()` pattern. Flags use integer IDs — the `--id` flag is `int` type (not string), matching `GetLegoPartCategory(id int)` and `GetLegoTheme(id int)`.

---

## Tests

### Unit tests

**File:** `cli/cmd/api/api_test.go`

Four table-driven tests using `httptest.NewServer` + `newClientWithBaseURL`:

- `TestGetLegoPartCategories` — cases: returns categories (200), server error (500)
- `TestGetLegoPartCategory` — cases: returns category (200), not found (404), server error (500)
- `TestGetLegoThemes` — cases: returns themes (200), server error (500)
- `TestGetLegoTheme` — cases: returns theme (200), not found (404), server error (500)

### Integration tests (txtar)

**File:** `testdata/lego_part_categories.txtar` (new)

```
# List all part categories
exec cli lego part-categories list
stdout '"count":'
! stderr .

# Get a specific part category by stable ID
exec cli lego part-categories get --id 1
stdout '"id": 1'
! stderr .
```

**File:** `testdata/lego_themes.txtar` (new)

```
# List all themes
exec cli lego themes list
stdout '"count":'
! stderr .

# Get a specific theme by stable ID
exec cli lego themes get --id 1
stdout '"id": 1'
! stderr .
```

ID `1` is stable for both resources (verified against live API: Part Category 1 = "Baseplates", Theme 1 = "Technic").

---

## Files Changed

| File | Change |
|------|--------|
| `cli/cmd/api/api.go` | Add `PartCategory`, `PartCategoriesResponse`, `Theme`, `ThemesResponse` |
| `cli/cmd/api/lego.go` | Add 4 client methods |
| `cli/cmd/api/api_test.go` | Add 4 unit tests |
| `cli/cmd/lego_part_categories.go` | New file — Cobra commands |
| `cli/cmd/lego_themes.go` | New file — Cobra commands |
| `cli/cmd/BUILD.bazel` | Add both new files to srcs |
| `testdata/lego_part_categories.txtar` | New integration test |
| `testdata/lego_themes.txtar` | New integration test |
| `requirements.md` | Mark all LEGO catalog endpoints as implemented |
