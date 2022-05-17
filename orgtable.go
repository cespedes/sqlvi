package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

func writeOrgTable(w io.Writer, columns []string, data [][]string) {
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

func readOrgLine(line string) []string {
	s := strings.Split(line, "|")
	if len(s) < 3 {
		return nil
	}
	s = s[1 : len(s)-1]
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	return s
}

func readOrgTable(r io.Reader) (columns []string, data [][]string, err error) {
	lineNo := 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		lineNo++
		if strings.Contains(s.Text(), `|---`) {
			break
		}
	}
	if !s.Scan() {
		return nil, nil, fmt.Errorf("No table found after reading %d lines of text.", lineNo)
	}
	lineNo++
	columns = readOrgLine(s.Text())
	if len(columns) == 0 {
		return nil, nil, fmt.Errorf("Wrong header for table in line %d.", lineNo)
	}
	if !s.Scan() {
		return nil, nil, fmt.Errorf("No table found after header in line %d.", lineNo)
	}
	lineNo++
	if !strings.Contains(s.Text(), `|---`) {
		return nil, nil, fmt.Errorf("Wrong table found after header in line %d.", lineNo)
	}
	for s.Scan() {
		lineNo++
		line := s.Text()
		if strings.Contains(line, `|---`) {
			break
		}
		s := readOrgLine(s.Text())
		if len(s) != len(columns) {
			return nil, nil, fmt.Errorf("Wrong number of columns in line %d.", lineNo)
		}
		data = append(data, s)
	}
	return columns, data, nil
}
