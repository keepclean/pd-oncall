package main

import (
	"log"
	"os"
	"time"
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
