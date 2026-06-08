# LEGO Part Categories and Themes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `lego part-categories list/get` and `lego themes list/get` CLI commands backed by the four remaining Rebrickable LEGO catalog endpoints.

**Architecture:** Add `PartCategory`, `PartCategoriesResponse`, `Theme`, and `ThemesResponse` types to `api.go`; add four client methods to `lego.go` (two paginated list, two direct GET); create `lego_part_categories.go` and `lego_themes.go` for Cobra commands; cover with four unit tests and two txtar integration tests. IDs are passed as strings (consistent with `GetLegoColor`, `GetLegoElement`).

**Tech Stack:** Go, Cobra, resty, Bazel, httptest (unit tests), txtar (integration tests)

---

### Task 1: Add data model and write failing unit tests

**Files:**
- Modify: `cli/cmd/api/api.go`
- Modify: `cli/cmd/api/api_test.go`

- [ ] **Step 1: Add four types to `api.go`**

Append to the end of `cli/cmd/api/api.go`:

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

- [ ] **Step 2: Add four failing tests to `api_test.go`**

Append to the end of `cli/cmd/api/api_test.go`:

```go
func TestGetLegoPartCategories(t *testing.T) {
	tests := []struct {
		name       string
		response   PartCategoriesResponse
		statusCode int
		wantErr    bool
	}{
		{"returns categories", PartCategoriesResponse{Count: 1, Results: []PartCategory{{ID: 1, Name: "Baseplates", PartCount: 243}}}, 200, false},
		{"server error", PartCategoriesResponse{}, 500, true},
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
			result, err := client.GetLegoPartCategories()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoPartCategories() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoPartCategories() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoPartCategory(t *testing.T) {
	tests := []struct {
		name       string
		response   PartCategory
		statusCode int
		wantErr    bool
	}{
		{"returns category", PartCategory{ID: 1, Name: "Baseplates", PartCount: 243}, 200, false},
		{"not found", PartCategory{}, 404, true},
		{"server error", PartCategory{}, 500, true},
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
			result, err := client.GetLegoPartCategory("1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoPartCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.ID != tt.response.ID {
				t.Errorf("GetLegoPartCategory() id = %v, want %v", result.ID, tt.response.ID)
			}
		})
	}
}

func TestGetLegoThemes(t *testing.T) {
	tests := []struct {
		name       string
		response   ThemesResponse
		statusCode int
		wantErr    bool
	}{
		{"returns themes", ThemesResponse{Count: 1, Results: []Theme{{ID: 1, Name: "Technic"}}}, 200, false},
		{"server error", ThemesResponse{}, 500, true},
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
			result, err := client.GetLegoThemes()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoThemes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoThemes() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoTheme(t *testing.T) {
	tests := []struct {
		name       string
		response   Theme
		statusCode int
		wantErr    bool
	}{
		{"returns theme", Theme{ID: 1, Name: "Technic"}, 200, false},
		{"not found", Theme{}, 404, true},
		{"server error", Theme{}, 500, true},
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
			result, err := client.GetLegoTheme("1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoTheme() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.ID != tt.response.ID {
				t.Errorf("GetLegoTheme() id = %v, want %v", result.ID, tt.response.ID)
			}
		})
	}
}
```

- [ ] **Step 3: Run tests to confirm they fail**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: compilation error — `GetLegoPartCategories`, `GetLegoPartCategory`, `GetLegoThemes`, `GetLegoTheme` undefined.

---

### Task 2: Implement API client methods

**Files:**
- Modify: `cli/cmd/api/lego.go`

- [ ] **Step 1: Append four methods to `lego.go`**

```go
func (c *Client) GetLegoPartCategories() (*PartCategoriesResponse, error) {
	count, results, err := fetchAllPages[PartCategory](c.http, "/lego/part_categories/")
	if err != nil {
		return nil, fmt.Errorf("get lego part categories: %w", err)
	}
	return &PartCategoriesResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoPartCategory(id string) (*PartCategory, error) {
	result := &PartCategory{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/part_categories/%s/", id))
	if err != nil {
		return nil, fmt.Errorf("get lego part category request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego part category failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) GetLegoThemes() (*ThemesResponse, error) {
	count, results, err := fetchAllPages[Theme](c.http, "/lego/themes/")
	if err != nil {
		return nil, fmt.Errorf("get lego themes: %w", err)
	}
	return &ThemesResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoTheme(id string) (*Theme, error) {
	result := &Theme{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/themes/%s/", id))
	if err != nil {
		return nil, fmt.Errorf("get lego theme request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego theme failed with status %d", resp.StatusCode())
	}
	return result, nil
}
```

- [ ] **Step 2: Run unit tests to confirm they pass**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: all tests pass including all four new `TestGetLegoPartCategor*` and `TestGetLegoTheme*` tests.

- [ ] **Step 3: Commit**

```bash
git add cli/cmd/api/api.go cli/cmd/api/lego.go cli/cmd/api/api_test.go
git commit -m "feat: add PartCategory/Theme types and API client methods"
```

---

### Task 3: Add Cobra commands

**Files:**
- Create: `cli/cmd/lego_part_categories.go`
- Create: `cli/cmd/lego_themes.go`
- Modify: `cli/cmd/BUILD.bazel`

- [ ] **Step 1: Create `cli/cmd/lego_part_categories.go`**

```go
package cmd

import (
	"github.com/spf13/cobra"
)

var partCategoryID string

func init() {
	legoPartCategoriesCommands()
}

func legoPartCategoriesCommands() {
	legoCmd.AddCommand(legoPartCategoriesCmd)
	legoPartCategoriesCmd.AddCommand(getLegoPartCategoriesCmd)
	legoPartCategoriesCmd.AddCommand(getLegoPartCategoryCmd)

	getLegoPartCategoryCmd.Flags().StringVarP(&partCategoryID, "id", "i", "", "Part category id")
	_ = getLegoPartCategoryCmd.MarkFlagRequired("id")
}

var legoPartCategoriesCmd = &cobra.Command{
	Use:   "part-categories",
	Short: "LEGO catalog part category actions",
}

var getLegoPartCategoriesCmd = &cobra.Command{
	Use:   "list",
	Short: "list all part categories",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPartCategories()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoPartCategoryCmd = &cobra.Command{
	Use:   "get",
	Short: "get a part category by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoPartCategory(partCategoryID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
```

- [ ] **Step 2: Create `cli/cmd/lego_themes.go`**

```go
package cmd

import (
	"github.com/spf13/cobra"
)

var themeID string

func init() {
	legoThemesCommands()
}

func legoThemesCommands() {
	legoCmd.AddCommand(legoThemesCmd)
	legoThemesCmd.AddCommand(getLegoThemesCmd)
	legoThemesCmd.AddCommand(getLegoThemeCmd)

	getLegoThemeCmd.Flags().StringVarP(&themeID, "id", "i", "", "Theme id")
	_ = getLegoThemeCmd.MarkFlagRequired("id")
}

var legoThemesCmd = &cobra.Command{
	Use:   "themes",
	Short: "LEGO catalog theme actions",
}

var getLegoThemesCmd = &cobra.Command{
	Use:   "list",
	Short: "list all themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoThemes()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var getLegoThemeCmd = &cobra.Command{
	Use:   "get",
	Short: "get a theme by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newLegoAPIClient(cmd)
		result, err := client.GetLegoTheme(themeID)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}
```

- [ ] **Step 3: Add both files to `cli/cmd/BUILD.bazel`**

In `cli/cmd/BUILD.bazel`, add `"lego_part_categories.go"` and `"lego_themes.go"` to the `srcs` list alphabetically:

```python
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
```

- [ ] **Step 4: Build to confirm it compiles**

```bash
bazel build //cli
```

Expected: build succeeds with no errors.

- [ ] **Step 5: Run all unit tests**

```bash
bazel test //cli/cmd/api:api_test --test_output=all
```

Expected: all tests pass.

- [ ] **Step 6: Commit**

```bash
git add cli/cmd/lego_part_categories.go cli/cmd/lego_themes.go cli/cmd/BUILD.bazel
git commit -m "feat: add lego part-categories and themes CLI commands"
```

---

### Task 4: Add integration tests

**Files:**
- Create: `testdata/lego_part_categories.txtar`
- Create: `testdata/lego_themes.txtar`

- [ ] **Step 1: Create `testdata/lego_part_categories.txtar`**

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

- [ ] **Step 2: Create `testdata/lego_themes.txtar`**

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

- [ ] **Step 3: Run integration tests**

```bash
bazel test //cli:cli_test --test_output=all
```

Expected: all integration tests pass. (Requires `REBRICKABLE_API_KEY` in the environment — loaded automatically via `.envrc` if direnv is active.)

- [ ] **Step 4: Commit**

```bash
git add testdata/lego_part_categories.txtar testdata/lego_themes.txtar
git commit -m "test: add lego part-categories and themes integration tests"
```

---

### Task 5: Update requirements doc

**Files:**
- Modify: `requirements.md`

- [ ] **Step 1: Update `requirements.md`**

Remove the entire `### Not yet implemented` section (it will be empty after this). Add `#### Part Categories ✅` and `#### Themes ✅` sections under `### Implemented ✅`. Update the summary line.

The top of the LEGO Catalog section should read:

```markdown
**Implemented:** Parts (5), Sets (6), Colors (2), Elements (1), Minifigs (4), Part Categories (2), Themes (2) — 20 of 20 endpoints.
```

Remove:
```markdown
### Not yet implemented

#### Part Categories

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/lego/part_categories/` | Change `extend-lego-catalog-part-categories` proposed, not yet applied |
| GET | `/api/v3/lego/part_categories/{id}/` | Change `extend-lego-catalog-part-categories` proposed, not yet applied |

#### Themes

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/themes/` |
| GET | `/api/v3/lego/themes/{id}/` |
```

Add after `#### Minifigs ✅` under `### Implemented ✅`:

```markdown
#### Part Categories ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/part_categories/` |
| GET | `/api/v3/lego/part_categories/{id}/` |

#### Themes ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/themes/` |
| GET | `/api/v3/lego/themes/{id}/` |
```

- [ ] **Step 2: Commit**

```bash
git add requirements.md
git commit -m "docs: mark LEGO part categories and themes as implemented"
```
