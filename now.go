package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/jedib0t/go-pretty/table"
)

func oncallNow(apiClient *Client, cf *PDTeams) {
	var data []table.Row
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	for _, shift := range cf.Teams {
		wg.Add(1)
		go func(sName, sID string) {
			defer wg.Done()
			_, err := apiClient.finalSchedule(sID, "", "")
			if err != nil {
				log.Printf("Failed to get schedule for %s: %v\n", sName, err)
				return
			}

			mutex.Lock()
			data = append(data, table.Row{sName, ""})
			mutex.Unlock()

		}(shift.ID, shift.Name)
	}
	wg.Wait()
}

func (c *Client) finalSchedule(ID, startdate, enddate string) ([]*PDTeam, error) {
	c.BaseURL.Path = "/schedules"
	q := c.BaseURL.Query()
	q.Set("include_oncall", "true")
	if startdate != "" {
		q.Set("since", startdate)
	}
	if enddate != "" {
		q.Set("until", enddate)
	}
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}
	defer resp.Body.Close()

	var tmp PDTeams
	if err = json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
		return nil, err
	}

	return tmp.Teams, nil
}
