package gql

import (
	"strconv"
	"strings"
)

type GraphqlError struct {
	Message   string `json:"message"`
	Locations []struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	} `json:"locations"`
}

func (e GraphqlError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	b.WriteString(" at [")
	for i := 0; i < len(e.Locations); i++ {
		b.WriteString(strconv.Itoa(e.Locations[i].Line))
		b.WriteString(":")
		b.WriteString(strconv.Itoa(e.Locations[i].Column))
		if i < len(e.Locations)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")
	return b.String()
}

type GraphqlErrors []GraphqlError

func (e GraphqlErrors) Error() string {
	var b strings.Builder
	for i := 0; i < len(e); i++ {
		b.WriteString(e[i].Error())
		if i < len(e)-1 {
			b.WriteString("; ")
		}
	}
	return b.String()
}
