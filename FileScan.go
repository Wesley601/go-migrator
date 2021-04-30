package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func fileScan(dumppath string, projectpath string) {
	file, err := os.Open(dumppath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)

	scanner.Buffer(buf, 1024*1024)

	var realTable RealTable

	newTable := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "--") == 0 || strings.Index(line, "/*") == 0 {
			continue
		}

		if strings.Contains(line, "CREATE TABLE") {
			la, ok := GetStringInBetween(line, "`", "`")
			if !ok {
				continue
			}
			// fmt.Println(la)
			if strings.Contains(la, "_") {
				var c string
				for _, v := range strings.Split(la, "_") {
					c += strings.Title(v)
				}
				realTable.ClassName = c
			} else {
				realTable.ClassName = strings.Title(la)
			}
			realTable.TableName = la
			newTable = true
			continue
		}

		if newTable {
			if strings.Contains(line, ";") {
				MakeTemplate(realTable, projectpath)
				newTable = false
				realTable = RealTable{}
				continue
			}

			realTable.TableFields = append(realTable.TableFields, ParseToField(line))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
