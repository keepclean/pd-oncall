package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

type ConfigFile struct {
	FileAsset
}

func (c ConfigFile) Create(apiClient *Client, cacheFile CacheFile) {
	if !cacheFile.Exist() || cacheFile.Stale() {
		if err := cacheFile.Create(apiClient); err != nil {
			errRm := cacheFile.Remove()
			log.Fatalln("can't create cache file:", err, errRm)
		}
	}

	pdSchedules, err := cacheFile.Read()
	if err != nil {
		log.Fatalln("fail to read schedules cache file:", err)
	}

	lenS := len(pdSchedules.Schedules)
	if lenS == 0 {
		errRm := cacheFile.Remove()
		log.Fatalln(
			"cache file for schedules is empty, can't create config file;",
			"cache file has been removed", errRm)
	}

	if err := c.CreateDirs(); err != nil {
		log.Fatalln(err)
	}

	printSchedulesAsTable(pdSchedules.Schedules)

	scheduleNumbers, err := getUserInput("Please select numbers, separate them by commas: ")
	if err != nil {
		log.Fatalln("[ConfigFile.Create] fail to read user input:", err)
	}

	pdSubset := make([]*Schedule, 0, len(scheduleNumbers))
	users := make(map[string]string)

	for _, n := range scheduleNumbers {
		if n > lenS || n < 1 {
			log.Println("There is no schedule with number in the list:", n)
			continue
		}

		item := pdSchedules.Schedules[n-1]
		pdSubset = append(pdSubset, &Schedule{item.ID, item.Name, item.Description, []*User{}})
		for _, user := range item.Users {
			if user.Deleted != "" {
				continue
			}

			if _, ok := users[user.ID]; !ok {
				users[user.ID] = user.Name
			}
		}
	}

	cf := Schedules{Schedules: pdSubset, Users: users}
	if err := c.Write(&cf); err != nil {
		log.Fatalln(err)
	}
}

func printSchedulesAsTable(t []*Schedule) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	lenT := len(t)
	rows := rowsNumber(lenT)
	columns := lenT / rows
	remainder := lenT % rows
	startFrom := 1

	for r := 0; r < rows; r++ {
		var s string
		for c := 0; c < columns; c++ {
			index := c*rows + r
			s += fmt.Sprintf("%d) %s\t", index+startFrom, t[index].Name)
		}

		if remainder > 0 {
			s += fmt.Sprintf("%d) %s\t", lenT-remainder+startFrom, t[lenT-remainder].Name)
			remainder--
		}

		fmt.Fprintln(w, s)
	}
}

func rowsNumber(i int) int {
	switch {
	case i > 150:
		return 50
	case i > 120:
		return 40
	case i > 90:
		return 30
	case i > 60:
		return 20
	case i < 15:
		return i
	default:
		return 15
	}
}

func getUserInput(prompt string) ([]int, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	input = strings.Trim(input, "\n")

	result := make([]int, 0)

	for _, s := range strings.Split(input, ",") {
		i, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Printf("[getUserInput/strconv.Atoi] error while converting a part of user input: %s - %v\n", s, err)
			continue
		}
		result = append(result, i)
	}

	return result, nil
}

func printUserAsTable(users map[string]string) []string {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	userKeys := make([]string, 0, len(users))
	for id := range users {
		userKeys = append(userKeys, id)
	}

	sort.Slice(userKeys, func(i, j int) bool { return userKeys[i] < userKeys[j] })

	lenUsers := len(users)
	rows := rowsNumber(lenUsers)
	columns := lenUsers / rows
	remainder := lenUsers % rows
	startFrom := 1

	for r := 0; r < rows; r++ {
		var s string
		for c := 0; c < columns; c++ {
			index := c*rows + r
			s += fmt.Sprintf("%d) %s\t", index+startFrom, users[userKeys[index]])
		}

		if remainder > 0 {
			s += fmt.Sprintf("%d) %s\t", lenUsers-remainder+startFrom, users[userKeys[lenUsers-remainder]])
			remainder--
		}

		fmt.Fprintln(w, s)
	}

	return userKeys
}
