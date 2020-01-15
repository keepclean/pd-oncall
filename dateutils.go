package main

import (
	"log"
	"time"

	"github.com/rickar/cal"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Dates struct {
	Since, Until string
}

func NewDates(c *kingpin.CmdClause) *Dates {
	t := time.Now()
	d := &Dates{}

	c.Flag("since", "The start of the date range over which you want to search.").Default(t.Format("2006-01-02")).StringVar(&d.Since)
	c.Flag("until", "The end of the date range over which you want to search.").Default(t.AddDate(0, 0, 7).Format("2006-01-02")).StringVar(&d.Until)

	return d
}

func (d *Dates) CheckDates() {
	t1, err := time.Parse("2006-01-02", d.Since)
	if err != nil {
		log.Fatalf("fail to parse date %s: %v\n", d.Since, err)
	}

	t2, err := time.Parse("2006-01-02", d.Until)
	if err != nil {
		log.Fatalf("fail to parse date %s: %v\n", d.Since, err)
	}

	if t2.Sub(t1) < 0 {
		t2 = t1.AddDate(0, 0, 7)
		log.Println("\"until\" time is in the past relatively to \"since\"; Using \"until\" as: \"since\" plus a week.")
	}

	d.Since = t1.Format("2006-01-02")
	d.Until = t2.Format("2006-01-02")
}

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
