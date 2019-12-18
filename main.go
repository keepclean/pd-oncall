package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("pd-oncall", "A command-line tool for represeting PagerDuty oncall schedule.")
	apiToken := app.Flag("api-token", "Auth API token; Might be an environment variable PAGERDUTY_API_TOKEN;").Envar("PAGERDUTY_API_TOKEN").Required().String()
	apiURL := app.Flag("api-url", "API URL.").Default("https://api.pagerduty.com/").URL()
	tableStyle := app.Flag("table-style", "Table style: rounded, box, colored").Default("rounded").String()

	version := "0.0.0"
	app.Version(version)

	config := app.Command("config", "sub command for managing a config file")
	configRm := config.Flag("rm", "remove config file").Bool()
	config.Flag("show", "show config file").Bool()
	// configFile := config.Flag("file", "Path to config file").Default("${HOME}/.config/pd-oncall/config.json").String()

	cache := app.Command("cache", "sub command for managing a cache file")
	cacheRm := cache.Flag("rm", "remove cache file").Bool()
	cache.Flag("show", "show cache file").Bool()

	now := app.Command("now", "list currently oncall for schedules in a config file")

	schedule := app.Command("schedule", "Oncall schedule information")
	scheduleDates := NewDates(schedule)

	report := app.Command("report", "generates report")
	reportDates := NewDates(report)

	roster := app.Command("roster", "roster for all known schedules")
	rosterDates := NewDates(roster)

	user := app.Command("user", "oncall schedule for user")
	userDates := NewDates(user)

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	apiClient := NewPDApiClient(*apiURL, version, *apiToken)

	var cf ConfigFile = "${HOME}/.config/pd-oncall/config.json"
	var sc CacheFile = "${HOME}/.cache/pd-oncall/schedules-cache.json"
	if !cf.Exist() {
		log.Printf("Config file %s doesn't exist;\n", cf)

		sc.Create(apiClient)
		pdSchedules, err := sc.Read()
		if err != nil {
			log.Fatalln("fail to read schedules cache file:", err)
		}
		cf.Create(pdSchedules)
	}
	schedules := cf.Read()

	switch cmd {
	case config.FullCommand():
		if *configRm {
			cf.Remove()
			return
		}
		cf.Show()
	case cache.FullCommand():
		if *cacheRm {
			sc.Remove()
			return
		}
		sc.Show()
	case now.FullCommand():
		oncallNow(apiClient, schedules, *tableStyle)
	case schedule.FullCommand():
		scheduleDates.CheckDates()
		oncallShift(apiClient, schedules, scheduleDates.Since, scheduleDates.Until, *tableStyle)
	case report.FullCommand():
		reportDates.CheckDates()
		oncallReport(apiClient, schedules, reportDates.Since, reportDates.Until, *tableStyle)
	case roster.FullCommand():
		rosterDates.CheckDates()
		oncallRoster(apiClient, schedules, rosterDates.Since, rosterDates.Until, *tableStyle)
	case user.FullCommand():
		userDates.CheckDates()
		oncallRoster(apiClient, schedules, rosterDates.Since, rosterDates.Until, *tableStyle)
	}
}
