package main

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
)

func printTable(data []table.Row, fields table.Row, style string) {
	var tableStyles = map[string]table.Style{
		"box":     table.StyleDefault,
		"rounded": table.StyleRounded,
		"colored": table.StyleColoredBright,
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(tableStyles[style])
	t.AppendHeader(fields)
	t.AppendRows(data)
	t.Render()
}
