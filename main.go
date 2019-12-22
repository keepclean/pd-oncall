package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("pd-oncall", "A command-line tool for represeting PagerDuty oncall schedules")
	apiToken := app.Flag("api-token", "Auth API token; Might be an environment variable PAGERDUTY_API_TOKEN").Envar("PAGERDUTY_API_TOKEN").Required().String()
	apiURL := app.Flag("api-url", "Pager Duty API URL").Default("https://api.pagerduty.com/").URL()
	tableStyle := app.Flag("table-style", "Available table styles: rounded, box, colored").Default("rounded").String()

	version := "0.0.0"
	app.Version(version)

	config := app.Command("config", "simple management of a config file")
	configRm := config.Flag("rm", "remove a config file").Bool()
	config.Flag("show", "show a config file").Bool()
	// configFile := config.Flag("file", "Path to config file").Default("${HOME}/.config/pd-oncall/config.json").String()

	cache := app.Command("cache", "simple management of a cache file")
	cacheRm := cache.Flag("rm", "remove a cache file").Bool()
	cache.Flag("show", "show a cache file").Bool()

	now := app.Command("now", "list currently oncall for schedules in a config file")

	schedule := app.Command("schedule", "oncall schedule information")
	scheduleDates := NewDates(schedule)

	report := app.Command("report", "generates a simple oncall report")
	reportDates := NewDates(report)

	roster := app.Command("roster", "roster for all known schedules")
	rosterDates := NewDates(roster)

	user := app.Command("user", "oncall schedule for a specific user")
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
	cfg := cf.Read()

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
		oncallNow(apiClient, cfg, *tableStyle)
	case schedule.FullCommand():
		scheduleDates.CheckDates()
		oncallShift(apiClient, cfg, scheduleDates.Since, scheduleDates.Until, *tableStyle)
	case report.FullCommand():
		reportDates.CheckDates()
		oncallReport(apiClient, cfg, reportDates.Since, reportDates.Until, *tableStyle)
	case roster.FullCommand():
		rosterDates.CheckDates()
		oncallRoster(apiClient, cfg, rosterDates.Since, rosterDates.Until, *tableStyle)
	case user.FullCommand():
		userDates.CheckDates()
		oncallUser(apiClient, cfg, userDates.Since, userDates.Until, *tableStyle)
	}
}
