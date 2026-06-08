# rebrickable-go Library — Design

**Date:** 2026-06-08

## Summary

Extract the Rebrickable API client from `rebrickable-cli/cli/cmd/api` into a standalone Go library at `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go`. The library is a clean Go module with no Cobra, no Bazel, no CLI concerns. It can be imported by the CLI, a future web app, or any other consumer.

This spec covers **only the library**. The CLI update (removing `cli/cmd/api` and importing the library) is a separate follow-up cycle.

---

## Repository Setup

- **Path:** `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go`
- **Module:** `github.com/bu3/rebrickable-go`
- **Package:** `rebrickable`
- **Go version:** 1.23 (match `rebrickable-cli`)
- **Dependency:** `github.com/go-resty/resty/v2` only
- **Build:** standard `go test ./...` — no Bazel

Initialize as a git repo with a `go.mod` and `go.sum`.

---

## File Structure

```
rebrickable-go/
  go.mod
  go.sum
  client.go       — Client struct, NewClient, NewAuthenticatedClient, fetchAllPages, userPath
  types.go        — all domain types (moved from cli/cmd/api/api.go)
  lego.go         — LEGO catalog methods (moved from cli/cmd/api/lego.go)
  lego_parts.go   — LEGO parts methods (moved from cli/cmd/api/lego_parts.go)
  user.go         — user endpoint methods (moved from cli/cmd/api/user.go)
  client_test.go  — all unit tests (moved from cli/cmd/api/api_test.go)
```

---

## `client.go`

Replaces `api.go`. Contains the `Client` struct, constructors, and the `fetchAllPages` helper.

**Removed from `api.go`:**
- `GetURL()` — CLI-specific utility, not needed in a library
- `NewClient(apiKey, authToken string)` — replaced by the two new constructors below
- `NewLegoClient(apiKey string)` — replaced by `NewClient(apiKey string)`

**Kept (unexported):**
- `fetchAllPages[T any]`
- `newClientWithBaseURL(apiKey, authToken, baseURL string) *Client` — used in tests
- `userPath(path string) string`

**New public API:**

```go
const defaultBaseURL = "https://rebrickable.com/api/v3"

type Client struct {
    http      *resty.Client
    authToken string
}

// NewClient creates a client for read-only LEGO catalog endpoints (no user token).
func NewClient(apiKey string) *Client {
    return newClientWithBaseURL(apiKey, "", defaultBaseURL)
}

// NewAuthenticatedClient creates a client for user endpoints by calling POST /users/_token/.
// Returns an error if authentication fails.
func NewAuthenticatedClient(apiKey, username, password string) (*Client, error) {
    c := newClientWithBaseURL(apiKey, "", defaultBaseURL)
    token, err := c.getUserToken(username, password)
    if err != nil {
        return nil, fmt.Errorf("authentication failed: %w", err)
    }
    return newClientWithBaseURL(apiKey, token, defaultBaseURL), nil
}
```

**`getUserToken` (unexported):** moves the login logic from `cli/cmd/user.go` into the library:

```go
func (c *Client) getUserToken(username, password string) (string, error) {
    type tokenResponse struct {
        UserToken string `json:"user_token"`
    }
    result := &tokenResponse{}
    resp, err := c.http.R().
        SetHeader("Content-Type", "application/x-www-form-urlencoded").
        SetFormData(map[string]string{"username": username, "password": password}).
        SetResult(result).
        Post("/users/_token/")
    if err != nil {
        return "", fmt.Errorf("token request failed: %w", err)
    }
    if resp.StatusCode() != 200 {
        return "", fmt.Errorf("token request failed with status %d", resp.StatusCode())
    }
    return result.UserToken, nil
}
```

---

## `types.go`

All domain types moved verbatim from `cli/cmd/api/api.go`, with `package rebrickable` at the top. No changes to struct fields or JSON tags.

Types to move: `Set`, `LegoSetsResponse`, `SetMinifig`, `SetMinifigsResponse`, `SetPart`, `Part`, `PartColor`, `SetPartsResponse`, `SetList`, `SetListsResponse`, `UserSet`, `SetsResponse`, `PartColors`, `PartColorsResponse`, `ColorsResponse`, `Element`, `Minifig`, `MinifigsResponse`, `PartCategory`, `PartCategoriesResponse`, `Theme`, `ThemesResponse`.

---

## `lego.go` and `lego_parts.go`

Move verbatim from `cli/cmd/api/lego.go` and `cli/cmd/api/lego_parts.go`. Only change: `package rebrickable` at the top. All method signatures and logic are identical.

---

## `user.go`

Move from `cli/cmd/api/user.go` with one important change: **remove all `fmt.Println` and `fmt.Printf` calls**. Libraries must not write to stdout — that is the CLI layer's responsibility.

Affected methods and what to remove:

| Method | Side effect to remove |
|--------|----------------------|
| `StoreUserSetList` | `fmt.Println("SetList saved")` |
| `UpdateUserSetList` | `fmt.Printf("Updated set list: %s\n", listID)` |
| `ReplaceUserSetList` | `fmt.Printf("Replaced set list: %s\n", listID)` |
| `DeleteUserSetList` | `fmt.Printf("Set list %s not found\n", id)` and `fmt.Printf("Deleted set list: %s\n", id)` |
| `StoreUserSetListSet` | `fmt.Println("Set added to set list")` |
| `DeleteUserSetListSet` | `fmt.Printf("Set %s not found in set list %s\n", ...)` and `fmt.Printf("Deleted set %s from set list %s\n", ...)` |
| `StoreUserSet` | `fmt.Println("Set saved")` |
| `ReplaceUserSet` | `fmt.Printf("Updated set: %s\n", setNum)` |
| `DeleteUserSet` | `fmt.Printf("Set %s not found\n", setNumber)` and `fmt.Printf("Deleted set: %s\n", setNumber)` |

`DeleteUserSetList`, `DeleteUserSetListSet`, and `DeleteUserSet` currently use `fmt.Printf` to signal a 404 instead of returning an error. In the library, 404 responses on delete operations return `nil` (idempotent delete — not an error). The CLI layer handles informing the user.

The `fmt` import is removed from `user.go` since it is no longer used after removing print calls. The `fmt` package in those methods is only used for `Sprintf` in URL construction (via `userPath`) and for the print statements — `Sprintf` calls in `userPath` calls remain, but `userPath` is defined in `client.go`, so `user.go` itself only needs `fmt` if there are remaining `fmt.Errorf` or `fmt.Sprintf` calls. Check at implementation time.

---

## `client_test.go`

All tests from `cli/cmd/api/api_test.go` move here, with:
- `package rebrickable` at the top (white-box — needed to access `newClientWithBaseURL`)
- Imports updated: remove `github.com/bu3/rebrickable-cli/cli/cmd/api`, no import needed since types are in the same package
- All test logic identical

Run with: `go test ./...`

---

## What Is NOT in the Library

- No Cobra, no CLI commands
- No `GetURL()` helper
- No Bazel BUILD files
- No txtar integration tests
- No `fmt.Println`/`fmt.Printf` in any method

---

## CLI Update (Second Cycle — Separate Spec)

After the library is published:
- Remove `cli/cmd/api/` entirely
- Add `github.com/bu3/rebrickable-go` to `go.mod`
- Replace `newLegoAPIClient(cmd)` with `rebrickable.NewClient(apiKey)`
- Replace `cli/cmd/user.go`'s `login()` with `rebrickable.NewAuthenticatedClient(apiKey, username, password)`
- Each Cobra command that calls a mutating user method adds its own output (e.g., `fmt.Println("Set saved")`)
- Update all BUILD.bazel files

---

## Files Changed

| File | Change |
|------|--------|
| `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/go.mod` | New — module `github.com/bu3/rebrickable-go` |
| `rebrickable-go/client.go` | New — Client, constructors, helpers |
| `rebrickable-go/types.go` | New — all domain types |
| `rebrickable-go/lego.go` | New — LEGO catalog methods |
| `rebrickable-go/lego_parts.go` | New — LEGO parts methods |
| `rebrickable-go/user.go` | New — user methods (no stdout side effects) |
| `rebrickable-go/client_test.go` | New — all unit tests |
