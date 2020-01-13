package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Schedules struct {
	Schedules []*Schedule       `json:"schedules"`
	Limit     int               `json:"limit,omitempty"`
	Offset    int               `json:"offset,omitempty"`
	More      bool              `json:"more,omitempty"`
	Users     map[string]string `json:"users,omitempty"`
}

type Schedule struct {
	ID          string  `json:"id"`
	Name        string  `json:"summary"`
	Description string  `json:"description,omitempty"`
	Users       []*User `json:"users,omitempty"`
}

func (c *Client) Schedules() (*Schedules, error) {
	c.BaseURL.Path = "/schedules"
	var offset int
	q := url.Values{}

	var schedules Schedules

	for {
		q.Set("offset", strconv.Itoa(offset))
		c.BaseURL.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", c.BaseURL.String(), nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
		req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.Token))
		req.Header.Set("User-Agent", c.UserAgent)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		respBody := bytes.NewBuffer(make([]byte, 0))
		for {
			_, err := io.CopyN(respBody, resp.Body, 1024)
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, fmt.Errorf("error while reading body: %v", err)
			}
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			log.Println("recieved 429 from API")
			time.Sleep(time.Second)
			continue
		} else if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("recieved response from API with %d status code: %s", resp.StatusCode, respBody.String())
		}

		var tmp Schedules
		if err = json.Unmarshal(respBody.Bytes(), &tmp); err != nil {
			return nil, err
		}
		schedules.Schedules = append(schedules.Schedules, tmp.Schedules...)

		if !tmp.More {
			break
		}

		offset += tmp.Limit
	}

	return &schedules, nil
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
	ID      string `json:"id"`
	Name    string `json:"summary"`
	Deleted string `json:"deleted_at,omitempty"`
}

func (c *Client) Schedule(id, startdate, enddate string) (*PDScheduleResponse, error) {
	path := &url.URL{Path: fmt.Sprint("/schedules/", id)}
	u := c.BaseURL.ResolveReference(path)
	q := url.Values{"include_oncall": []string{"true"}}
	if startdate != "" {
		q.Set("since", startdate)
	}
	if enddate != "" {
		q.Set("until", enddate)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.Token))
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := bytes.NewBuffer(make([]byte, 0))
	for {
		_, err := io.CopyN(respBody, resp.Body, 1024)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error while reading body: %v", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recieved response from API with %d status code: %s", resp.StatusCode, respBody.String())
	}

	var r PDScheduleResponse
	if err = json.Unmarshal(respBody.Bytes(), &r); err != nil {
		return nil, err
	}

	return &r, nil
}
