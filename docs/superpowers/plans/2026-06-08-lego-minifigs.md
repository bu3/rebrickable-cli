# LEGO Minifigs API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `lego minifigs list/get/parts/sets` CLI commands backed by the four Rebrickable `/lego/minifigs/` endpoints.

**Architecture:** Add `Minifig` and `MinifigsResponse` types to `api.go`; add four client methods to `lego.go` (two paginated, two direct GET); create `lego_minifigs.go` for Cobra commands; cover with four unit tests and a txtar integration test. The `/parts/` and `/sets/` sub-endpoints reuse existing `SetPartsResponse` and `LegoSetsResponse` types — no extra types needed for those.

**Tech Stack:** Go, Cobra, resty, Bazel, httptest (unit tests), txtar (integration tests)

---

### Task 1: Add data model and write failing unit tests

**Files:**
- Modify: `cli/cmd/api/api.go`
- Modify: `cli/cmd/api/api_test.go`

- [ ] **Step 1: Add `Minifig` and `MinifigsResponse` to `api.go`**

Append to the end of `cli/cmd/api/api.go`:

```go
type Minifig struct {
	SetNum         string `json:"set_num"`
	Name           string `json:"name"`
	NumParts       int    `json:"num_parts"`
	SetImgURL      string `json:"set_img_url"`
	SetURL         string `json:"set_url"`
	LastModifiedDt string `json:"last_modified_dt"`
}

type MinifigsResponse struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Minifig `json:"results"`
}
```

- [ ] **Step 2: Add four failing tests to `api_test.go`**

Append to the end of `cli/cmd/api/api_test.go`:

```go
func TestGetLegoMinifigs(t *testing.T) {
	tests := []struct {
		name       string
		response   MinifigsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns minifigs", MinifigsResponse{Count: 1, Results: []Minifig{{SetNum: "fig-000001", Name: "Spaceman"}}}, 200, false},
		{"server error", MinifigsResponse{}, 500, true},
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
			result, err := client.GetLegoMinifigs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoMinifigs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoMinifigs() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoMinifig(t *testing.T) {
	tests := []struct {
		name       string
		response   Minifig
		statusCode int
		wantErr    bool
	}{
		{"returns minifig", Minifig{SetNum: "fig-000001", Name: "Spaceman", NumParts: 4}, 200, false},
		{"not found", Minifig{}, 404, true},
		{"server error", Minifig{}, 500, true},
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
			result, err := client.GetLegoMinifig("fig-000001")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoMinifig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.SetNum != tt.response.SetNum {
				t.Errorf("GetLegoMinifig() set_num = %v, want %v", result.SetNum, tt.response.SetNum)
			}
		})
	}
}

func TestGetLegoMinifigParts(t *testing.T) {
	tests := []struct {
		name       string
		response   SetPartsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns parts", SetPartsResponse{Count: 2, Results: []SetPart{{Quantity: 1, Part: Part{PartNum: "3001"}}}}, 200, false},
		{"server error", SetPartsResponse{}, 500, true},
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
			result, err := client.GetLegoMinifigParts("fig-000001")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoMinifigParts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoMinifigParts() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoMinifigSets(t *testing.T) {
	tests := []struct {
		name       string
		response   LegoSetsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns sets", LegoSetsResponse{Count: 1, Results: []Set{{SetNum: "10497-1", Name: "Galaxy Explorer"}}}, 200, false},
		{"server error", LegoSetsResponse{}, 500, true},
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
			result, err := client.GetLegoMinifigSets("fig-000001")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoMinifigSets() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoMinifigSets() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}
```

- [ ] **Step 3: Run tests to confirm they fail**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: compilation error — `GetLegoMinifigs`, `GetLegoMinifig`, `GetLegoMinifigParts`, `GetLegoMinifigSets` undefined.

---

### Task 2: Implement API client methods

**Files:**
- Modify: `cli/cmd/api/lego.go`

- [ ] **Step 1: Append four methods to `lego.go`**

```go
func (c *Client) GetLegoMinifigs() (*MinifigsResponse, error) {
	count, results, err := fetchAllPages[Minifig](c.http, "/lego/minifigs/")
	if err != nil {
		return nil, fmt.Errorf("get lego minifigs: %w", err)
	}
	return &MinifigsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoMinifig(figNum string) (*Minifig, error) {
	result := &Minifig{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/minifigs/%s/", figNum))
	if err != nil {
		return nil, fmt.Errorf("get lego minifig request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego minifig failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) GetLegoMinifigParts(figNum string) (*SetPartsResponse, error) {
	count, results, err := fetchAllPages[SetPart](c.http, fmt.Sprintf("/lego/minifigs/%s/parts/", figNum))
	if err != nil {
		return nil, fmt.Errorf("get lego minifig parts: %w", err)
	}
	return &SetPartsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoMinifigSets(figNum string) (*LegoSetsResponse, error) {
	count, results, err := fetchAllPages[Set](c.http, fmt.Sprintf("/lego/minifigs/%s/sets/", figNum))
	if err != nil {
		return nil, fmt.Errorf("get lego minifig sets: %w", err)
	}
	return &LegoSetsResponse{Count: count, Results: results}, nil
}
```

- [ ] **Step 2: Run unit tests to confirm they pass**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: all tests pass including all four new `TestGetLegoMinifig*` tests.

- [ ] **Step 3: Commit**

```bash
git add cli/cmd/api/api.go cli/cmd/api/lego.go cli/cmd/api/api_test.go
git commit -m "feat: add Minifig types and GetLegoMinifig* API client methods"
```

---

### Task 3: Add Cobra commands

**Files:**
- Create: `cli/cmd/lego_minifigs.go`
- Modify: `cli/cmd/BUILD.bazel`

- [ ] **Step 1: Create `cli/cmd/lego_minifigs.go`**

```go
package cmd

import (
	"github.com/spf13/cobra"
)

var figNum string

func init() {
	legoMinifigsCommands()
}

func legoMinifigsCommands() {
	legoCmd.AddCommand(legoMinifigsCmd)
	legoMinifigsCmd.AddCommand(getLegoMinifigsCmd)
	legoMinifigsCmd.AddCommand(getLegoMinifigCmd)
	legoMinifigsCmd.AddCommand(getLegoMinifigPartsCmd)
	legoMinifigsCmd.AddCommand(getLegoMinifigSetsCmd)

	getLegoMinifigCmd.Flags().StringVarP(&figNum, "fig_num", "f", "", "Minifig number")
	_ = getLegoMinifigCmd.MarkFlagRequired("fig_num")

	getLegoMinifigPartsCmd.Flags().StringVarP(&figNum, "fig_num", "f", "", "Minifig number")
	_ = getLegoMinifigPartsCmd.MarkFlagRequired("fig_num")

	getLegoMinifigSetsCmd.Flags().StringVarP(&figNum, "fig_num", "f", "", "Minifig number")
	_ = getLegoMinifigSetsCmd.MarkFlagRequired("fig_num")
}

var legoMinifigsCmd = &cobra.Command{
	Use:   "minifigs",
	Short: "LEGO catalog minifig actions",
}

var getLegoMinifigsCmd = &cobra.Command{
	Use:   "list",
	Short: "list all minifigs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoMinifigs()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoMinifigCmd = &cobra.Command{
	Use:   "get",
	Short: "get a minifig by fig_num",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoMinifig(figNum)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoMinifigPartsCmd = &cobra.Command{
	Use:   "parts",
	Short: "list parts of a minifig",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoMinifigParts(figNum)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoMinifigSetsCmd = &cobra.Command{
	Use:   "sets",
	Short: "list sets containing a minifig",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoMinifigSets(figNum)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
```

- [ ] **Step 2: Add `lego_minifigs.go` to `cli/cmd/BUILD.bazel`**

In `cli/cmd/BUILD.bazel`, add `"lego_minifigs.go"` to the `srcs` list alphabetically between `lego_elements.go` and `lego_parts.go`:

```python
    srcs = [
        "lego.go",
        "lego_colors.go",
        "lego_elements.go",
        "lego_minifigs.go",
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
git add cli/cmd/lego_minifigs.go cli/cmd/BUILD.bazel
git commit -m "feat: add lego minifigs CLI commands"
```

---

### Task 4: Add integration test

**Files:**
- Create: `testdata/lego_minifigs.txtar`

- [ ] **Step 1: Create `testdata/lego_minifigs.txtar`**

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

- [ ] **Step 2: Run integration tests**

```bash
bazel test //cli:cli_test --test_output=all
```

Expected: all integration tests pass. (Requires `REBRICKABLE_API_KEY`, `REBRICKABLE_USERNAME`, `REBRICKABLE_PASSWORD` in the environment — loaded automatically via `.envrc` if direnv is active.)

- [ ] **Step 3: Commit**

```bash
git add testdata/lego_minifigs.txtar
git commit -m "test: add lego minifigs integration test"
```

---

### Task 5: Update requirements doc

**Files:**
- Modify: `requirements.md`

- [ ] **Step 1: Update `requirements.md`**

Remove the entire Minifigs section from `### Not yet implemented`. Add a `#### Minifigs ✅` section under `### Implemented ✅`:

```markdown
#### Minifigs ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/minifigs/` |
| GET | `/api/v3/lego/minifigs/{set_num}/` |
| GET | `/api/v3/lego/minifigs/{set_num}/parts/` |
| GET | `/api/v3/lego/minifigs/{set_num}/sets/` |
```

Update the summary line at the top:

```
**Implemented:** Parts (5), Sets (6), Colors (2), Elements (1), Minifigs (4) — 18 of 20 endpoints.
**Not yet implemented:** Part Categories (2), Themes (2) — 2 of 20 endpoints.
```

- [ ] **Step 2: Commit**

```bash
git add requirements.md
git commit -m "docs: mark LEGO minifigs endpoints as implemented"
```
