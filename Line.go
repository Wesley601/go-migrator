package main

import (
	"strings"
)

type Line interface {
	New(a string) Line
	GetLaravelMigration() string
}

func ParseToField(line string) string {
	line = strings.TrimSpace(line)

	var newLine Line

	if line[0] == '`' {
		newLine = TableColumn{}
	} else {
		newLine = KeyColumn{}
	}

	ln := newLine.New(line)
	return ln.GetLaravelMigration()
}
