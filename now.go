package main

import (
	"fmt"
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
		go func(sID, sName string) {
			defer wg.Done()
			schedule, err := apiClient.finalSchedule(sID, "", "")
			if err != nil {
				log.Printf("Failed to get schedule for %s: %v\n", sName, err)
				return
			}

			fmt.Printf("%s - %s: %v\n", sID, sName, schedule)
			mutex.Lock()
			data = append(data, table.Row{sName, schedule.Schedule.Oncall.EntryUser.Name})
			mutex.Unlock()

		}(shift.ID, shift.Name)
	}
	wg.Wait()

	sort.Slice(data, func(i, j int) bool { return data[i][0].(string) < data[j][0].(string) })

	fields := table.Row{"SHIFT", "ENGINEER"}
	printTable(data, fields, tableStyle)
}
