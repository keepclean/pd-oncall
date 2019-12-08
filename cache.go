package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type CacheFile string

func (c CacheFile) String() string {
	return fmt.Sprint(string(c))
}

func (c CacheFile) Create(apiClient *Client) {
	if c.Exist() {
		return
	}

	log.Printf("Cache file %s doesn't exist; Creating it...\n", c)

	if err := os.MkdirAll(c.DirName(), 0755); err != nil {
		log.Fatalln("can't create directory for cache file: ", err)
	}

	pdTeams, err := apiClient.Teams()
	if err != nil {
		log.Println("failt to query PD API: ", err)
	}

	c.Write(pdTeams)

}

func (c CacheFile) Exist() bool {
	fInfo, err := os.Stat(c.ExpandPath())
	if os.IsNotExist(err) {
		return false
	}
	// if modification of file more than a week, refresh it
	if time.Since(fInfo.ModTime()) > (time.Hour * 168) {
		return false
	}

	return true
}

func (c CacheFile) DirName() string {
	return filepath.Dir(c.ExpandPath())
}

func (c CacheFile) ExpandPath() string {
	return os.ExpandEnv(c.String())
}

func (c CacheFile) Write(t []*PDTeam) {
	f, err := os.Create(c.ExpandPath())
	if err != nil {
		log.Println("can't create file: ", err)
	}

	if err = json.NewEncoder(f).Encode(t); err != nil {
		log.Println("can't write json: ", err)
	}
}

func (c CacheFile) Read() ([]*PDTeam, error) {
	f, err := os.Open(c.ExpandPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var t []*PDTeam
	if err = json.NewDecoder(f).Decode(&t); err != nil {
		return nil, err
	}

	return t, nil
}
