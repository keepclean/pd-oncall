package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Schedules struct {
	Schedules []*Schedule `json:"schedules"`
	Limit     int         `json:"limit,omitempty"`
	Offset    int         `json:"offset,omitempty"`
	More      bool        `json:"more,omitempty"`
}

type Schedule struct {
	ID          string `json:"id"`
	Name        string `json:"summary"`
	Description string `json:"description"`
}

func (c *Client) Schedules() ([]*Schedule, error) {
	c.BaseURL.Path = "/schedules"
	q := c.BaseURL.Query()
	var offset int
	q.Set("offset", strconv.Itoa(offset))
	c.BaseURL.RawQuery = q.Encode()

	var schedules Schedules

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

		var tmp Schedules
		if err = json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
			return nil, err
		}
		schedules.Schedules = append(schedules.Schedules, tmp.Schedules...)

		if !tmp.More {
			break
		}

		offset += tmp.Limit
		q.Set("offset", strconv.Itoa(offset))
		c.BaseURL.RawQuery = q.Encode()
		req.URL = c.BaseURL
	}

	return schedules.Schedules, nil
}
