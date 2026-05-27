# LEGO Elements API — Design

**Date:** 2026-05-27

## Summary

Implement the single LEGO elements catalog endpoint:

- `GET /api/v3/lego/elements/{element_id}/` — get a specific element by ID

An "element" is a specific part+color combination as sold by LEGO, identified by a numeric element ID. This is a read-only, unauthenticated endpoint (API key only, no user token).

---

## Data Model

**File:** `cli/cmd/api/api.go`

Add a new `Element` struct embedding the existing `Part` and `PartColor` types:

```go
type Element struct {
    ElementID string    `json:"element_id"`
    Part      Part      `json:"part"`
    Color     PartColor `json:"color"`
    DesignID  string    `json:"design_id"`
}
```

No new response wrapper type is needed — the endpoint returns a single object, not a paginated list.

---

## API Client

**File:** `cli/cmd/api/lego.go`

One new method on `*Client`:

- `GetLegoElement(elementID string) (*Element, error)` — direct GET to `/lego/elements/{elementID}/`, returns error on non-200. Same shape as `GetLegoColor`.

---

## CLI Commands

**File:** `cli/cmd/lego_elements.go` (new file)

New `lego elements` subcommand registered on `legoCmd`, with one sub-subcommand:

| Command | Flags | Action |
|---------|-------|--------|
| `lego elements get` | `--id` (required) | calls `GetLegoElement(id)`, prints JSON |

Registration follows the same `legoElementsCommands()` + `init()` pattern as `lego_colors.go`. Uses `newLegoAPIClient(cmd)` and `printJSON()`.

---

## Tests

### Unit tests

**File:** `cli/cmd/api/api_test.go`

One table-driven test using `httptest.NewServer` + `newClientWithBaseURL`:

- `TestGetLegoElement` — cases: returns element (200), not found (404), server error (500)

### Integration test (txtar)

**File:** `testdata/lego_elements.txtar` (new file)

```
# Get a specific element by stable ID
exec cli lego elements get --id 4119739
stdout '"element_id"'
! stderr .
```

Element ID `4119739` is a stable, well-known identifier in the Rebrickable dataset (Blue 2×4 Brick).

---

## Files Changed

| File | Change |
|------|--------|
| `cli/cmd/api/api.go` | Add `Element` struct |
| `cli/cmd/api/lego.go` | Add `GetLegoElement` |
| `cli/cmd/api/api_test.go` | Add `TestGetLegoElement` |
| `cli/cmd/lego_elements.go` | New file — Cobra commands |
| `cli/cmd/BUILD.bazel` | Add `lego_elements.go` to srcs |
| `testdata/lego_elements.txtar` | New integration test |
| `requirements.md` | Mark Elements endpoint as implemented |
