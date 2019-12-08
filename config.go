package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/jedib0t/go-pretty/text"
)

type ConfigFile string

func (c ConfigFile) String() string {
	return fmt.Sprint(string(c))
}

func (c ConfigFile) Create(t []*PDTeam) {
	if c.Exist() {
		return
	}

	log.Printf("Config file %s doesn't exist;\n", c)
	if err := os.MkdirAll(c.DirName(), 0755); err != nil {
		log.Fatalln("can't create directory for cache file: ", err)
	}

	printTeamsAsTable(t)

	teamsNumbers, err := getUserInput()
	if err != nil {
		log.Fatalln("[ConfigFile.Create] fail to read user input:", err)
	}

	c.Write(t, teamsNumbers)

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

func (c ConfigFile) Write(t []*PDTeam, teamsNumbers []int) {
	len_t := len(t)
	var tSubset []*PDTeam

	for _, n := range teamsNumbers {
		if n > len_t || n < 1 {
			log.Println("There is no a team with number in the list:", n)
			continue
		}

		tSubset = append(tSubset, t[n-1])
	}

	f, err := os.Create(c.ExpandPath())
	if err != nil {
		log.Println("can't create file: ", err)
	}

	var cf PDTeams = PDTeams{Teams: tSubset}
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

func (c ConfigFile) Show() {
	f, err := os.Open(c.ExpandPath())
	if err != nil {
		log.Fatalln("can not open the config file", c.ExpandPath(), err)
	}
	defer f.Close()

	var cf PDTeams
	if err = json.NewDecoder(f).Decode(&cf); err != nil {
		log.Fatalln("can not decode the config file", c.ExpandPath(), err)
	}

	jsonPrettyPrinter := text.NewJSONTransformer("", "  ")
	fmt.Println(jsonPrettyPrinter(cf))
}

func printTeamsAsTable(t []*PDTeam) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer w.Flush()

	rows := 15
	len_t := len(t)
	columns := len_t / rows
	remainder := len_t % rows

	for r := 0; r < rows; r++ {
		var s string
		for c := 0; c < columns; c++ {
			index := c*rows + r
			s = s + fmt.Sprintf("%d) %s\t", index+1, t[index].Name)
		}

		if remainder > 0 {
			s = s + fmt.Sprintf("%d) %s\t", len_t-remainder+1, t[-remainder].Name)
			remainder--
		}

		fmt.Fprintln(w, s)
	}
}

func getUserInput() ([]int, error) {
	fmt.Print("Please select command numbers, separate them by commas: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	input = strings.Trim(input, "\n")

	var result []int

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
