package main

import (
	"log"
	"sort"
	"sync"

	"github.com/jedib0t/go-pretty/table"
)

func oncallNow(apiClient *Client, cf *Schedules, tableStyle string) {
	var data []table.Row
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	for _, shift := range cf.Schedules {
		wg.Add(1)
		go func(ID, Name string) {
			defer wg.Done()
			schedule, err := apiClient.Schedule(ID, "", "")
			if err != nil {
				log.Printf("Failed to get schedule for %s: %v\n", Name, err)
				return
			}

			mutex.Lock()
			data = append(data, table.Row{Name, schedule.Schedule.Oncall.User.Name})
			mutex.Unlock()

		}(shift.ID, shift.Name)
	}
	wg.Wait()

	sort.Slice(data, func(i, j int) bool { return data[i][0].(string) < data[j][0].(string) })

	fields := table.Row{"SHIFT", "ENGINEER"}
	printTable(data, fields, "", tableStyle)
}
