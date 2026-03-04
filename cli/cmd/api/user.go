package api

import (
	"fmt"
	"strings"
)

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
	fmt.Println("SetList saved")
	return nil
}

func (c *Client) GetUserSetLists() (*SetListsResponse, error) {
	result := &SetListsResponse{}
	resp, err := c.http.R().
		SetResult(result).
		Get(c.userPath("/setlists"))

	if err != nil {
		return nil, fmt.Errorf("get set lists request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get set lists failed with status %d", resp.StatusCode())
	}
	return result, nil
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
	fmt.Printf("Updated set list: %s\n", listID)
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
	fmt.Printf("Replaced set list: %s\n", listID)
	return nil
}

func (c *Client) DeleteUserSetList(id string) error {
	resp, err := c.http.R().
		Delete(c.userPath(fmt.Sprintf("/setlists/%s/", id)))

	if err != nil {
		return fmt.Errorf("delete set list request failed: %w", err)
	}
	if resp.StatusCode() == 404 {
		fmt.Printf("Set list %s not found\n", id)
		return nil
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("delete set list failed with status %d", resp.StatusCode())
	}
	fmt.Printf("Deleted set list: %s\n", id)
	return nil
}

func (c *Client) GetUserSetListSets(listID string) (*SetsResponse, error) {
	result := &SetsResponse{}
	resp, err := c.http.R().
		SetResult(result).
		Get(c.userPath(fmt.Sprintf("/setlists/%s/sets/", listID)))

	if err != nil {
		return nil, fmt.Errorf("get set list sets request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get set list sets failed with status %d", resp.StatusCode())
	}
	return result, nil
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
	fmt.Println("Set added to set list")
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
		fmt.Printf("Set %s not found in set list %s\n", setNum, listID)
		return nil
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("delete set list set failed with status %d", resp.StatusCode())
	}
	fmt.Printf("Deleted set %s from set list %s\n", setNum, listID)
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
	fmt.Println("Set saved")
	return nil
}

func (c *Client) GetUserSets() (*SetsResponse, error) {
	result := &SetsResponse{}
	resp, err := c.http.R().
		SetResult(result).
		Get(c.userPath("/sets"))

	if err != nil {
		return nil, fmt.Errorf("get sets request failed: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("get sets failed with status %d", resp.StatusCode())
	}
	return result, nil
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
	fmt.Printf("Updated set: %s\n", setNum)
	return nil
}

func (c *Client) DeleteUserSet(setNumber string) error {
	path := c.userPath(fmt.Sprintf("/sets/%s/", setNumber))
	fmt.Println("Calling URL:", strings.ReplaceAll(apiBaseURI+path, c.authToken, "#token#"))

	resp, err := c.http.R().Delete(path)

	if err != nil {
		return fmt.Errorf("delete set request failed: %w", err)
	}
	if resp.StatusCode() == 404 {
		fmt.Printf("Set %s not found\n", setNumber)
		return nil
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("delete set failed with status %d", resp.StatusCode())
	}
	fmt.Printf("Deleted set: %s\n", setNumber)
	return nil
}
