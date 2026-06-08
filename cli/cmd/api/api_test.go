package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

func TestGetURL(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"foo path", args{path: "/foo"}, "https://rebrickable.com/api/v3/foo"},
		{"foo path", args{path: "foo"}, "https://rebrickable.com/api/v3/foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetURL(tt.args.path); got != tt.want {
				t.Errorf("GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		wantNil    bool
		wantErr    bool
	}{
		{"returns part", PartDetail{PartNum: "3001", Name: "Brick 2 x 4"}, 200, false, false},
		{"not found", PartDetail{}, 404, true, false},
		{"server error", PartDetail{}, 500, true, true},
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
			if (result == nil) != tt.wantNil {
				t.Errorf("GetLegoPart() nilResult = %v, wantNil %v", result == nil, tt.wantNil)
			}
			if !tt.wantErr && !tt.wantNil && result.PartNum != tt.response.PartNum {
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
		wantNil    bool
		wantErr    bool
	}{
		{"returns combination", PartColorDetail{ColorID: 4, ColorName: "Red", NumSets: 12}, 200, false, false},
		{"not found", PartColorDetail{}, 404, true, false},
		{"server error", PartColorDetail{}, 500, true, true},
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
			if (result == nil) != tt.wantNil {
				t.Errorf("GetLegoPartColor() nilResult = %v, wantNil %v", result == nil, tt.wantNil)
			}
			if !tt.wantErr && !tt.wantNil && result.ColorName != tt.response.ColorName {
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
