package main

import (
	"log"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

func oncallShift(apiClient *Client, cf *Schedules, since, until, tableStyle string) {
	printSchedulesAsTable(cf.Schedules)
	items, _ := getUserInput()
	ID := cf.Schedules[items[0]-1].ID
	Name := cf.Schedules[items[0]-1].Name

	schedule, err := apiClient.Schedule(ID, since, until)
	if err != nil {
		log.Fatalf("Failed to get schedule: %v", err)
	}

	var data []table.Row
	for _, entry := range schedule.Schedule.FinalSchedule.RenderedScheduleEntries {
		start, err := convertTime(entry.Start, "")
		if err != nil {
			start = entry.Start
		}
		day, _ := weekday(entry.Start)
		holiday, _ := holidays(entry.Start)

		data = append(data, table.Row{start, day.String(), entry.User.Name, strings.Join(holiday, ", ")})
	}

	fields := table.Row{"SINCE", "WEEKDAY", "ENGINEER", "HOLIDAY"}
	printTable(data, fields, Name, tableStyle)
}
