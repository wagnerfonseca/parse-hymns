package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	s "strings"
)

const sourcePath = "./raw"

var hymns map[int]*Hymn

func main() {

	d, _ := os.Open(sourcePath)
	files, _ := d.Readdir(-1)

	for _, fi := range files {
		fmt.Println(fi.Name())

		filePath := sourcePath + "/" + fi.Name()

		lines := Readlines(filePath)

		hymns = make(map[int]*Hymn)
		for _, value := range lines {
			re := regexp.MustCompile(`(?P<Number>[0-9]+\.) (?P<Title>[a-zA-Zà-úÀ-Ú0-9 \,\']+)`)

			if re.MatchString(value) {
				m := re.FindStringSubmatch(value)

				h := new(Hymn)
				h.Number, _ = strconv.Atoi(s.Replace(m[1], ".", "", -1))
				h.Title = s.Trim(m[2], " ")

				hymns[h.Number] = h
			}
		}

	}
	d.Close()

}

// Readlines read lines file
func Readlines(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return lines
}

type Hymn struct {
	Number int
	Title  string
}
