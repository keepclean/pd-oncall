package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jedib0t/go-pretty/text"
)

type CacheFile string

func (c CacheFile) String() string {
	return string(c)
}

func (c CacheFile) Create(apiClient *Client) error {
	if c.Exist() {
		return nil
	}

	log.Printf("Cache file %s doesn't exist; Creating it...\n", c)

	if err := os.MkdirAll(c.DirName(), 0755); err != nil {
		log.Fatalln("can't create directory for cache file: ", err)
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

func (c CacheFile) Exist() bool {
	fInfo, err := os.Stat(c.ExpandPath())
	if err != nil && os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Println("non-IsNotExist error upon calling os.Stat:", err)
		return false
	}

	// if modification of file more than four weeks, refresh it
	if time.Since(fInfo.ModTime()) > (time.Hour * 672) {
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

func (c CacheFile) Write(t []*Schedule) error {
	f, err := os.Create(c.ExpandPath())
	if err != nil {
		log.Println("can't create file: ", err)
		return err
	}

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

func (c CacheFile) Remove() error {
	cf := c.ExpandPath()
	if err := os.Remove(cf); err != nil {
		log.Println("can not remove cache file:", err)
		return err
	}

	log.Println("Cache file", cf, "has been removed")
	return nil
}

func (c CacheFile) Show() {
	cache, err := c.Read()
	if err != nil {
		log.Fatalln("can not read the cache file", c.ExpandPath(), err)
	}

	jsonPrettyPrinter := text.NewJSONTransformer("", "  ")
	fmt.Println(jsonPrettyPrinter(cache))
}
