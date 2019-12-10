package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("pd-oncall", "A command-line tool for represeting PagerDuty oncall schedule.")
	apiToken := app.Flag("api-token", "Auth API token; Might be an environment variable PAGERDUTY_API_TOKEN;").Envar("PAGERDUTY_API_TOKEN").Required().String()
	apiURL := app.Flag("api-url", "API URL.").Default("https://api.pagerduty.com/").URL()

	version := "0.0.0"
	app.Version(version)

	config := app.Command("config", "sub command for managing a config file")
	configRm := config.Flag("rm", "remove config file").Bool()
	config.Flag("show", "show config file").Bool()

	now := app.Command("now", "list currently oncall")

	schedules := app.Command("schedules", "list currently schedules")

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	apiClient := NewPDApiClient(*apiURL, version, *apiToken)

	var cf ConfigFile = "${HOME}/.config/pd-oncall/config.json"
	if !cf.Exist() {
		log.Printf("Config file %s doesn't exist;\n", cf)

		var teamsCache CacheFile = "${HOME}/.cache/pd-oncall/teams-cache.json"
		teamsCache.Create(apiClient)
		pdTeams, err := teamsCache.Read()
		if err != nil {
			log.Fatalln("fail to read teams cache file:", err)
		}
		cf.Create(pdTeams)
	}
	teams := cf.Read()

	switch cmd {
	case config.FullCommand():
		if *configRm {
			cf.Remove()
			return
		}
		cf.Show()
	case now.FullCommand():
		oncallNow(apiClient, teams)
	case schedules.FullCommand():
		scs, err := apiClient.Schedules()
		if err != nil {
			log.Fatalln("fail to get schedules from PD:", err)
		}
		for _, s := range scs {
			fmt.Println(s.ID, s.Name)
		}
	}
}
