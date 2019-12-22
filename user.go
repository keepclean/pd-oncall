package main

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/jedib0t/go-pretty/table"
)

func oncallUser(apiClient *Client, cf *Schedules, since, until, tableStyle string) {
	userIDs := printUserAsTable(cf.Users)
	index, err := getUserInput("Please select an user id: ")
	if err != nil {
		log.Println("Failed to get user ID:", err)
		return
	}

	id := index[0]
	if id < 0 || id > len(cf.Users) {
		log.Println("Specified user ID is out of range:", id)
		return
	}
	userID := userIDs[id-1]
	userName := cf.Users[userID]

	var data []table.Row
	fields := table.Row{"START", "DAY", "SHIFT"}
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup

	for _, shift := range cf.Schedules {
		wg.Add(1)
		go func(id, name, t1, t2 string) {
			defer wg.Done()
			schedule, err := apiClient.Schedule(id, t1, t2)
			if err != nil {
				log.Printf("Failed to get schedule for %s: %v\n", name, err)
				return
			}

			for _, entry := range schedule.Schedule.FinalSchedule.RenderedScheduleEntries {
				if entry.User.ID != userID {
					continue
				}

				start, err := convertTime(entry.Start, "")
				if err != nil {
					start = entry.Start
				}
				day, _ := weekday(entry.Start)

				mutex.Lock()
				data = append(data, table.Row{start, day, name})
				mutex.Unlock()
			}
		}(shift.ID, shift.Name, since, until)
	}
	wg.Wait()

	sort.Slice(data, func(i, j int) bool { return data[i][0].(string) < data[j][0].(string) })

	title := fmt.Sprintf("%s (%s - %s)", userName, since, until)
	printTable(data, fields, title, tableStyle, "", -1)
}
