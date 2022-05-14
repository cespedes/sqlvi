package main

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

func writeOrgtable(w io.Writer, columns []string, data [][]string) {
	widths := make([]int, len(columns))
	for i, x := range columns {
		widths[i] = utf8.RuneCountInString(x)
	}
	for _, x := range data {
		for i, y := range x {
			if utf8.RuneCountInString(y) > widths[i] {
				widths[i] = utf8.RuneCountInString(y)
			}
		}
	}
	line := fmt.Sprint("|", strings.Repeat("-", widths[0]+2))
	for i := range columns[1:] {
		line += "+" + strings.Repeat("-", widths[i+1]+2)
	}
	line += "|"
	fmt.Fprint(w, line, "\n|")
	for i, x := range columns {
		fmt.Fprintf(w, " %-*s |", widths[i], x)
	}
	fmt.Fprint(w, "\n", line, "\n")
	for _, x := range data {
		fmt.Fprintf(w, "|")
		for i, y := range x {
			fmt.Fprintf(w, " %-*s |", widths[i], y)
		}
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintln(w, line)
}

func readOrgtable(r io.Reader) (columns []string, data [][]string) {
	return columns, data
}
