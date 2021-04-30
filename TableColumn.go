package main

import (
	"fmt"
	"strings"
)

type TableColumn struct {
	Name          string
	Tp            string
	TpRange       string
	DefaultValue  string
	IsNullable    bool
	IsUnsigned    bool
	splittedValue []string
	rawValue      string
}

func (t TableColumn) New(value string) Line {
	value = strings.TrimRight(value, ",")
	t.splittedValue = strings.Split(value, " ")
	t.rawValue = value

	tp, rg := t.GetType()

	t.Name = t.GetName()
	t.Tp = tp
	t.TpRange = rg
	t.DefaultValue = t.getDefaultValue()
	t.IsNullable = t.isNullable()
	t.IsUnsigned = t.isUnsigned()

	return t
}

func (l TableColumn) GetName() string {
	if l.Name != "" {
		return l.Name
	}
	return strings.ReplaceAll(l.splittedValue[0], "`", "")
}

func (l TableColumn) GetType() (string, string) {

	if l.Tp != "" || l.TpRange != "" {
		return l.Tp, l.TpRange
	}

	a, _ := GetStringInBetween(l.splittedValue[1], "(", ")")
	b := strings.Split(l.splittedValue[1], "(")[0]
	return b, a
}

func (l *TableColumn) isUnsigned() bool {
	return l.splittedValue[2] == "unsigned"
}

func (l *TableColumn) isNullable() bool {
	nullable := true
	for i, value := range l.splittedValue {
		if value == "NOT" && l.splittedValue[i+1] == "NULL" {
			nullable = false
			break
		}
	}
	return nullable
}

func (l *TableColumn) getDefaultValue() string {
	var defaultValue string
	for i, value := range l.splittedValue {
		if value == "DEFAULT" {
			defaultValue = l.splittedValue[i+1]
			break
		}
	}
	return defaultValue
}

func (t TableColumn) GetLaravelMigration() string {
	m := fmt.Sprintf("$table->%s('%s')", ParseTypesToLaravel(t.Tp), t.Name)

	if t.TpRange != "" {
		m = strings.Replace(m, ")", ", "+t.TpRange+")", 1)
	}

	if t.IsNullable {
		m += "->nullable()"
	}

	if t.IsUnsigned {
		m += "->unsigned()"
	}

	if t.DefaultValue != "" && !t.IsNullable {
		if t.DefaultValue == "CURRENT_TIMESTAMP" {
			m += "->useCurrent()"
		} else {
			m += "->default(" + t.DefaultValue + ")"
		}
	}

	return m + ";"
}
