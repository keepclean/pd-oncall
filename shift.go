package main

import (
	"log"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

func oncallShift(apiClient *Client, cf *Schedules, since, until, tableStyle string) {
	printSchedulesAsTable(cf.Schedules)
	shift, err := getUserInput("Please enter shift ID: ")
	if err != nil {
		log.Println("Failed to get shift ID:", err)
		return
	}

	index := shift[0]
	if index < 0 || index > len(cf.Schedules) {
		log.Println("Specified shift ID is out of range:", index)
		return
	}

	ID := cf.Schedules[index-1].ID
	Name := cf.Schedules[index-1].Name

	schedule, err := apiClient.Schedule(ID, since, until)
	if err != nil {
		log.Println("Failed to get schedule:", err)
		return
	}

	data := make([]table.Row, 0)
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
	printTable(data, fields, Name, tableStyle, "")
}
