# Code Review Analysis

**Last Updated:** 2026-02-05

## Overview

This is a well-structured Go CLI for the Rebrickable LEGO API. The architecture follows good separation of concerns with distinct `cmd` and `api` packages. However, there are several areas that need improvement.

**Progress:** 1 issue fixed, 1 partially fixed, 6 remaining.

---

## Critical Issues

### 1. ~~Login Function Returns Success on Failure~~ ✅ FIXED

**Location:** `cli/cmd/user.go:61-66`

**Status:** This issue has been resolved. The login function now properly returns errors:

```go
if err != nil {
    return nil, fmt.Errorf("login request failed: %w", err)
}
if resp.StatusCode() != 200 {
    return nil, fmt.Errorf("login failed with status %d", resp.StatusCode())
}
```

### 2. Silent Error Discarding Throughout API Package (Partially Fixed)

**Locations still affected:**
- `cli/cmd/api/api.go:20` - `StoreUserSetList`
- `cli/cmd/api/api.go:34` - `GetUserSetLists`
- `cli/cmd/api/api.go:48` - `DeleteUserSetList`
- `cli/cmd/api/api.go:60` - `StoreUserSet`
- `cli/cmd/api/api.go:90` - `DeleteUserSet`

**Fixed:**
- ✅ `cli/cmd/api/api.go:71-85` - `GetUserSets` now returns errors properly

**Problem:** Network failures, auth errors, and API issues are silently swallowed in the remaining functions. Users won't know why operations fail.

**Fix:** Remaining API functions should return `error` and handle it in callers.

### 3. Silent Error Discarding in Commands

**Location:** `cli/cmd/sets.go:95-96`

```go
setsResponse, _ := api.GetUserSets(client, apiKey, authToken)
output, _ := json.MarshalIndent(setsResponse, "", "\t")
```

**Problem:** Errors are explicitly ignored with `_`.

---

## Design Issues

### 4. Weak Type Safety for Results

**Location:** `cli/cmd/api/api.go:100-103`

```go
type SetsResponse struct {
    Count   int              `json:"count"`
    Results []map[string]any `json:"results"`
}
```

**Problem:** Using `map[string]any` loses type safety and makes it hard to work with the data.

**Fix:** Define proper structs for the API response:

```go
type Set struct {
    SetNum   string `json:"set_num"`
    Name     string `json:"name"`
    Year     int    `json:"year"`
    NumParts int    `json:"num_parts"`
    // ... other fields
}

type SetsResponse struct {
    Count   int   `json:"count"`
    Results []Set `json:"results"`
}
```

### 5. Unsafe Context Type Assertions

**Location:** `cli/cmd/sets.go` (8 places)

```go
authToken := cmd.Context().Value(AuthToken).(string)
apiKey := cmd.Context().Value(ApiKey).(string)
```

**Problem:** If context values are missing, this panics.

**Fix:** Use safe type assertions with error handling:

```go
authToken, ok := cmd.Context().Value(AuthToken).(string)
if !ok || authToken == "" {
    return fmt.Errorf("not authenticated")
}
```

### 6. Context Keys Should Be Typed

**Location:** `cli/cmd/user.go:17-20`

```go
const (
    ApiKey    = "api_key"
    AuthToken = "auth_token"
)
```

**Problem:** String keys can collide with other packages. Go idiom is to use unexported typed keys.

**Fix:**

```go
type contextKey string

const (
    apiKeyKey    contextKey = "api_key"
    authTokenKey contextKey = "auth_token"
)
```

---

## Code Duplication

### 7. HTTP Client Created Multiple Times

A new `resty.New()` client is created in every command handler. Consider injecting a shared client or creating it once in `PersistentPreRunE`.

### 8. Repeated Header Setup

**Location:** `cli/cmd/api/api.go`

Every API function manually sets the same headers:

```go
SetHeader("Content-Type", "application/json").
SetHeader("Authorization", fmt.Sprintf("key %s", apiKey)).
```

**Fix:** Create a configured client once:

```go
func NewAuthenticatedClient(apiKey string) *resty.Client {
    return resty.New().
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", fmt.Sprintf("key %s", apiKey))
}
```

---

## Minor Issues

### 9. Unnecessary `fmt.Sprintf`

**Locations:** `cli/cmd/api/api.go:15, 54, 96`

```go
return fmt.Sprintf(apiBaseURI + path)  // Sprintf not needed
fmt.Println(fmt.Sprintf("Deleted set: %s", id))  // Redundant
```

**Fix:**

```go
return apiBaseURI + path
fmt.Printf("Deleted set: %s\n", id)
```

### 10. Global Mutable State for Flags

**Location:** `cli/cmd/sets.go:12-13`

```go
var setNumber string
var setListName string
```

**Problem:** Module-level mutable variables can cause issues with testing and concurrent usage.

### 11. Undocumented Magic Suffix

**Location:** `cli/cmd/sets.go:126-132`

```go
func adjustedSetNumber() string {
    if !strings.HasSuffix(setNumber, "-1") {
        return setNumber + "-1"
    }
    return setNumber
}
```

**Problem:** The "-1" suffix requirement is unexplained. A comment would help explain this is the Rebrickable set variant convention.

### 12. Missing Required Flag Validation

Commands like `saveSetsCmd` don't validate that required flags are provided before making API calls.

---

## Summary Table

| Priority | Issue | Location | Status |
|----------|-------|----------|--------|
| ~~**Critical**~~ | ~~Login returns nil on failure~~ | ~~`user.go:60-63`~~ | ✅ Fixed |
| **Critical** | Errors silently discarded | 5 functions in `api.go` | ⚠️ Partial (`GetUserSets` fixed) |
| **High** | Unsafe type assertions | `sets.go` (8 places) | ❌ Open |
| **High** | Weak type safety | `api.go:100-103` | ❌ Open |
| **Medium** | Code duplication | HTTP client/headers | ❌ Open |
| **Medium** | Untyped context keys | `user.go:17-20` | ❌ Open |
| **Low** | Global flag variables | `sets.go:12-13` | ❌ Open |
| **Low** | Missing documentation | `adjustedSetNumber()` | ❌ Open |

---

## What's Done Well

- Clean separation between `cmd` and `api` packages
- Proper use of Cobra's `PersistentPreRunE` for authentication
- Credentials kept in environment variables (not hardcoded)
- Bazel build setup with proper dependency management
- Integration tests using testscript framework
- Token masking in debug output (`api.go:89`)

---

## Maintenance

To refresh this document after making code changes:

```
check CODE_REVIEW.md, review the code and update the doc
```
