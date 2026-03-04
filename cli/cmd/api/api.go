package api

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

const apiBaseURI = "https://rebrickable.com/api/v3"

func GetURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return apiBaseURI + path
}

type Client struct {
	http      *resty.Client
	authToken string
}

func NewClient(apiKey, authToken string) *Client {
	return newClientWithBaseURL(apiKey, authToken, apiBaseURI)
}

func NewLegoClient(apiKey string) *Client {
	return newClientWithBaseURL(apiKey, "", apiBaseURI)
}

func newClientWithBaseURL(apiKey, authToken, baseURL string) *Client {
	http := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("key %s", apiKey))

	return &Client{http: http, authToken: authToken}
}

func (c *Client) userPath(path string) string {
	return fmt.Sprintf("/users/%s%s", c.authToken, path)
}

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
	ID   int    `json:"id"`
	Name string `json:"name"`
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
