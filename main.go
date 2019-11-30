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
	apiURL := app.Flag("api-host", "API host.").Default("https://api.pagerduty.com/").URL()
	version := "0.0.0"

	teams := app.Command("teams", "show and store in file cache teams which are defined in PD")

	app.Version(version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case teams.FullCommand():
		c := NewPDApiClient(*apiURL, version, *apiToken)
		teams, err := c.Teams()
		if err != nil {
			log.Fatalf("fail to query teams from %s: %v", *apiURL, err)
		}
		for _, team := range teams {
			fmt.Println(team.Name)
		}
	}
}
