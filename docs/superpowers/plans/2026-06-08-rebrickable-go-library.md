# rebrickable-go Library Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extract the Rebrickable API client from `rebrickable-cli/cli/cmd/api/` into a standalone Go library at `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go`.

**Architecture:** A clean Go module (`github.com/bu3/rebrickable-go`, package `rebrickable`) with no Cobra, no Bazel, no CLI concerns. All domain types, HTTP client logic, and endpoint methods move verbatim from the CLI; the only changes are: package name → `rebrickable`, two new public constructors replace the old ones, `fmt.Print*` side effects are removed from all methods, and 404 on GET methods returns an error (not `nil, nil`).

**Tech Stack:** Go 1.23, `github.com/go-resty/resty/v2`, standard `go test ./...`

---

## File Structure

| File | Responsibility |
|------|---------------|
| `go.mod` / `go.sum` | Module definition |
| `client.go` | `Client` struct, `NewClient`, `NewAuthenticatedClient`, `getUserToken`, `fetchAllPages`, `userPath`, `newClientWithBaseURL` |
| `types.go` | All domain types (verbatim from `cli/cmd/api/api.go`) |
| `lego.go` | LEGO catalog methods (verbatim from `cli/cmd/api/lego.go`) |
| `lego_parts.go` | LEGO parts methods (from `cli/cmd/api/lego_parts.go`; `fmt.Printf` removed; 404 returns error) |
| `user.go` | User endpoint methods (from `cli/cmd/api/user.go`; all stdout removed) |
| `client_test.go` | All tests (from `cli/cmd/api/api_test.go`; `TestGetURL` dropped; `wantNil` cases updated) |

---

### Task 1: Repository Bootstrap

**Files:**
- Create: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/go.mod`

- [ ] **Step 1: Initialize git repo and go module**

```bash
cd /Users/fabio.mangione/workspace/my-stuff
mkdir rebrickable-go
cd rebrickable-go
git init
go mod init github.com/bu3/rebrickable-go
```

Expected: `go.mod` created with `module github.com/bu3/rebrickable-go` and `go 1.23.X` (match output).

- [ ] **Step 2: Add resty dependency**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go get github.com/go-resty/resty/v2
```

Expected: `go.sum` created, `go.mod` has `require github.com/go-resty/resty/v2 v2.x.x`.

- [ ] **Step 3: Commit**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
git add go.mod go.sum
git commit -m "chore: initialize rebrickable-go module"
```

---

### Task 2: `types.go` — Domain Types

**Files:**
- Create: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/types.go`

No tests needed — pure data structs.

- [ ] **Step 1: Create `types.go`**

Copy all types verbatim from `cli/cmd/api/api.go`, changing only the package name:

```go
package rebrickable

type Set struct {
	SetNum         string `json:"set_num"`
	Name           string `json:"name"`
	Year           int    `json:"year"`
	ThemeID        int    `json:"theme_id"`
	NumParts       int    `json:"num_parts"`
	SetImgURL      string `json:"set_img_url"`
	SetURL         string `json:"set_url"`
	LastModifiedDt string `json:"last_modified_dt"`
}

type UserSet struct {
	ListID        int  `json:"list_id"`
	Quantity      int  `json:"quantity"`
	IncludeSpares bool `json:"include_spares"`
	Set           Set  `json:"set"`
}

type SetsResponse struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []UserSet `json:"results"`
}

type SetList struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	NumSets int    `json:"num_sets"`
	IsBuild bool   `json:"is_build_list"`
}

type SetListsResponse struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []SetList `json:"results"`
}

type LegoSetsResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Set  `json:"results"`
}

type SetMinifig struct {
	ID       int    `json:"id"`
	SetNum   string `json:"set_num"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	NumParts int    `json:"num_parts"`
	ImgURL   string `json:"set_img_url"`
}

type SetMinifigsResponse struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []SetMinifig `json:"results"`
}

type Part struct {
	PartNum    string `json:"part_num"`
	Name       string `json:"name"`
	PartCatID  int    `json:"part_cat_id"`
	PartURL    string `json:"part_url"`
	PartImgURL string `json:"part_img_url"`
}

type PartColor struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	RGB     string `json:"rgb"`
	IsTrans bool   `json:"is_trans"`
}

type ColorsResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []PartColor `json:"results"`
}

type SetPart struct {
	ID        int       `json:"id"`
	InvPartID int       `json:"inv_part_id"`
	Part      Part      `json:"part"`
	Color     PartColor `json:"color"`
	Quantity  int       `json:"quantity"`
	IsSpare   bool      `json:"is_spare"`
	NumSets   int       `json:"num_sets"`
}

type SetPartsResponse struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []SetPart `json:"results"`
}

type Element struct {
	ElementID string    `json:"element_id"`
	Part      Part      `json:"part"`
	Color     PartColor `json:"color"`
	DesignID  string    `json:"design_id"` // not in OpenAPI spec but returned by the API
}

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

- [ ] **Step 2: Verify it compiles**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go build ./...
```

Expected: no errors, no output.

- [ ] **Step 3: Commit**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
git add types.go
git commit -m "feat: add domain types"
```

---

### Task 3: `client.go` — Client, Constructors, Helpers

**Files:**
- Create: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/client.go`
- Create: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/client_test.go` (partial — constructors and fetchAllPages)

- [ ] **Step 1: Write failing tests for `NewClient` and `NewAuthenticatedClient`**

```go
package rebrickable

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_SetsAuthorizationHeader(t *testing.T) {
	var capturedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(LegoSetsResponse{Count: 0})
	}))
	defer server.Close()

	c := newClientWithBaseURL("mykey", "", server.URL)
	_, _ = c.GetLegoSets()
	if capturedAuth != "key mykey" {
		t.Errorf("Authorization = %q, want %q", capturedAuth, "key mykey")
	}
}

func TestNewAuthenticatedClient_Success(t *testing.T) {
	type tokenResp struct {
		UserToken string `json:"user_token"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(tokenResp{UserToken: "tok123"})
	}))
	defer server.Close()

	c := newClientWithBaseURL("apikey", "", server.URL)
	token, err := c.getUserToken("user", "pass")
	if err != nil {
		t.Fatalf("getUserToken() error = %v", err)
	}
	if token != "tok123" {
		t.Errorf("getUserToken() = %q, want %q", token, "tok123")
	}
}

func TestNewAuthenticatedClient_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	c := newClientWithBaseURL("apikey", "", server.URL)
	_, err := c.getUserToken("user", "wrong")
	if err == nil {
		t.Error("getUserToken() expected error for 401, got nil")
	}
}
```

Save this to `client_test.go`.

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./...
```

Expected: FAIL — `undefined: newClientWithBaseURL`, `undefined: Client`, etc.

- [ ] **Step 3: Create `client.go`**

```go
package rebrickable

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const defaultBaseURL = "https://rebrickable.com/api/v3"

type Client struct {
	http      *resty.Client
	authToken string
}

func NewClient(apiKey string) *Client {
	return newClientWithBaseURL(apiKey, "", defaultBaseURL)
}

func NewAuthenticatedClient(apiKey, username, password string) (*Client, error) {
	c := newClientWithBaseURL(apiKey, "", defaultBaseURL)
	token, err := c.getUserToken(username, password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	return newClientWithBaseURL(apiKey, token, defaultBaseURL), nil
}

func newClientWithBaseURL(apiKey, authToken, baseURL string) *Client {
	http := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey))

	return &Client{http: http, authToken: authToken}
}

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

func (c *Client) userPath(path string) string {
	return fmt.Sprintf("/users/%s%s", c.authToken, path)
}

func fetchAllPages[T any](httpClient *resty.Client, firstURL string) (int, []T, error) {
	type page struct {
		Count   int    `json:"count"`
		Next    string `json:"next"`
		Results []T    `json:"results"`
	}

	var all []T
	totalCount := 0
	isFirst := true
	url := firstURL

	for {
		p := &page{}
		resp, err := httpClient.R().SetResult(p).Get(url)
		if err != nil {
			return 0, nil, fmt.Errorf("request failed: %w", err)
		}
		if resp.StatusCode() != 200 {
			return 0, nil, fmt.Errorf("unexpected status %d", resp.StatusCode())
		}
		if isFirst {
			totalCount = p.Count
			isFirst = false
		}
		all = append(all, p.Results...)
		if p.Next == "" {
			break
		}
		url = p.Next
	}
	return totalCount, all, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./...
```

Expected: PASS — `ok github.com/bu3/rebrickable-go`

- [ ] **Step 5: Commit**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
git add client.go client_test.go
git commit -m "feat: add Client struct, constructors, and fetchAllPages helper"
```

---

### Task 4: `lego.go` — LEGO Catalog Methods + Tests

**Files:**
- Create: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/lego.go`
- Modify: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/client_test.go` (append lego tests)

Source: `rebrickable-cli/cli/cmd/api/lego.go` (177 lines) — copy verbatim, change package only.

- [ ] **Step 1: Write failing tests for lego catalog methods**

Append to `client_test.go`:

```go
func TestGetLegoSets(t *testing.T) {
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
			result, err := client.GetLegoSets()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoSets() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoSets() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoSet(t *testing.T) {
	tests := []struct {
		name       string
		response   Set
		statusCode int
		wantErr    bool
	}{
		{"returns set", Set{SetNum: "10497-1", Name: "Galaxy Explorer", Year: 2022}, 200, false},
		{"not found", Set{}, 404, true},
		{"server error", Set{}, 500, true},
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
			result, err := client.GetLegoSet("10497-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoSet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.SetNum != tt.response.SetNum {
				t.Errorf("GetLegoSet() set_num = %v, want %v", result.SetNum, tt.response.SetNum)
			}
		})
	}
}

func TestGetLegoSetAlternates(t *testing.T) {
	tests := []struct {
		name       string
		response   LegoSetsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns alternates", LegoSetsResponse{Count: 1, Results: []Set{{SetNum: "moc-1234"}}}, 200, false},
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
			result, err := client.GetLegoSetAlternates("10497-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoSetAlternates() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoSetAlternates() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoSetMinifigs(t *testing.T) {
	tests := []struct {
		name       string
		response   SetMinifigsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns minifigs", SetMinifigsResponse{Count: 2, Results: []SetMinifig{{SetNum: "fig-001", Name: "Astronaut"}}}, 200, false},
		{"server error", SetMinifigsResponse{}, 500, true},
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
			result, err := client.GetLegoSetMinifigs("10497-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoSetMinifigs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoSetMinifigs() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoSetParts(t *testing.T) {
	tests := []struct {
		name       string
		response   SetPartsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns parts", SetPartsResponse{Count: 3, Results: []SetPart{{Quantity: 2, Part: Part{PartNum: "3001"}}}}, 200, false},
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
			result, err := client.GetLegoSetParts("10497-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoSetParts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoSetParts() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoSetSets(t *testing.T) {
	tests := []struct {
		name       string
		response   LegoSetsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns sub-sets", LegoSetsResponse{Count: 1, Results: []Set{{SetNum: "75192-1"}}}, 200, false},
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
			result, err := client.GetLegoSetSets("10497-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoSetSets() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoSetSets() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

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

func TestGetLegoSetsPagination(t *testing.T) {
	page1Set := Set{SetNum: "10497-1", Name: "Galaxy Explorer"}
	page2Set := Set{SetNum: "75192-1", Name: "Millennium Falcon"}

	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if r.URL.Query().Get("page") == "2" {
			_ = json.NewEncoder(w).Encode(LegoSetsResponse{Count: 2, Results: []Set{page2Set}})
		} else {
			_ = json.NewEncoder(w).Encode(LegoSetsResponse{Count: 2, Next: serverURL + "/?page=2", Results: []Set{page1Set}})
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := newClientWithBaseURL("key", "", server.URL)
	result, err := client.GetLegoSets()

	if err != nil {
		t.Fatalf("GetLegoSets() unexpected error: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("GetLegoSets() count = %d, want 2", result.Count)
	}
	if len(result.Results) != 2 {
		t.Errorf("GetLegoSets() len(results) = %d, want 2", len(result.Results))
	}
	if result.Results[0].SetNum != page1Set.SetNum {
		t.Errorf("GetLegoSets() results[0].SetNum = %q, want %q", result.Results[0].SetNum, page1Set.SetNum)
	}
	if result.Results[1].SetNum != page2Set.SetNum {
		t.Errorf("GetLegoSets() results[1].SetNum = %q, want %q", result.Results[1].SetNum, page2Set.SetNum)
	}
}

func TestGetLegoSetsPaginationErrorOnSecondPage(t *testing.T) {
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("page") == "2" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(LegoSetsResponse{Count: 2, Next: serverURL + "/?page=2", Results: []Set{{SetNum: "10497-1"}}})
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := newClientWithBaseURL("key", "", server.URL)
	_, err := client.GetLegoSets()

	if err == nil {
		t.Error("GetLegoSets() expected error on second page, got nil")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./...
```

Expected: FAIL — methods undefined.

- [ ] **Step 3: Create `lego.go`**

Copy verbatim from `rebrickable-cli/cli/cmd/api/lego.go`, changing only `package api` → `package rebrickable`:

```go
package rebrickable

import "fmt"

func (c *Client) GetLegoSets() (*LegoSetsResponse, error) {
	count, results, err := fetchAllPages[Set](c.http, "/lego/sets/")
	if err != nil {
		return nil, fmt.Errorf("get lego sets: %w", err)
	}
	return &LegoSetsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoSet(setNum string) (*Set, error) {
	result := &Set{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/sets/%s/", setNum))

	if err != nil {
		return nil, fmt.Errorf("get lego set request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego set failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) GetLegoSetAlternates(setNum string) (*LegoSetsResponse, error) {
	count, results, err := fetchAllPages[Set](c.http, fmt.Sprintf("/lego/sets/%s/alternates/", setNum))
	if err != nil {
		return nil, fmt.Errorf("get lego set alternates: %w", err)
	}
	return &LegoSetsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoSetMinifigs(setNum string) (*SetMinifigsResponse, error) {
	count, results, err := fetchAllPages[SetMinifig](c.http, fmt.Sprintf("/lego/sets/%s/minifigs/", setNum))
	if err != nil {
		return nil, fmt.Errorf("get lego set minifigs: %w", err)
	}
	return &SetMinifigsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoSetParts(setNum string) (*SetPartsResponse, error) {
	count, results, err := fetchAllPages[SetPart](c.http, fmt.Sprintf("/lego/sets/%s/parts/", setNum))
	if err != nil {
		return nil, fmt.Errorf("get lego set parts: %w", err)
	}
	return &SetPartsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoSetSets(setNum string) (*LegoSetsResponse, error) {
	count, results, err := fetchAllPages[Set](c.http, fmt.Sprintf("/lego/sets/%s/sets/", setNum))
	if err != nil {
		return nil, fmt.Errorf("get lego set sets: %w", err)
	}
	return &LegoSetsResponse{Count: count, Results: results}, nil
}

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

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./...
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
git add lego.go client_test.go
git commit -m "feat: add LEGO catalog endpoint methods"
```

---

### Task 5: `lego_parts.go` — Parts Methods + Tests

**Files:**
- Create: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/lego_parts.go`
- Modify: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/client_test.go` (append parts tests)

Key change from source: `GetLegoPart` and `GetLegoPartColor` previously printed to stdout on 404 and returned `nil, nil`. In the library they return an error on 404, consistent with all other GET methods.

- [ ] **Step 1: Write failing tests for lego parts methods**

Append to `client_test.go`:

```go
func TestGetLegoParts(t *testing.T) {
	tests := []struct {
		name       string
		response   LegoPartsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns parts", LegoPartsResponse{Count: 1, Results: []PartDetail{{PartNum: "3001", Name: "Brick 2 x 4"}}}, 200, false},
		{"server error", LegoPartsResponse{}, 500, true},
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
			result, err := client.GetLegoParts(PartsFilter{})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoParts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoParts() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoPartsAppliesFilters(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(LegoPartsResponse{Count: 0})
	}))
	defer server.Close()

	client := newClientWithBaseURL("key", "", server.URL)
	_, err := client.GetLegoParts(PartsFilter{PartCatID: "5", ColorID: "4", Search: "brick"})
	if err != nil {
		t.Fatalf("GetLegoParts() unexpected error: %v", err)
	}

	for _, want := range []string{"part_cat_id=5", "color_id=4", "search=brick"} {
		if !strings.Contains(capturedQuery, want) {
			t.Errorf("GetLegoParts() query = %q, missing %q", capturedQuery, want)
		}
	}
}

func TestGetLegoPartsPagination(t *testing.T) {
	page1 := PartDetail{PartNum: "3001"}
	page2 := PartDetail{PartNum: "3002"}

	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if r.URL.Query().Get("page") == "2" {
			_ = json.NewEncoder(w).Encode(LegoPartsResponse{Count: 2, Results: []PartDetail{page2}})
		} else {
			_ = json.NewEncoder(w).Encode(LegoPartsResponse{Count: 2, Next: serverURL + "/?page=2", Results: []PartDetail{page1}})
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := newClientWithBaseURL("key", "", server.URL)
	result, err := client.GetLegoParts(PartsFilter{})
	if err != nil {
		t.Fatalf("GetLegoParts() unexpected error: %v", err)
	}
	if len(result.Results) != 2 {
		t.Fatalf("GetLegoParts() len = %d, want 2", len(result.Results))
	}
	if result.Results[0].PartNum != page1.PartNum || result.Results[1].PartNum != page2.PartNum {
		t.Errorf("GetLegoParts() pagination order wrong: %+v", result.Results)
	}
}

func TestGetLegoPart(t *testing.T) {
	tests := []struct {
		name       string
		response   PartDetail
		statusCode int
		wantErr    bool
	}{
		{"returns part", PartDetail{PartNum: "3001", Name: "Brick 2 x 4"}, 200, false},
		{"not found", PartDetail{}, 404, true},
		{"server error", PartDetail{}, 500, true},
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
			result, err := client.GetLegoPart("3001")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoPart() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.PartNum != tt.response.PartNum {
				t.Errorf("GetLegoPart() part_num = %v, want %v", result.PartNum, tt.response.PartNum)
			}
		})
	}
}

func TestGetLegoPartColors(t *testing.T) {
	tests := []struct {
		name       string
		response   PartColorsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns colors", PartColorsResponse{Count: 1, Results: []PartColorDetail{{ColorID: 4, ColorName: "Red"}}}, 200, false},
		{"server error", PartColorsResponse{}, 500, true},
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
			result, err := client.GetLegoPartColors("3001")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoPartColors() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoPartColors() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoPartColorsPagination(t *testing.T) {
	page1 := PartColorDetail{ColorID: 4, ColorName: "Red"}
	page2 := PartColorDetail{ColorID: 5, ColorName: "Blue"}

	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if r.URL.Query().Get("page") == "2" {
			_ = json.NewEncoder(w).Encode(PartColorsResponse{Count: 2, Results: []PartColorDetail{page2}})
		} else {
			_ = json.NewEncoder(w).Encode(PartColorsResponse{Count: 2, Next: serverURL + "/?page=2", Results: []PartColorDetail{page1}})
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := newClientWithBaseURL("key", "", server.URL)
	result, err := client.GetLegoPartColors("3001")
	if err != nil {
		t.Fatalf("GetLegoPartColors() unexpected error: %v", err)
	}
	if len(result.Results) != 2 {
		t.Fatalf("GetLegoPartColors() len = %d, want 2", len(result.Results))
	}
}

func TestGetLegoPartColor(t *testing.T) {
	tests := []struct {
		name       string
		response   PartColorDetail
		statusCode int
		wantErr    bool
	}{
		{"returns combination", PartColorDetail{ColorID: 4, ColorName: "Red", NumSets: 12}, 200, false},
		{"not found", PartColorDetail{}, 404, true},
		{"server error", PartColorDetail{}, 500, true},
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
			result, err := client.GetLegoPartColor("3001", "4")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoPartColor() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.ColorName != tt.response.ColorName {
				t.Errorf("GetLegoPartColor() color_name = %v, want %v", result.ColorName, tt.response.ColorName)
			}
		})
	}
}

func TestGetLegoPartColorSets(t *testing.T) {
	tests := []struct {
		name       string
		response   LegoSetsResponse
		statusCode int
		wantErr    bool
	}{
		{"returns sets", LegoSetsResponse{Count: 1, Results: []Set{{SetNum: "10497-1"}}}, 200, false},
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
			result, err := client.GetLegoPartColorSets("3001", "4")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLegoPartColorSets() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetLegoPartColorSets() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetLegoPartColorSetsPagination(t *testing.T) {
	page1 := Set{SetNum: "10497-1"}
	page2 := Set{SetNum: "75192-1"}

	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if r.URL.Query().Get("page") == "2" {
			_ = json.NewEncoder(w).Encode(LegoSetsResponse{Count: 2, Results: []Set{page2}})
		} else {
			_ = json.NewEncoder(w).Encode(LegoSetsResponse{Count: 2, Next: serverURL + "/?page=2", Results: []Set{page1}})
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := newClientWithBaseURL("key", "", server.URL)
	result, err := client.GetLegoPartColorSets("3001", "4")
	if err != nil {
		t.Fatalf("GetLegoPartColorSets() unexpected error: %v", err)
	}
	if len(result.Results) != 2 {
		t.Fatalf("GetLegoPartColorSets() len = %d, want 2", len(result.Results))
	}
}
```

Also add `"strings"` to the imports in `client_test.go` (needed for `TestGetLegoPartsAppliesFilters`).

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./...
```

Expected: FAIL — `PartDetail`, `PartsFilter`, etc. undefined.

- [ ] **Step 3: Create `lego_parts.go`**

Copy from `rebrickable-cli/cli/cmd/api/lego_parts.go` with these changes:
- `package api` → `package rebrickable`
- Remove `fmt.Printf` from `GetLegoPart` 404 case; return error instead
- Remove `fmt.Printf` from `GetLegoPartColor` 404 case; return error instead
- Remove `fmt` from imports if only used for the removed prints (keep it — `fmt.Sprintf` and `fmt.Errorf` still used)

```go
package rebrickable

import (
	"fmt"
	"net/url"
)

type PartDetail struct {
	PartNum    string           `json:"part_num"`
	Name       string           `json:"name"`
	PartCatID  int              `json:"part_cat_id"`
	PartURL    string           `json:"part_url"`
	PartImgURL string           `json:"part_img_url"`
	ExternalIDs map[string][]any `json:"external_ids,omitempty"`
	PrintOf    string           `json:"print_of,omitempty"`
	YearFrom   int              `json:"year_from,omitempty"`
	YearTo     int              `json:"year_to,omitempty"`
}

type PartColorDetail struct {
	ColorID     int      `json:"color_id"`
	ColorName   string   `json:"color_name"`
	NumSets     int      `json:"num_sets"`
	NumSetParts int      `json:"num_set_parts"`
	PartImgURL  string   `json:"part_img_url"`
	Elements    []string `json:"elements,omitempty"`
}

type LegoPartsResponse struct {
	Count    int          `json:"count"`
	Next     string       `json:"next"`
	Previous string       `json:"previous"`
	Results  []PartDetail `json:"results"`
}

type PartColorsResponse struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []PartColorDetail `json:"results"`
}

type PartsFilter struct {
	PartNum     string
	PartNums    string
	PartCatID   string
	ColorID     string
	BricklinkID string
	BrickowlID  string
	LegoID      string
	LdrawID     string
	Ordering    string
	Search      string
}

func (f PartsFilter) queryString() string {
	q := url.Values{}
	pairs := []struct {
		key string
		val string
	}{
		{"part_num", f.PartNum},
		{"part_nums", f.PartNums},
		{"part_cat_id", f.PartCatID},
		{"color_id", f.ColorID},
		{"bricklink_id", f.BricklinkID},
		{"brickowl_id", f.BrickowlID},
		{"lego_id", f.LegoID},
		{"ldraw_id", f.LdrawID},
		{"ordering", f.Ordering},
		{"search", f.Search},
	}
	for _, p := range pairs {
		if p.val != "" {
			q.Set(p.key, p.val)
		}
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

func (c *Client) GetLegoParts(filter PartsFilter) (*LegoPartsResponse, error) {
	count, results, err := fetchAllPages[PartDetail](c.http, "/lego/parts/"+filter.queryString())
	if err != nil {
		return nil, fmt.Errorf("get lego parts: %w", err)
	}
	return &LegoPartsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoPart(partNum string) (*PartDetail, error) {
	result := &PartDetail{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/parts/%s/", partNum))

	if err != nil {
		return nil, fmt.Errorf("get lego part request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego part failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) GetLegoPartColors(partNum string) (*PartColorsResponse, error) {
	count, results, err := fetchAllPages[PartColorDetail](c.http, fmt.Sprintf("/lego/parts/%s/colors/", partNum))
	if err != nil {
		return nil, fmt.Errorf("get lego part colors: %w", err)
	}
	return &PartColorsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetLegoPartColor(partNum, colorID string) (*PartColorDetail, error) {
	result := &PartColorDetail{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/parts/%s/colors/%s/", partNum, colorID))

	if err != nil {
		return nil, fmt.Errorf("get lego part color request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego part color failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) GetLegoPartColorSets(partNum, colorID string) (*LegoSetsResponse, error) {
	count, results, err := fetchAllPages[Set](c.http, fmt.Sprintf("/lego/parts/%s/colors/%s/sets/", partNum, colorID))
	if err != nil {
		return nil, fmt.Errorf("get lego part color sets: %w", err)
	}
	return &LegoSetsResponse{Count: count, Results: results}, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./...
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
git add lego_parts.go client_test.go
git commit -m "feat: add LEGO parts endpoint methods"
```

---

### Task 6: `user.go` — User Endpoint Methods + Tests

**Files:**
- Create: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/user.go`
- Modify: `/Users/fabio.mangione/workspace/my-stuff/rebrickable-go/client_test.go` (append user tests)

Key changes from source: all `fmt.Println` and `fmt.Printf` calls removed. Delete methods (404 → `nil`, idempotent). `fmt` import stays (still used by `fmt.Errorf` and `fmt.Sprintf`).

- [ ] **Step 1: Write failing tests for user methods**

Append to `client_test.go`:

```go
func TestStoreUserSetList(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"created successfully", 201, false},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.StoreUserSetList("My List")
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreUserSetList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserSetLists(t *testing.T) {
	tests := []struct {
		name       string
		response   SetListsResponse
		statusCode int
		wantErr    bool
	}{
		{
			"returns set lists",
			SetListsResponse{Count: 1, Results: []SetList{{ID: 42, Name: "Technic"}}},
			200,
			false,
		},
		{"server error", SetListsResponse{}, 500, true},
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

			client := newClientWithBaseURL("key", "token", server.URL)
			result, err := client.GetUserSetLists()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSetLists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetUserSetLists() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetUserSetList(t *testing.T) {
	tests := []struct {
		name       string
		response   SetList
		statusCode int
		wantErr    bool
	}{
		{"returns set list", SetList{ID: 42, Name: "Technic"}, 200, false},
		{"not found", SetList{}, 404, true},
		{"server error", SetList{}, 500, true},
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

			client := newClientWithBaseURL("key", "token", server.URL)
			result, err := client.GetUserSetList("42")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSetList() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Name != tt.response.Name {
				t.Errorf("GetUserSetList() name = %v, want %v", result.Name, tt.response.Name)
			}
		})
	}
}

func TestUpdateUserSetList(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"updated successfully", 200, false},
		{"not found", 404, true},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.UpdateUserSetList("42", "New Name")
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserSetList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReplaceUserSetList(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"replaced successfully", 200, false},
		{"not found", 404, true},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.ReplaceUserSetList("42", "New Name")
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceUserSetList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteUserSetList(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"deleted successfully", 204, false},
		{"not found", 404, false},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.DeleteUserSetList("123")
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserSetList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserSetListSets(t *testing.T) {
	tests := []struct {
		name       string
		response   SetsResponse
		statusCode int
		wantErr    bool
	}{
		{
			"returns sets in set list",
			SetsResponse{Count: 1, Results: []UserSet{{Quantity: 1}}},
			200,
			false,
		},
		{"server error", SetsResponse{}, 500, true},
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

			client := newClientWithBaseURL("key", "token", server.URL)
			result, err := client.GetUserSetListSets("42")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSetListSets() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetUserSetListSets() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestStoreUserSetListSet(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"created successfully", 201, false},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.StoreUserSetListSet("42", "10274-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreUserSetListSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserSetListSet(t *testing.T) {
	tests := []struct {
		name       string
		response   UserSet
		statusCode int
		wantErr    bool
	}{
		{"returns set", UserSet{Quantity: 2}, 200, false},
		{"not found", UserSet{}, 404, true},
		{"server error", UserSet{}, 500, true},
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

			client := newClientWithBaseURL("key", "token", server.URL)
			result, err := client.GetUserSetListSet("42", "10274-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSetListSet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Quantity != tt.response.Quantity {
				t.Errorf("GetUserSetListSet() quantity = %v, want %v", result.Quantity, tt.response.Quantity)
			}
		})
	}
}

func TestDeleteUserSetListSet(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"deleted successfully", 204, false},
		{"not found", 404, false},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.DeleteUserSetListSet("42", "10274-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserSetListSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoreUserSet(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"created successfully", 201, false},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.StoreUserSet("42043-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreUserSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserSets(t *testing.T) {
	tests := []struct {
		name       string
		response   SetsResponse
		statusCode int
		wantErr    bool
	}{
		{
			"returns user sets",
			SetsResponse{Count: 2, Results: []UserSet{{Quantity: 1}, {Quantity: 3}}},
			200,
			false,
		},
		{"server error", SetsResponse{}, 500, true},
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

			client := newClientWithBaseURL("key", "token", server.URL)
			result, err := client.GetUserSets()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSets() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Count != tt.response.Count {
				t.Errorf("GetUserSets() count = %v, want %v", result.Count, tt.response.Count)
			}
		})
	}
}

func TestGetUserSet(t *testing.T) {
	tests := []struct {
		name       string
		response   UserSet
		statusCode int
		wantErr    bool
	}{
		{"returns set", UserSet{Quantity: 1, Set: Set{SetNum: "10274-1", Name: "Ghost"}}, 200, false},
		{"not found", UserSet{}, 404, true},
		{"server error", UserSet{}, 500, true},
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

			client := newClientWithBaseURL("key", "token", server.URL)
			result, err := client.GetUserSet("10274-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserSet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result.Set.SetNum != tt.response.Set.SetNum {
				t.Errorf("GetUserSet() set_num = %v, want %v", result.Set.SetNum, tt.response.Set.SetNum)
			}
		})
	}
}

func TestReplaceUserSet(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"replaced successfully", 200, false},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.ReplaceUserSet("10274-1", 2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceUserSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteUserSet(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"deleted successfully", 204, false},
		{"not found", 404, false},
		{"server error", 500, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := newClientWithBaseURL("key", "token", server.URL)
			err := client.DeleteUserSet("42043-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUserSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserSetsPagination(t *testing.T) {
	page1Set := UserSet{Quantity: 1, Set: Set{SetNum: "10497-1"}}
	page2Set := UserSet{Quantity: 2, Set: Set{SetNum: "75192-1"}}

	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if r.URL.Query().Get("page") == "2" {
			_ = json.NewEncoder(w).Encode(SetsResponse{Count: 2, Results: []UserSet{page2Set}})
		} else {
			_ = json.NewEncoder(w).Encode(SetsResponse{Count: 2, Next: serverURL + "/?page=2", Results: []UserSet{page1Set}})
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := newClientWithBaseURL("key", "token", server.URL)
	result, err := client.GetUserSets()

	if err != nil {
		t.Fatalf("GetUserSets() unexpected error: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("GetUserSets() count = %d, want 2", result.Count)
	}
	if len(result.Results) != 2 {
		t.Errorf("GetUserSets() len(results) = %d, want 2", len(result.Results))
	}
	if result.Results[0].Set.SetNum != page1Set.Set.SetNum {
		t.Errorf("GetUserSets() results[0].SetNum = %q, want %q", result.Results[0].Set.SetNum, page1Set.Set.SetNum)
	}
	if result.Results[1].Set.SetNum != page2Set.Set.SetNum {
		t.Errorf("GetUserSets() results[1].SetNum = %q, want %q", result.Results[1].Set.SetNum, page2Set.Set.SetNum)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./...
```

Expected: FAIL — user methods undefined.

- [ ] **Step 3: Create `user.go`**

Copy from `rebrickable-cli/cli/cmd/api/user.go` with these changes:
- `package api` → `package rebrickable`
- Remove all `fmt.Println` and `fmt.Printf` calls
- `fmt` import kept (still used by `fmt.Errorf`, `fmt.Sprintf`)

```go
package rebrickable

import "fmt"

func (c *Client) StoreUserSetList(name string) error {
	resp, err := c.http.R().
		SetBody(map[string]string{"name": name}).
		Post(c.userPath("/setlists/"))

	if err != nil {
		return fmt.Errorf("store set list request failed: %w", err)
	}
	if resp.StatusCode() != 201 {
		return fmt.Errorf("store set list failed with status %d", resp.StatusCode())
	}
	return nil
}

func (c *Client) GetUserSetLists() (*SetListsResponse, error) {
	count, results, err := fetchAllPages[SetList](c.http, c.userPath("/setlists"))
	if err != nil {
		return nil, fmt.Errorf("get user set lists: %w", err)
	}
	return &SetListsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetUserSetList(listID string) (*SetList, error) {
	result := &SetList{}
	resp, err := c.http.R().
		SetResult(result).
		Get(c.userPath(fmt.Sprintf("/setlists/%s/", listID)))

	if err != nil {
		return nil, fmt.Errorf("get set list request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get set list failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) UpdateUserSetList(listID, name string) error {
	resp, err := c.http.R().
		SetBody(map[string]string{"name": name}).
		Patch(c.userPath(fmt.Sprintf("/setlists/%s/", listID)))

	if err != nil {
		return fmt.Errorf("update set list request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("update set list failed with status %d", resp.StatusCode())
	}
	return nil
}

func (c *Client) ReplaceUserSetList(listID, name string) error {
	resp, err := c.http.R().
		SetBody(map[string]string{"name": name}).
		Put(c.userPath(fmt.Sprintf("/setlists/%s/", listID)))

	if err != nil {
		return fmt.Errorf("replace set list request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("replace set list failed with status %d", resp.StatusCode())
	}
	return nil
}

func (c *Client) DeleteUserSetList(id string) error {
	resp, err := c.http.R().
		Delete(c.userPath(fmt.Sprintf("/setlists/%s/", id)))

	if err != nil {
		return fmt.Errorf("delete set list request failed: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("delete set list failed with status %d", resp.StatusCode())
	}
	return nil
}

func (c *Client) GetUserSetListSets(listID string) (*SetsResponse, error) {
	count, results, err := fetchAllPages[UserSet](c.http, c.userPath(fmt.Sprintf("/setlists/%s/sets/", listID)))
	if err != nil {
		return nil, fmt.Errorf("get user set list sets: %w", err)
	}
	return &SetsResponse{Count: count, Results: results}, nil
}

func (c *Client) StoreUserSetListSet(listID, setNum string) error {
	resp, err := c.http.R().
		SetBody(map[string]string{"set_num": setNum, "quantity": "1"}).
		Post(c.userPath(fmt.Sprintf("/setlists/%s/sets/", listID)))

	if err != nil {
		return fmt.Errorf("store set list set request failed: %w", err)
	}
	if resp.StatusCode() != 201 {
		return fmt.Errorf("store set list set failed with status %d", resp.StatusCode())
	}
	return nil
}

func (c *Client) GetUserSetListSet(listID, setNum string) (*UserSet, error) {
	result := &UserSet{}
	resp, err := c.http.R().
		SetResult(result).
		Get(c.userPath(fmt.Sprintf("/setlists/%s/sets/%s/", listID, setNum)))

	if err != nil {
		return nil, fmt.Errorf("get set list set request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get set list set failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) DeleteUserSetListSet(listID, setNum string) error {
	resp, err := c.http.R().
		Delete(c.userPath(fmt.Sprintf("/setlists/%s/sets/%s/", listID, setNum)))

	if err != nil {
		return fmt.Errorf("delete set list set request failed: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("delete set list set failed with status %d", resp.StatusCode())
	}
	return nil
}

func (c *Client) StoreUserSet(setNumber string) error {
	resp, err := c.http.R().
		SetBody(map[string]string{"set_num": setNumber, "quantity": "1"}).
		Post(c.userPath("/sets/"))

	if err != nil {
		return fmt.Errorf("store set request failed: %w", err)
	}
	if resp.StatusCode() != 201 {
		return fmt.Errorf("store set failed with status %d", resp.StatusCode())
	}
	return nil
}

func (c *Client) GetUserSets() (*SetsResponse, error) {
	count, results, err := fetchAllPages[UserSet](c.http, c.userPath("/sets"))
	if err != nil {
		return nil, fmt.Errorf("get user sets: %w", err)
	}
	return &SetsResponse{Count: count, Results: results}, nil
}

func (c *Client) GetUserSet(setNum string) (*UserSet, error) {
	result := &UserSet{}
	resp, err := c.http.R().
		SetResult(result).
		Get(c.userPath(fmt.Sprintf("/sets/%s/", setNum)))

	if err != nil {
		return nil, fmt.Errorf("get set request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get set failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) ReplaceUserSet(setNum string, quantity int) error {
	resp, err := c.http.R().
		SetBody(map[string]int{"quantity": quantity}).
		Put(c.userPath(fmt.Sprintf("/sets/%s/", setNum)))

	if err != nil {
		return fmt.Errorf("replace set request failed: %w", err)
	}
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
		return fmt.Errorf("replace set failed with status %d", resp.StatusCode())
	}
	return nil
}

func (c *Client) DeleteUserSet(setNumber string) error {
	path := c.userPath(fmt.Sprintf("/sets/%s/", setNumber))
	resp, err := c.http.R().Delete(path)

	if err != nil {
		return fmt.Errorf("delete set request failed: %w", err)
	}
	if resp.StatusCode() == 404 {
		return nil
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("delete set failed with status %d", resp.StatusCode())
	}
	return nil
}
```

- [ ] **Step 4: Run all tests to verify they pass**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./...
```

Expected: PASS — all tests green.

- [ ] **Step 5: Commit**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
git add user.go client_test.go
git commit -m "feat: add user endpoint methods"
```

---

### Task 7: Final Verification and Tag

**Files:** none new

- [ ] **Step 1: Run all tests with verbose output**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go test ./... -v
```

Expected: All tests PASS. Count: should be ~45 test functions.

- [ ] **Step 2: Verify no stdout output from library methods**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
grep -rn "fmt\.Print" .
```

Expected: no output (no fmt.Print, fmt.Printf, fmt.Println in any non-test file).

- [ ] **Step 3: Verify build is clean**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
go build ./...
go vet ./...
```

Expected: no errors, no output.

- [ ] **Step 4: Tag initial version**

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
git tag v0.1.0
```

- [ ] **Step 5: Final commit (if anything changed)**

If `go vet` or `go build` revealed anything, fix and commit. Otherwise:

```bash
cd /Users/fabio.mangione/workspace/my-stuff/rebrickable-go
git log --oneline
```

Expected output (5 commits):
```
<sha> feat: add user endpoint methods
<sha> feat: add LEGO parts endpoint methods
<sha> feat: add LEGO catalog endpoint methods
<sha> feat: add Client struct, constructors, and fetchAllPages helper
<sha> feat: add domain types
<sha> chore: initialize rebrickable-go module
```
