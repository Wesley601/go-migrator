package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

func MakeTemplate(table RealTable, projectpath string) {
	table.TableFields = handleTT(table.TableFields)
	tml, err := template.New("migration_template").Parse(migration_template)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(
		fmt.Sprintf(projectpath+"/database/migrations/%s_create_%s_table.php", FormatDateFileName(), table.TableName),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = tml.Execute(f, table)
	if err != nil {
		log.Fatal(err)
	}
}

func handleTT(fl []string) []string {
	var f2 []string
	a, ok := findCreatedAt(fl)
	if ok {
		f2 = append(fl[:a], fl[a+1:]...)
	}

	b, ok1 := findUpdatedAt(fl)
	if ok1 {
		f2 = append(f2[:b], f2[b+1:]...)
	}

	if ok && ok1 {
		f2 = append(f2, "$table->timestamps();")
		return f2
	}

	return fl
}

func findCreatedAt(table []string) (int, bool) {
	ok := false
	var index int
	for i, field := range table {
		if strings.Contains(field, "created_at") {
			index = i
			ok = true
		}
	}
	return index, ok
}

func findUpdatedAt(table []string) (int, bool) {
	ok := false
	var index int
	for i, field := range table {
		if strings.Contains(field, "updated_at") {
			index = i
			ok = true
		}
	}
	return index, ok
}
