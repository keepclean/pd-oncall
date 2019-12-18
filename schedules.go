package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Schedules struct {
	Schedules []*Schedule `json:"schedules"`
	Limit     int         `json:"limit,omitempty"`
	Offset    int         `json:"offset,omitempty"`
	More      bool        `json:"more,omitempty"`
}

type Schedule struct {
	ID          string  `json:"id"`
	Name        string  `json:"summary"`
	Description string  `json:"description"`
	Teams       []*Team `json:"teams"`
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

type PDScheduleResponse struct {
	Schedule *PDScheduleItems `json:"schedule"`
}

type PDScheduleItems struct {
	FinalSchedule *FinalSchedule `json:"final_schedule"`
	Oncall        *OnCall        `json:"oncall"`
	Users         []*User        `json:"users"`
}

type FinalSchedule struct {
	RenderedScheduleEntries []*ScheduleEntry `json:"rendered_schedule_entries"`
}

type ScheduleEntry struct {
	Start string `json:"start"`
	End   string `json:"end"`
	User  *User  `json:"user"`
}

type OnCall struct {
	User *User `json:"user"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"summary"`
}

func (c *Client) Schedule(id, startdate, enddate string) (*PDScheduleResponse, error) {
	path := &url.URL{Path: fmt.Sprint("/schedules/", id)}
	u := c.BaseURL.ResolveReference(path)
	q := u.Query()
	q.Set("include_oncall", "true")
	if startdate != "" {
		q.Set("since", startdate)
	}
	if enddate != "" {
		q.Set("until", enddate)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return &PDScheduleResponse{}, err
	}

	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.Token))
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &PDScheduleResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &PDScheduleResponse{}, fmt.Errorf("%s", resp.Status)
	}
	defer resp.Body.Close()

	var r PDScheduleResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
