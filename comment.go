package pbast

import "strings"

type Comment []string

func (c Comment) String() string {
	var a []string
	c = a
	return strings.Join(a, "\n")
}
