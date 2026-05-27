# LEGO Colors API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `lego colors list` and `lego colors get` CLI commands backed by the Rebrickable `/lego/colors/` and `/lego/colors/{id}/` API endpoints.

**Architecture:** Expand `PartColor` in `api.go` with `RGB`/`IsTrans` fields and add a `ColorsResponse` type; add two client methods to `lego.go`; add a new `lego_colors.go` file with Cobra commands; cover with unit tests in `api_test.go` and an integration test in `testdata/lego_colors.txtar`.

**Tech Stack:** Go, Cobra, resty, Bazel, httptest (unit tests), txtar (integration tests)

---

### Task 1: Expand data model and write failing unit tests

**Files:**
- Modify: `cli/cmd/api/api.go`
- Modify: `cli/cmd/api/api_test.go`

- [ ] **Step 1: Expand `PartColor` and add `ColorsResponse` in `api.go`**

Find the `PartColor` struct (currently around line 149) and replace it:

```go
type PartColor struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	RGB     string `json:"rgb"`
	IsTrans bool   `json:"is_trans"`
}
```

Add `ColorsResponse` immediately after (alongside the other `*Response` types):

```go
type ColorsResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []PartColor `json:"results"`
}
```

- [ ] **Step 2: Add failing tests to `api_test.go`**

Append to the end of `cli/cmd/api/api_test.go`:

```go
func TestGetLegoColors(t *testing.T) {
	tests := []struct {
		name       string
		response   ColorsResponse
		statusCode int
		wantErr    bool
	}{
		{
			"returns colors",
			ColorsResponse{Count: 2, Results: []PartColor{
				{ID: 0, Name: "Black", RGB: "05131D", IsTrans: false},
				{ID: 1, Name: "Blue", RGB: "0055BF", IsTrans: false},
			}},
			200, false,
		},
		{"server error", ColorsResponse{}, 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == 200 {
					_ = json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "", server.URL)
			result, err := client.GetLegoColors()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoColors() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoColors() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoColor(t *testing.T) {
	tests := []struct {
		name       string
		response   PartColor
		statusCode int
		wantErr    bool
	}{
		{"returns color", PartColor{ID: 0, Name: "Black", RGB: "05131D", IsTrans: false}, 200, false},
		{"not found", PartColor{}, 404, true},
		{"server error", PartColor{}, 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == 200 {
					_ = json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "", server.URL)
			result, err := client.GetLegoColor("0")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoColor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.ID != tt.response.ID {
				t.Errorf("GetLegoColor() id = %v, want %v", result.ID, tt.response.ID)
			}
		})
	}
}
```

- [ ] **Step 3: Run tests to confirm they fail**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: compilation error — `GetLegoColors` and `GetLegoColor` undefined. If tests compile and pass, something is wrong.

---

### Task 2: Implement API client methods

**Files:**
- Modify: `cli/cmd/api/lego.go`

- [ ] **Step 1: Add `GetLegoColors` and `GetLegoColor` to `lego.go`**

Append to the end of `cli/cmd/api/lego.go`:

```go
func (c *Client) GetLegoColors() (*ColorsResponse, error) {
	count, results, err := fetchAllPages[PartColor](c.http, "/lego/colors/")
	if err != nil {
		return nil, fmt.Errorf("get lego colors: %w", err)
	}
	return &ColorsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoColor(id string) (*PartColor, error) {
	result := &PartColor{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/colors/%s/", id))
	if err != nil {
		return nil, fmt.Errorf("get lego color request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego color failed with status %d", resp.StatusCode())
	}
	return result, nil
}
```

- [ ] **Step 2: Run unit tests to confirm they pass**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: all tests pass, including `TestGetLegoColors` and `TestGetLegoColor`.

- [ ] **Step 3: Commit**

```bash
git add cli/cmd/api/api.go cli/cmd/api/lego.go cli/cmd/api/api_test.go
git commit -m "feat: add GetLegoColors and GetLegoColor API client methods"
```

---

### Task 3: Add Cobra commands

**Files:**
- Create: `cli/cmd/lego_colors.go`
- Modify: `cli/cmd/BUILD.bazel`

- [ ] **Step 1: Create `cli/cmd/lego_colors.go`**

```go
package cmd

import (
	"github.com/spf13/cobra"
)

var colorID string

func init() {
	legoColorsCommands()
}

func legoColorsCommands() {
	legoCmd.AddCommand(legoColorsCmd)
	legoColorsCmd.AddCommand(getLegoColorsCmd)
	legoColorsCmd.AddCommand(getLegoColorCmd)

	getLegoColorCmd.Flags().StringVarP(&colorID, "id", "i", "", "Color id")
	_ = getLegoColorCmd.MarkFlagRequired("id")
}

var legoColorsCmd = &cobra.Command{
	Use:   "colors",
	Short: "LEGO catalog color actions",
}

var getLegoColorsCmd = &cobra.Command{
	Use:   "list",
	Short: "list all colors",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoColors()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoColorCmd = &cobra.Command{
	Use:   "get",
	Short: "get a color by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoColor(colorID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
```

- [ ] **Step 2: Add `lego_colors.go` to `cli/cmd/BUILD.bazel`**

In `cli/cmd/BUILD.bazel`, add `"lego_colors.go"` to the `srcs` list:

```python
    srcs = [
        "lego.go",
        "lego_colors.go",
        "lego_parts.go",
        "lego_sets.go",
        "root.go",
        "sets.go",
        "user.go",
    ],
```

- [ ] **Step 3: Build to confirm it compiles**

```bash
bazel build //cli
```

Expected: build succeeds with no errors.

- [ ] **Step 4: Run all unit tests**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: all tests pass.

- [ ] **Step 5: Commit**

```bash
git add cli/cmd/lego_colors.go cli/cmd/BUILD.bazel
git commit -m "feat: add lego colors CLI commands"
```

---

### Task 4: Add integration test

**Files:**
- Create: `testdata/lego_colors.txtar`

- [ ] **Step 1: Create `testdata/lego_colors.txtar`**

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

- [ ] **Step 2: Run integration tests to confirm they pass**

```bash
bazel test //cli:cli_test --test_output=all
```

Expected: all integration tests pass, including the new color tests. (Requires `REBRICKABLE_API_KEY`, `REBRICKABLE_USERNAME`, `REBRICKABLE_PASSWORD` in the environment — loaded automatically via `.envrc` if direnv is active.)

- [ ] **Step 3: Commit**

```bash
git add testdata/lego_colors.txtar
git commit -m "test: add lego colors integration test"
```

---

### Task 5: Update requirements doc

**Files:**
- Modify: `requirements.md`

- [ ] **Step 1: Mark Colors as implemented in `requirements.md`**

In the `#### Colors` section under `### Not yet implemented`, remove the two color rows. Then add a `#### Colors ✅` section under `### Implemented ✅`:

```markdown
#### Colors ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/colors/` |
| GET | `/api/v3/lego/colors/{id}/` |
```

Also update the summary line at the top:

```
**Implemented:** Parts (5), Sets (6), Colors (2) — 13 of 20 endpoints.
**Not yet implemented:** Elements (1), Minifigs (4), Part Categories (2), Themes (2) — 7 of 20 endpoints.
```

- [ ] **Step 2: Commit**

```bash
git add requirements.md
git commit -m "docs: mark LEGO colors endpoints as implemented"
```
