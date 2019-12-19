package main

import "log"

func oncallUser(apiClient *Client, cf *Schedules, since, until, tableStyle string) {
	printUserAsTable(cf.Users)
	userIDs, err := getUserInput("Please select an user id: ")
	if err != nil {
		log.Println("Failed to get user ID:", err)
		return
	}

	id := userIDs[0]
	if id < 0 || id > len(cf.Users) {
		log.Println("Specified shift ID is out of range:", id)
		return
	}
}
