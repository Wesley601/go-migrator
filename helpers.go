package main

import (
	"fmt"
	"strings"
	"time"
)

const migration_template = `<?php
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class Create{{ .ClassName }}Table extends Migration
{
	/**
	* Run the migrations.
	*
	* @return void
	*/
	public function up()
	{
		Schema::table('{{ .TableName }}', function(Blueprint $table) {
		{{- range .TableFields }}
			{{ . }}
		{{- end }}
		});
	}

	/**
	* Reverse the migrations.
	*
	* @return void
	*/
	public function down()
	{
		Schema::DropColumn('{{ .TableName }}');
	}
}`

func GetStringInBetween(str string, start string, end string) (result string, ok bool) {
	s := strings.Index(str, start)
	if s == -1 {
		return "", false
	}

	s += 1

	b := strings.Index(str[s:], end)
	if b == -1 {
		return "", false
	}

	return str[s : s+b], true
}

func ParseTypesToLaravel(field string) string {
	types := map[string]string{
		"int":       "integer",
		"tinyint":   "tinyInteger",
		"smallint":  "smallInteger",
		"mediumint": "mediumInteger",
		"bigint":    "bigInteger",
		"decimal":   "decimal",
		"float":     "float",
		"boolean":   "boolean",
		"varchar":   "string",
		"timestamp": "timestamp",
		"text":      "text",
		"longtext":  "longText",
		"enum":      "enum",
		"point":     "point",
	}

	tp, ok := types[field]

	if !ok {
		fmt.Println("Tipo desconhecido ", field)
	}

	return tp
}

func ParseDefaultToLaravel(field string) string {

	switch field {
	case "NULL":
		return "null"
	default:
		return field
	}
}

func FormatDateFileName() string {
	t := time.Now().Add(time.Hour * 3).Format(time.RFC3339Nano)
	b := strings.Split(t, "T")

	a := strings.ReplaceAll(b[0], "-", "_")
	x := strings.ReplaceAll(b[1], ".", "_")
	v := strings.ReplaceAll(x, ":", "")
	a += "_" + v

	return a
}
