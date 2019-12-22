package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/jedib0t/go-pretty/text"
)

type ConfigFile string

func (c ConfigFile) String() string {
	return string(c)
}

func (c ConfigFile) Create(t []*Schedule) {
	if err := os.MkdirAll(c.DirName(), 0755); err != nil {
		log.Fatalln("can't create directory for cache file: ", err)
	}

	printSchedulesAsTable(t)

	scheduleNumbers, err := getUserInput("Please select numbers, separate them by commas: ")
	if err != nil {
		log.Fatalln("[ConfigFile.Create] fail to read user input:", err)
	}

	c.Write(t, scheduleNumbers)
}

func (c ConfigFile) Exist() bool {
	if _, err := os.Stat(c.ExpandPath()); os.IsNotExist(err) {
		return false
	}

	return true
}

func (c ConfigFile) DirName() string {
	return filepath.Dir(c.ExpandPath())
}

func (c ConfigFile) ExpandPath() string {
	return os.ExpandEnv(c.String())
}

func (c ConfigFile) Write(t []*Schedule, scheduleNumbers []int) {
	lenT := len(t)
	tSubset := make([]*Schedule, 0)
	users := make(map[string]string)

	for _, n := range scheduleNumbers {
		if n > lenT || n < 1 {
			log.Println("There is no schedule with number in the list:", n)
			continue
		}

		item := t[n-1]
		tSubset = append(tSubset, &Schedule{item.ID, item.Name, item.Description, []*User{}})
		for _, user := range item.Users {
			if user.Deleted != "" {
				continue
			}

			if _, ok := users[user.ID]; !ok {
				users[user.ID] = user.Name
			}
		}
	}

	f, err := os.Create(c.ExpandPath())
	if err != nil {
		log.Println("can't create file: ", err)
	}

	var cf Schedules = Schedules{Schedules: tSubset, Users: users}
	if err = json.NewEncoder(f).Encode(cf); err != nil {
		log.Println("can't write json: ", err)
	}
}

func (c ConfigFile) Remove() {
	cf := c.ExpandPath()
	if err := os.Remove(cf); err != nil {
		log.Fatalln("can not remove config file", cf, err)
	}

	log.Println("Config file", cf, "has been removed")
}

func (c ConfigFile) Read() *Schedules {
	f, err := os.Open(c.ExpandPath())
	if err != nil {
		log.Fatalln("can not open the config file", c.ExpandPath(), err)
	}
	defer f.Close()

	var cf Schedules
	if err = json.NewDecoder(f).Decode(&cf); err != nil {
		log.Fatalln("can not decode the config file", c.ExpandPath(), err)
	}

	return &cf
}

func (c ConfigFile) Show() {
	cf := c.Read()

	jsonPrettyPrinter := text.NewJSONTransformer("", "  ")
	fmt.Println(jsonPrettyPrinter(cf))
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

	userKeys := []string{}
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
