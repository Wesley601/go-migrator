package main

import (
	"fmt"
	"strings"
)

type KeyColumn struct {
	KeyType       string
	originalValue string
	Fields        string
	Reference     [2]string
	splittedValue []string
	rawValue      string
}

func (k KeyColumn) New(value string) Line {
	k.originalValue = value
	k.splittedValue = strings.Split(value, " ")
	k.rawValue = value
	k.KeyType = k.getKeyType()
	k.Fields = k.getFields()
	k.Reference = k.getReference()
	return k
}

func (l KeyColumn) GetName() string {
	return ""
}

func (l KeyColumn) GetType() (string, string) {
	return "KEY", ""
}

func (k *KeyColumn) getReference() [2]string {
	reference := [2]string{"", ""}
	if k.KeyType != "CONSTRAINT" {
		return reference
	}

	reference[0], _ = GetStringInBetween(k.splittedValue[6], "`", "`")
	reference[1], _ = GetStringInBetween(k.splittedValue[7], "`", "`")

	return reference
}

func (k *KeyColumn) getFields() string {
	RawField, _ := GetStringInBetween(k.originalValue, "(", ")")
	return strings.ReplaceAll(RawField, "`", "'")
}

func (k *KeyColumn) getKeyType() string {
	return k.splittedValue[0]
}

func (k KeyColumn) GetLaravelMigration() string {
	if k.splittedValue[0] == "KEY" {
		return ""
	}
	m := "$table->"
	var fn string

	switch keyType := k.KeyType; keyType {
	case "PRIMARY":
		if strings.Contains(k.Fields, ",") {
			fn = fmt.Sprintf("primary([%s])", k.Fields)
		} else {
			fn = fmt.Sprintf("primary(%s)", k.Fields)
		}
	case "UNIQUE":
		fn = fmt.Sprintf("unique(%s)", k.Fields)
	case "CONSTRAINT":
		fn = fmt.Sprintf(
			"foreign(%s)->references('%s')->on('%s')",
			k.Fields,
			k.Reference[0],
			k.Reference[1],
		)
	default:
	}

	return m + fn + ";"
}
