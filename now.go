package main

import (
	"log"
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

			var user string
			if schedule.Schedule.Oncall != nil {
				user = schedule.Schedule.Oncall.User.Name
			}

			mutex.Lock()
			data = append(data, table.Row{Name, user})
			mutex.Unlock()
		}(shift.ID, shift.Name)
	}
	wg.Wait()

	fields := table.Row{"SHIFT", "ENGINEER"}
	printTable(data, fields, "", tableStyle, "SHIFT")
}
