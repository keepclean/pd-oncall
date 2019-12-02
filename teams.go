package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type PDTeams struct {
	Teams  []*PDTeam `json:"teams"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
	More   bool      `json:"more"`
}

type PDTeam struct {
	ID   string `json:"id"`
	Name string `json:"summary"`
}

func (c *Client) Teams() ([]*PDTeam, error) {
	c.BaseURL.Path = "/teams"
	q := c.BaseURL.Query()
	var offset int
	q.Set("offset", strconv.Itoa(offset))
	c.BaseURL.RawQuery = q.Encode()

	var teams PDTeams

	req, err := http.NewRequest("GET", c.BaseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.Token))
	req.Header.Set("User-Agent", c.UserAgent)

	for {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var tmp PDTeams
		if err = json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
			return nil, err
		}
		teams.Teams = append(teams.Teams, tmp.Teams...)

		if !tmp.More {
			break
		}

		offset += tmp.Limit
		q.Set("offset", strconv.Itoa(offset))
		c.BaseURL.RawQuery = q.Encode()
		req.URL = c.BaseURL
	}

	return teams.Teams, nil
}
