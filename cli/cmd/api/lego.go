package api

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
