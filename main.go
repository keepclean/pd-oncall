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

	now := app.Command("now", "list currently oncall")

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	apiClient := NewPDApiClient(*apiURL, version, *apiToken)

	var cf ConfigFile = "${HOME}/.config/pd-oncall/config.json"
	if !cf.Exist() {
		log.Printf("Config file %s doesn't exist;\n", cf)

		var schedulesCache CacheFile = "${HOME}/.cache/pd-oncall/schedules-cache.json"
		schedulesCache.Create(apiClient)
		pdSchedules, err := schedulesCache.Read()
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
	case now.FullCommand():
		oncallNow(apiClient, schedules, *tableStyle)
	}
}
