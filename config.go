package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	printTeamsAsTable(t)

	fmt.Print("Please choose team names: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	fmt.Println(input)
}

func (c ConfigFile) Exist() bool {
	if _, err := os.Stat(c.String()); os.IsNotExist(err) {
		return false
	}

	return true
}

func printTeamsAsTable(t []*PDTeam) {
	for _, team := range t {
		fmt.Println(team.Name)
	}
	fmt.Println(len(t))
}
