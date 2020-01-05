package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/text"
)

type CacheFile struct {
	FileAsset
}

func (c CacheFile) Create(apiClient *Client) error {
	if err := c.CreateDirs(); err != nil {
		return err
	}

	schedules, err := apiClient.Schedules()
	if err != nil {
		log.Println("fail to query API: ", err)
		return err
	}

	if err := c.Write(schedules); err != nil {
		return err
	}

	return nil
}

func (c CacheFile) Stale() bool {
	fInfo, err := os.Stat(c.ExpandPath())
	if err != nil {
		log.Printf("couldn't check stat of %q file: %v", c, err)
		return false
	}

	// if cache file was modified more than four weeks ago, refresh that file
	if time.Since(fInfo.ModTime()) > (time.Hour * 24 * 7 * 4) {
		return false
	}

	return true
}

func (c CacheFile) Write(t []*Schedule) error {
	f, err := os.Create(c.ExpandPath())
	if err != nil {
		log.Println("can't create file: ", err)
		return err
	}
	defer f.Close()

	if err = json.NewEncoder(f).Encode(t); err != nil {
		log.Println("can't write json: ", err)
		return err
	}

	return nil
}

func (c CacheFile) Read() ([]*Schedule, error) {
	f, err := os.Open(c.ExpandPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var t []*Schedule
	if err = json.NewDecoder(f).Decode(&t); err != nil {
		return nil, err
	}

	return t, nil
}

func (c CacheFile) Show() {
	cache, err := c.Read()
	if err != nil {
		log.Fatalln("can not read the cache file", c.ExpandPath(), err)
	}

	jsonPrettyPrinter := text.NewJSONTransformer("", "  ")
	fmt.Println(jsonPrettyPrinter(cache))
}
