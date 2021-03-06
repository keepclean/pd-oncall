package main

import (
	"fmt"
	"log"

	"github.com/jedib0t/go-pretty/table"
)

func oncallReport(apiClient *Client, cf *Schedules, since, until, tableStyle string) {
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

	parsedSchedule := make(map[string]map[string]uint)
	for _, entry := range schedule.Schedule.FinalSchedule.RenderedScheduleEntries {
		if parsedSchedule[entry.User.Name] == nil {
			parsedSchedule[entry.User.Name] = map[string]uint{}
		}

		parsedSchedule[entry.User.Name]["oncall"]++
		day, _ := weekday(entry.Start)
		if day == 0 || day == 6 {
			parsedSchedule[entry.User.Name]["weekends"]++
		}

		holiday, _ := holidays(entry.Start)
		if len(holiday) > 0 {
			parsedSchedule[entry.User.Name]["holidays"]++
		}
	}

	data := make([]table.Row, 0)
	for user, days := range parsedSchedule {
		data = append(data, table.Row{user, days["weekends"], days["holidays"], days["oncall"]})
	}

	fields := table.Row{"ENGINEER", "WEEKEND", "HOLIDAY", "TOTAL"}
	title := fmt.Sprintf("%s (%s - %s)", Name, since, until)
	printTable(data, fields, title, tableStyle, "TOTAL", table.DscNumeric)
}
