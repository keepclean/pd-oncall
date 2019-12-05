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

	teams := app.Command("teams", "show and store in file cache teams which are defined in PD")

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	apiClient := NewPDApiClient(*apiURL, version, *apiToken)

	// Create cache file for PD teams
	var teamsCache CacheFile = "${HOME}/.cache/pd-oncall/teams-cache.json"
	teamsCache.Create(apiClient)
	pdTeams, err := teamsCache.Read()
	if err != nil {
		log.Fatalf("fail to read teams cache file: %v", err)
	}

	// Create config file
	var cf ConfigFile = "${HOME}/.config/pd-oncall/config.json"
	cf.Create(pdTeams)

	switch cmd {
	case teams.FullCommand():
		fmt.Println(cmd)
	}
}