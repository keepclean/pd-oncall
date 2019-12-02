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
	var teamsCache CacheFile = "/tmp/pd-teams-cache.json"

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case teams.FullCommand():
		apiClient := NewPDApiClient(*apiURL, version, *apiToken)
		teamsCache.Create(apiClient)
	}

	pdTeams, err := teamsCache.Read()
	if err != nil {
		log.Fatalf("fail to read teams cache file: %v", err)
	}

	for _, team := range pdTeams {
		fmt.Println(team.Name)
	}
}
