package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/jedib0t/go-pretty/table"
)

func oncallRoster(apiClient *Client, cf *Schedules, since, until, tableStyle string) {
	fields := table.Row{"START", "DAY"}
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	rawData := make(map[string]map[string]string)

	for _, shift := range cf.Schedules {
		fields = append(fields, prettyField(shift.Name))

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

	data := make([]table.Row, 0, len(rawData))
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

func prettyField(f string) string {
	maxLineLenght := 20
	if len(f) < maxLineLenght {
		return f
	}

	var result, s []string
	var l int
	for _, item := range strings.Fields(f) {
		if item == "-" {
			continue
		}

		if l+len(item) > maxLineLenght {
			result = append(result, strings.Join(s, " "))
			s = []string{item}
			l = len(item)
			continue
		}

		s = append(s, item)
		l += len(item)
	}
	result = append(result, strings.Join(s, " "))

	return strings.Join(result, "\n")
}
