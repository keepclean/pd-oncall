package main

import (
	"log"
	"time"

	"github.com/rickar/cal"
)

func convertTime(t, f string) (string, error) {
	if f == "" {
		f = "2006-01-02 15:04"
	}

	parsedTime, err := parseTime(t)
	if err != nil {
		return "", err
	}

	return parsedTime.Format(f), nil
}

func parseTime(t string) (pt time.Time, err error) {
	pt, err = time.Parse(time.RFC3339, t)
	return
}

func holidays(t string) ([]string, error) {
	var cUK = cal.NewCalendar()
	var cUS = cal.NewCalendar()
	cal.AddBritishHolidays(cUK)
	cal.AddUsHolidays(cUS)

	var countries []string

	parsedTime, err := parseTime(t)
	if err != nil {
		return countries, err
	}

	if cUK.IsHoliday(parsedTime) {
		countries = append(countries, "UK")
	}
	if cUS.IsHoliday(parsedTime) {
		countries = append(countries, "US")
	}

	return countries, nil
}

func weekday(t string) (time.Weekday, error) {
	parsedTime, err := parseTime(t)
	if err != nil {
		return 0, err
	}

	return parsedTime.Weekday(), nil
}

func checkDate(s string) string {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalln(err)
	}
	return t.Format("2006-01-02")
}
