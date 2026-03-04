package api

import "fmt"

func (c *Client) GetLegoSets() (*LegoSetsResponse, error) {
	result := &LegoSetsResponse{}
	resp, err := c.http.R().
		SetResult(result).
		Get("/lego/sets/")

	if err != nil {
		return nil, fmt.Errorf("get lego sets request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego sets failed with status %d", resp.StatusCode())
	}
	return result, nil
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
	result := &LegoSetsResponse{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/sets/%s/alternates/", setNum))

	if err != nil {
		return nil, fmt.Errorf("get lego set alternates request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego set alternates failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) GetLegoSetMinifigs(setNum string) (*SetMinifigsResponse, error) {
	result := &SetMinifigsResponse{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/sets/%s/minifigs/", setNum))

	if err != nil {
		return nil, fmt.Errorf("get lego set minifigs request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego set minifigs failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) GetLegoSetParts(setNum string) (*SetPartsResponse, error) {
	result := &SetPartsResponse{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/sets/%s/parts/", setNum))

	if err != nil {
		return nil, fmt.Errorf("get lego set parts request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego set parts failed with status %d", resp.StatusCode())
	}
	return result, nil
}

func (c *Client) GetLegoSetSets(setNum string) (*LegoSetsResponse, error) {
	result := &LegoSetsResponse{}
	resp, err := c.http.R().
		SetResult(result).
		Get(fmt.Sprintf("/lego/sets/%s/sets/", setNum))

	if err != nil {
		return nil, fmt.Errorf("get lego set sets request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get lego set sets failed with status %d", resp.StatusCode())
	}
	return result, nil
}
