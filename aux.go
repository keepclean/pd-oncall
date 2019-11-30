package main

import (
	"os"
	"time"
)

func Exists(f string) bool {
	fInfo, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	// if modification of file more than a week, refresh it
	if time.Since(fInfo.ModTime()) > (time.Hour * 168) {
		return false
	}

	return true
}
