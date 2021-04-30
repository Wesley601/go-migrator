package main

import (
	"flag"
	"fmt"
	"time"
)

type RealTable struct {
	ClassName   string
	TableName   string
	TableFields []string
}

type Table []string

func main() {
	s := time.Now()

	inline := flag.String("inline", "", "parse a single line of a drump")
	dumppath := flag.String("dumppath", "", "dump path to parse")
	projectpath := flag.String("projectpath", "", "project path")

	flag.Parse()

	if *inline != "" && *dumppath != "" {
		panic("no arguments passed")
	}

	if *inline != "" {
		fmt.Println(*inline)
		fmt.Println(ParseToField(*inline))
	}

	if *dumppath != "" {
		fileScan(*dumppath, *projectpath)
	}

	fmt.Println(time.Since(s))
}
