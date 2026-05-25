package api

import (
	"fmt"
	"net/url"
)

type PartDetail struct {
	PartNum         string            `json:"part_num"`
	Name            string            `json:"name"`
	PartCatID       int               `json:"part_cat_id"`
	PartURL         string            `json:"part_url"`
	PartImgURL      string            `json:"part_img_url"`
	ExternalIDs     map[string][]any  `json:"external_ids,omitempty"`
	PrintOf         string            `json:"print_of,omitempty"`
	YearFrom        int               `json:"year_from,omitempty"`
	YearTo          int               `json:"year_to,omitempty"`
}

type PartColorDetail struct {
	ColorID     int    `json:"color_id"`
	ColorName   string `json:"color_name"`
	NumSets     int    `json:"num_sets"`
	NumSetParts int    `json:"num_set_parts"`
	PartImgURL  string `json:"part_img_url"`
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
	if resp.StatusCode() == 404 {
		fmt.Printf("Part %s not found\n", partNum)
		return nil, nil
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
	if resp.StatusCode() == 404 {
		fmt.Printf("Part %s in color %s not found\n", partNum, colorID)
		return nil, nil
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
