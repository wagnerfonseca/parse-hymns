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

var (
	hymns map[int]*Hymn
	re    = regexp.MustCompile(`(?P<Number>[0-9]+\.) (?P<Title>[a-zA-Zà-úÀ-Ú0-9 \,\']+)`)
)

func main() {

	d, _ := os.Open(sourcePath)
	files, _ := d.Readdir(-1)

	for _, fi := range files {
		fmt.Println(fi.Name())

		filePath := sourcePath + "/" + fi.Name()

		lines := Readlines(filePath)

		hymns = make(map[int]*Hymn)
		for _, value := range lines {

			h := new(Hymn)

			isTitle := re.MatchString(value)
			if isTitle {
				m := re.FindStringSubmatch(value)

				h.Number, _ = strconv.Atoi(s.Replace(m[1], ".", "", -1))
				h.Title = s.Trim(m[2], " ")

			} else {
				rev := regexp.MustCompile(`(?P<Number>[0-9])(?P<Line>[a-zA-Zà-úÀ-Ú0-9 \,\'\!\;\-]+)`)

				fmt.Printf("%q\n", rev.FindString(value))
			}

			hymns[h.Number] = h
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

// Hymn type
type Hymn struct {
	Number int
	Title  string
	Verse  string
	Chorus string
}
