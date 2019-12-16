package main

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/jedib0t/go-pretty/table"
)

func oncallRoster(apiClient *Client, cf *Schedules, since, until, tableStyle string) {
	var data []table.Row
	fields := table.Row{"START", "DAY"}
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	rawData := make(map[string]map[string]string)

	for _, shift := range cf.Schedules {
		fields = append(fields, shift.Name)

		wg.Add(1)
		go func(id, name, t1, t2 string) {
			defer wg.Done()
			schedule, err := apiClient.Schedule(id, t1, t2)
			if err != nil {
				log.Printf("Failed to get schedule for %s: %v\n", name, err)
				return
			}

			for _, entry := range schedule.Schedule.FinalSchedule.RenderedScheduleEntries {
				start, err := convertTime(entry.Start, "")
				if err != nil {
					start = entry.Start
				}
				day, _ := weekday(entry.Start)

				mutex.Lock()
				if rawData[start] == nil {
					rawData[start] = map[string]string{}
				}

				rawData[start]["day"] = day.String()
				rawData[start][name] = entry.User.Name
				mutex.Unlock()
			}
		}(shift.ID, shift.Name, since, until)
	}
	wg.Wait()

	for k, v := range rawData {
		tmp := table.Row{k, v["day"]}
		for _, s := range cf.Schedules {
			tmp = append(tmp, v[s.Name])
		}

		data = append(data, tmp)
	}

	sort.Slice(data, func(i, j int) bool { return data[i][0].(string) < data[j][0].(string) })

	title := fmt.Sprintf("%s - %s", since, until)
	printTable(data, fields, title, tableStyle, "", -1)
}
