# CLI Migration to rebrickable-go Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace `cli/cmd/api/` with the standalone `github.com/bu3/rebrickable-go` library, restoring mutating-command output in the CLI layer.

**Architecture:** The CLI imports `github.com/bu3/rebrickable-go`, stores a `*rebrickable.Client` in the Cobra command context (instead of raw token strings), and each Cobra command prints its own confirmation for mutating operations. The `cli/cmd/api/` package is deleted entirely.

**Tech Stack:** Go 1.23, Cobra, Bazel (rules_go + Gazelle), `github.com/bu3/rebrickable-go@v0.1.0`

---

### File Map

| File | Change |
|------|--------|
| `go.mod` | Add `github.com/bu3/rebrickable-go v0.1.0` |
| `MODULE.bazel` | Add `com_github_bu3_rebrickable_go` to `use_repo` |
| `cli/cmd/user.go` | Replace `login()`/resty with `rebrickable.NewAuthenticatedClient`; store `*rebrickable.Client` in context |
| `cli/cmd/lego.go` | Replace `api.NewLegoClient` with `rebrickable.NewClient`; store client in context |
| `cli/cmd/sets.go` | Extract client from context; add `fmt.Println` to all mutating commands |
| `cli/cmd/lego_parts.go` | Change `api.PartsFilter` → `rebrickable.PartsFilter`; remove dead nil-check branches |
| `cli/cmd/BUILD.bazel` | Swap `//cli/cmd/api` + resty → `@com_github_bu3_rebrickable_go//:rebrickable-go` |
| `cli/cmd/api/` | Delete entirely |

---

### Task 1: Add rebrickable-go to the Go and Bazel module graphs

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`
- Modify: `MODULE.bazel`

- [ ] **Step 1: Add the dependency**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-cli
go get github.com/bu3/rebrickable-go@v0.1.0
go mod tidy
```

Expected: `go.mod` now lists `github.com/bu3/rebrickable-go v0.1.0` under `require`.

- [ ] **Step 2: Register the new repo in MODULE.bazel**

Update `MODULE.bazel` — add `"com_github_bu3_rebrickable_go"` to `use_repo`:

```python
module(
    name = "rebrickable-cli",
    version = "0.0.0",
)

bazel_dep(name = "rules_go", version = "0.59.0")
bazel_dep(name = "gazelle", version = "0.47.0")

go_sdk = use_extension("@rules_go//go:extensions.bzl", "go_sdk")
go_sdk.download(version = "1.23.4")

go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
go_deps.from_file(go_mod = "//:go.mod")
use_repo(
    go_deps,
    "com_github_bu3_rebrickable_go",
    "com_github_go_resty_resty_v2",
    "com_github_rogpeppe_go_internal",
    "com_github_spf13_cobra",
    "org_golang_x_net",
)
```

- [ ] **Step 3: Verify the Bazel target name**

```bash
bazel query @com_github_bu3_rebrickable_go//...
```

Expected output will include a target like `@com_github_bu3_rebrickable_go//:rebrickable-go` or `@com_github_bu3_rebrickable_go//:rebrickable`. Note the exact name — it is used in Task 6's BUILD.bazel update.

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum MODULE.bazel MODULE.bazel.lock
git commit -m "chore: add github.com/bu3/rebrickable-go dependency"
```

---

### Task 2: Rewrite `cli/cmd/user.go`

**Files:**
- Modify: `cli/cmd/user.go`

The current file calls `login()` via a raw resty client and stores two string context values (`AuthToken`, `ApiKey`). After this task it calls `rebrickable.NewAuthenticatedClient` and stores a single `*rebrickable.Client` in context under a typed key, removing resty and the `api` package from this file entirely.

The typed context key (`contextKey` type + `rebrickableClient` constant) lives here and is shared across the `cmd` package — `lego.go` and `sets.go` use the same key to read the client back.

- [ ] **Step 1: Write the new `cli/cmd/user.go`**

```go
package cmd

import (
	"os"

	rebrickable "github.com/bu3/rebrickable-go"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type contextKey string

const rebrickableClient contextKey = "rebrickable_client"

func init() {
	rootCmd.AddCommand(user)
}

var user = &cobra.Command{
	Use:   "user",
	Short: "user actions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("REBRICKABLE_API_KEY")
		username := os.Getenv("REBRICKABLE_USERNAME")
		password := os.Getenv("REBRICKABLE_PASSWORD")
		client, err := rebrickable.NewAuthenticatedClient(apiKey, username, password)
		if err != nil {
			return err
		}
		cmd.SetContext(context.WithValue(cmd.Context(), rebrickableClient, client))
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
```

- [ ] **Step 2: Verify compilation (BUILD.bazel still references `//cli/cmd/api` — that is OK for now, fix in Task 6)**

```bash
bazel build //cli/cmd:cmd 2>&1 | head -40
```

Expected: compile errors because `sets.go` and `lego.go` still import `api`. That is expected — they are fixed in subsequent tasks. Verify that the errors are NOT about `user.go` itself.

- [ ] **Step 3: Commit**

```bash
git add cli/cmd/user.go
git commit -m "refactor: replace login() with rebrickable.NewAuthenticatedClient in user.go"
```

---

### Task 3: Rewrite `cli/cmd/lego.go`

**Files:**
- Modify: `cli/cmd/lego.go`

Remove `const LegoApiKey`, remove the `api` import, call `rebrickable.NewClient(apiKey)`, store the client in context under the shared `rebrickableClient` key defined in `user.go`. Update `newLegoAPIClient` to extract the client from context.

- [ ] **Step 1: Write the new `cli/cmd/lego.go`**

```go
package cmd

import (
	"os"

	rebrickable "github.com/bu3/rebrickable-go"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func init() {
	rootCmd.AddCommand(legoCmd)
}

var legoCmd = &cobra.Command{
	Use:   "lego",
	Short: "LEGO catalog actions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("REBRICKABLE_API_KEY")
		client := rebrickable.NewClient(apiKey)
		cmd.SetContext(context.WithValue(cmd.Context(), rebrickableClient, client))
		return nil
	},
}

func newLegoAPIClient(cmd *cobra.Command) *rebrickable.Client {
	return cmd.Context().Value(rebrickableClient).(*rebrickable.Client)
}
```

- [ ] **Step 2: Commit**

```bash
git add cli/cmd/lego.go
git commit -m "refactor: replace api.NewLegoClient with rebrickable.NewClient in lego.go"
```

---

### Task 4: Update `cli/cmd/sets.go`

**Files:**
- Modify: `cli/cmd/sets.go`

Replace `newAPIClient` to extract `*rebrickable.Client` from context. Remove the `api` import. Add `fmt.Println`/`fmt.Printf` output to every mutating command — the library no longer prints anything on mutation.

- [ ] **Step 1: Update `newAPIClient` and imports**

Replace the import block and `newAPIClient` function. Remove:
```go
"github.com/bu3/rebrickable-cli/cli/cmd/api"
```
Add:
```go
rebrickable "github.com/bu3/rebrickable-go"
```

Replace `newAPIClient`:
```go
func newAPIClient(cmd *cobra.Command) *rebrickable.Client {
	return cmd.Context().Value(rebrickableClient).(*rebrickable.Client)
}
```

The full import block becomes:
```go
import (
	"encoding/json"
	"fmt"
	"strings"

	rebrickable "github.com/bu3/rebrickable-go"
	"github.com/spf13/cobra"
)
```

- [ ] **Step 2: Add output to `saveSetListCmd`**

```go
var saveSetListCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.StoreUserSetList(setListName); err != nil {
			return err
		}
		fmt.Println("SetList saved")
		return nil
	},
}
```

- [ ] **Step 3: Add output to `updateSetListCmd`**

```go
var updateSetListCmd = &cobra.Command{
	Use:   "update",
	Short: "update",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.UpdateUserSetList(setListID, setListName); err != nil {
			return err
		}
		fmt.Printf("Updated set list: %s\n", setListID)
		return nil
	},
}
```

- [ ] **Step 4: Add output to `replaceSetListCmd`**

```go
var replaceSetListCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.ReplaceUserSetList(setListID, setListName); err != nil {
			return err
		}
		fmt.Printf("Replaced set list: %s\n", setListID)
		return nil
	},
}
```

- [ ] **Step 5: Add output to `deleteSetListsCmd`**

```go
var deleteSetListsCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.DeleteUserSetList(setNumber); err != nil {
			return err
		}
		fmt.Printf("Deleted set list: %s\n", setNumber)
		return nil
	},
}
```

- [ ] **Step 6: Add output to `saveSetListSetCmd`**

```go
var saveSetListSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.StoreUserSetListSet(setListID, adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Println("Set added to set list")
		return nil
	},
}
```

- [ ] **Step 7: Add output to `deleteSetListSetCmd`**

```go
var deleteSetListSetCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.DeleteUserSetListSet(setListID, adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Printf("Deleted %s from set list %s\n", adjustedSetNumber(), setListID)
		return nil
	},
}
```

- [ ] **Step 8: Add output to `saveSetsCmd`**

```go
var saveSetsCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.StoreUserSet(adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Println("Set saved")
		return nil
	},
}
```

- [ ] **Step 9: Add output to `replaceSetCmd`**

```go
var replaceSetCmd = &cobra.Command{
	Use:   "replace",
	Short: "replace",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.ReplaceUserSet(adjustedSetNumber(), quantity); err != nil {
			return err
		}
		fmt.Printf("Updated set: %s\n", adjustedSetNumber())
		return nil
	},
}
```

- [ ] **Step 10: Add output to `deleteSetsCmd`**

```go
var deleteSetsCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newAPIClient(cmd)
		if err := client.DeleteUserSet(adjustedSetNumber()); err != nil {
			return err
		}
		fmt.Printf("Deleted set: %s\n", adjustedSetNumber())
		return nil
	},
}
```

- [ ] **Step 11: Commit**

```bash
git add cli/cmd/sets.go
git commit -m "refactor: migrate sets.go to rebrickable-go client; restore mutation output"
```

---

### Task 5: Update `cli/cmd/lego_parts.go`

**Files:**
- Modify: `cli/cmd/lego_parts.go`

Replace `api.PartsFilter` with `rebrickable.PartsFilter`. Remove the dead nil-check branches in `getLegoPartCmd` and `getLegoPartColorCmd` — the library now returns an error on 404 (not `nil, nil`), so those branches can never be reached.

- [ ] **Step 1: Update import and `partsFilter` variable**

Replace:
```go
import (
	"encoding/json"
	"fmt"

	"github.com/bu3/rebrickable-cli/cli/cmd/api"
	"github.com/spf13/cobra"
)

var (
	partNumber  string
	partColorID string
	partsFilter api.PartsFilter
)
```

With:
```go
import (
	"encoding/json"
	"fmt"

	rebrickable "github.com/bu3/rebrickable-go"
	"github.com/spf13/cobra"
)

var (
	partNumber  string
	partColorID string
	partsFilter rebrickable.PartsFilter
)
```

- [ ] **Step 2: Remove dead nil-check in `getLegoPartCmd`**

Replace:
```go
var getLegoPartCmd = &cobra.Command{
	Use:   "get",
	Short: "get a part by part_num",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPart(partNumber)
		if err != nil {
			return err
		}
		if result == nil {
			return nil
		}
		return printJSON(result)
	},
}
```

With:
```go
var getLegoPartCmd = &cobra.Command{
	Use:   "get",
	Short: "get a part by part_num",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPart(partNumber)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
```

- [ ] **Step 3: Remove dead nil-check in `getLegoPartColorCmd`**

Replace:
```go
var getLegoPartColorCmd = &cobra.Command{
	Use:   "colorDetail",
	Short: "get a specific part/color combination",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPartColor(partNumber, partColorID)
		if err != nil {
			return err
		}
		if result == nil {
			return nil
		}
		return printJSON(result)
	},
}
```

With:
```go
var getLegoPartColorCmd = &cobra.Command{
	Use:   "colorDetail",
	Short: "get a specific part/color combination",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPartColor(partNumber, partColorID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
```

- [ ] **Step 4: Commit**

```bash
git add cli/cmd/lego_parts.go
git commit -m "refactor: migrate lego_parts.go to rebrickable-go PartsFilter type"
```

---

### Task 6: Update `cli/cmd/BUILD.bazel`

**Files:**
- Modify: `cli/cmd/BUILD.bazel`

Remove `//cli/cmd/api` and `@com_github_go_resty_resty_v2//:resty` from deps. Add the rebrickable-go Bazel dep. Note: use the exact target name verified in Task 1, Step 3.

- [ ] **Step 1: Write the new BUILD.bazel**

Use the target name from `bazel query @com_github_bu3_rebrickable_go//...` (Task 1, Step 3). Expected target: `@com_github_bu3_rebrickable_go//:rebrickable-go`.

```python
load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = [
        "lego.go",
        "lego_colors.go",
        "lego_elements.go",
        "lego_minifigs.go",
        "lego_part_categories.go",
        "lego_parts.go",
        "lego_sets.go",
        "lego_themes.go",
        "root.go",
        "sets.go",
        "user.go",
    ],
    importpath = "github.com/bu3/rebrickable-cli/cli/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_bu3_rebrickable_go//:rebrickable-go",
        "@com_github_spf13_cobra//:cobra",
        "@org_golang_x_net//context",
    ],
)
```

If the query in Task 1 returned a different name (e.g., `:rebrickable`), use that instead.

- [ ] **Step 2: Verify the build compiles (api/ still exists at this point)**

```bash
bazel build //cli/cmd:cmd
```

Expected: success. The `//cli/cmd/api` dep is gone, so Bazel fetches and links the external module.

- [ ] **Step 3: Commit**

```bash
git add cli/cmd/BUILD.bazel
git commit -m "build: swap //cli/cmd/api for rebrickable-go in cmd BUILD.bazel"
```

---

### Task 7: Delete `cli/cmd/api/` and run full test suite

**Files:**
- Delete: `cli/cmd/api/api.go`
- Delete: `cli/cmd/api/api_test.go`
- Delete: `cli/cmd/api/lego.go`
- Delete: `cli/cmd/api/lego_parts.go`
- Delete: `cli/cmd/api/user.go`
- Delete: `cli/cmd/api/BUILD.bazel`

The tests from `api_test.go` have already been migrated into `rebrickable-go/client_test.go`. Nothing is lost.

- [ ] **Step 1: Delete the api package**

```bash
rm -rf /Users/fabio.mangione/workspace/my-stuff/rebrickable-cli/cli/cmd/api
```

- [ ] **Step 2: Run the unit test suite**

```bash
bazel test //cli/cmd/api:api_test 2>&1 | head -5
```

Expected: `ERROR: no such target '//cli/cmd/api:api_test'` — the target is gone. That is correct.

```bash
bazel build //...
```

Expected: success with no references to `//cli/cmd/api`.

- [ ] **Step 3: Run all tests**

```bash
bazel test //... --test_output=all
```

Expected: all tests pass. If integration tests (`//cli:cli_test`) fail with 429 (rate limiting), re-run just the unit tests:

```bash
bazel test //cli/cmd:cmd_test --test_output=all 2>/dev/null || echo "No unit tests in cmd (expected)"
bazel build //cli
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "refactor: delete cli/cmd/api — replaced by github.com/bu3/rebrickable-go"
```

---

## Self-Review

**Spec coverage check against `docs/superpowers/specs/2026-06-08-rebrickable-go-library-design.md` (CLI Update section):**

| Spec requirement | Task |
|-----------------|------|
| Remove `cli/cmd/api/` entirely | Task 7 |
| Add `github.com/bu3/rebrickable-go` to `go.mod` | Task 1 |
| Replace `newLegoAPIClient(cmd)` with `rebrickable.NewClient(apiKey)` | Task 3 |
| Replace `login()` with `rebrickable.NewAuthenticatedClient(apiKey, username, password)` | Task 2 |
| Each mutating Cobra command adds its own output | Task 4 |
| Update all BUILD.bazel files | Task 6 |

All requirements covered. No spec gaps found.

**Placeholder scan:** No TBDs, no "implement later", no vague steps. All code blocks are complete.

**Type consistency:**
- `*rebrickable.Client` used consistently across Tasks 2–5
- `rebrickableClient` context key defined once in Task 2 (`user.go`), referenced in Tasks 3–4
- `rebrickable.PartsFilter` used in Task 5, matches library's exported type name
