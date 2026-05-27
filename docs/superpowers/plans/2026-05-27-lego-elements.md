# LEGO Elements API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `lego elements get` CLI command backed by the Rebrickable `GET /api/v3/lego/elements/{element_id}/` endpoint.

**Architecture:** Add an `Element` struct (embedding existing `Part` and `PartColor`) to `api.go`; add `GetLegoElement` to `lego.go`; create a new `lego_elements.go` for the Cobra command; cover with a unit test and a txtar integration test.

**Tech Stack:** Go, Cobra, resty, Bazel, httptest (unit tests), txtar (integration tests)

---

### Task 1: Add data model and write failing unit test

**Files:**
- Modify: `cli/cmd/api/api.go`
- Modify: `cli/cmd/api/api_test.go`

- [ ] **Step 1: Add `Element` struct to `api.go`**

Append to the end of `cli/cmd/api/api.go` (after `SetPartsResponse`):

```go
type Element struct {
	ElementID string    `json:"element_id"`
	Part      Part      `json:"part"`
	Color     PartColor `json:"color"`
	DesignID  string    `json:"design_id"`
}
```

- [ ] **Step 2: Add failing test to `api_test.go`**

Append to the end of `cli/cmd/api/api_test.go`:

```go
func TestGetLegoElement(t *testing.T) {
	tests := []struct {
		name       string
		response   Element
		statusCode int
		wantErr    bool
	}{
		{
			"returns element",
			Element{
				ElementID: "4119739",
				Part:      Part{PartNum: "3001", Name: "Brick 2 x 4"},
				Color:     PartColor{ID: 1, Name: "Blue"},
				DesignID:  "3001",
			},
			200, false,
		},
		{"not found", Element{}, 404, true},
		{"server error", Element{}, 500, true},
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
			result, err := client.GetLegoElement("4119739")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoElement() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.ElementID != tt.response.ElementID {
				t.Errorf("GetLegoElement() element_id = %v, want %v", result.ElementID, tt.response.ElementID)
			}
		})
	}
}
```

- [ ] **Step 3: Run tests to confirm they fail**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: compilation error — `GetLegoElement` undefined. If it compiles and passes, something is wrong.

---

### Task 2: Implement API client method

**Files:**
- Modify: `cli/cmd/api/lego.go`

- [ ] **Step 1: Add `GetLegoElement` to `lego.go`**

Append to the end of `cli/cmd/api/lego.go`:

```go
func (c *Client) GetLegoElement(elementID string) (*Element, error) {
	result := &Element{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/elements/%s/", elementID))
	if err != nil {
		return nil, fmt.Errorf("get lego element request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego element failed with status %d", resp.StatusCode())
	}
	return result, nil
}
```

- [ ] **Step 2: Run unit tests to confirm they pass**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: all tests pass including `TestGetLegoElement`.

- [ ] **Step 3: Commit**

```bash
git add cli/cmd/api/api.go cli/cmd/api/lego.go cli/cmd/api/api_test.go
git commit -m "feat: add Element type and GetLegoElement API client method"
```

---

### Task 3: Add Cobra command

**Files:**
- Create: `cli/cmd/lego_elements.go`
- Modify: `cli/cmd/BUILD.bazel`

- [ ] **Step 1: Create `cli/cmd/lego_elements.go`**

```go
package cmd

import (
	"github.com/spf13/cobra"
)

var elementID string

func init() {
	legoElementsCommands()
}

func legoElementsCommands() {
	legoCmd.AddCommand(legoElementsCmd)
	legoElementsCmd.AddCommand(getLegoElementCmd)

	getLegoElementCmd.Flags().StringVarP(&elementID, "id", "i", "", "Element id")
	_ = getLegoElementCmd.MarkFlagRequired("id")
}

var legoElementsCmd = &cobra.Command{
	Use:   "elements",
	Short: "LEGO catalog element actions",
}

var getLegoElementCmd = &cobra.Command{
	Use:   "get",
	Short: "get an element by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoElement(elementID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
```

- [ ] **Step 2: Add `lego_elements.go` to `cli/cmd/BUILD.bazel`**

In `cli/cmd/BUILD.bazel`, add `"lego_elements.go"` to the `srcs` list (alphabetically between `lego_colors.go` and `lego_parts.go`):

```python
    srcs = [
        "lego.go",
        "lego_colors.go",
        "lego_elements.go",
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
git add cli/cmd/lego_elements.go cli/cmd/BUILD.bazel
git commit -m "feat: add lego elements CLI command"
```

---

### Task 4: Add integration test

**Files:**
- Create: `testdata/lego_elements.txtar`

- [ ] **Step 1: Create `testdata/lego_elements.txtar`**

```
# Get a specific element by stable ID
exec cli lego elements get --id 4119739
stdout '"element_id"'
! stderr .
```

- [ ] **Step 2: Run integration tests to confirm they pass**

```bash
bazel test //cli:cli_test --test_output=all
```

Expected: all integration tests pass including the new element test. (Requires `REBRICKABLE_API_KEY`, `REBRICKABLE_USERNAME`, `REBRICKABLE_PASSWORD` in the environment — loaded automatically via `.envrc` if direnv is active.)

- [ ] **Step 3: Commit**

```bash
git add testdata/lego_elements.txtar
git commit -m "test: add lego elements integration test"
```

---

### Task 5: Update requirements doc

**Files:**
- Modify: `requirements.md`

- [ ] **Step 1: Update `requirements.md`**

Remove the Elements section from `### Not yet implemented`. Add a `#### Elements ✅` section under `### Implemented ✅`:

```markdown
#### Elements ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/elements/{element_id}/` |
```

Update the summary line at the top:

```
**Implemented:** Parts (5), Sets (6), Colors (2), Elements (1) — 14 of 20 endpoints.
**Not yet implemented:** Minifigs (4), Part Categories (2), Themes (2) — 6 of 20 endpoints.
```

- [ ] **Step 2: Commit**

```bash
git add requirements.md
git commit -m "docs: mark LEGO elements endpoint as implemented"
```
