package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	s "strings"
)

const sourcePath = "./raw_test"

var (
	hymns    map[int]Hymn
	metrics  map[int]*position
	re       = regexp.MustCompile(`(?P<Number>[0-9]+\.) (?P<Title>[a-zA-Zà-úÀ-Ú0-9 \,\']+)`)
	reIniDig = regexp.MustCompile(`^[1-9]`)
)

type (
	// Hymn type
	Hymn struct {
		Number int
		Title  string
		Verse  []Verse
		Chorus []string
	}

	// Verse strophes of the hymnal
	Verse struct {
		Number int
		Verse  string
	}

	position struct {
		start int
		end   int
	}
)

func main() {

	d, _ := os.Open(sourcePath)
	files, _ := d.Readdir(-1)

	for _, fi := range files {
		fmt.Println(fi.Name())

		filePath := sourcePath + "/" + fi.Name()

		lines := Readlines(filePath)

		// Delimited the start and end each of hymn
		metrics = delimeted(lines)
		// for k, v := range metrics {
		// 	fmt.Printf("N: %3d - s: %5d e: %5d\n", k, v.start, v.end)
		// }

		//
		var h Hymn
		var numberHymn, numberVerse int
		var isInitVerse, isChorus bool
		var title string
		var chorus []string

		hymns = make(map[int]Hymn)
		for idx, value := range lines {

			isTitle := re.MatchString(value)
			if isTitle {
				isChorus = false
				chorus = nil

				m := re.FindStringSubmatch(value)
				numberHymn, _ = strconv.Atoi(s.Replace(m[1], ".", "", -1))
				title = s.Trim(m[2], " ")
			} else {
				// defines if is verse or chorus
				// TODO: create function return if is verse or chorus
				if len(value) > 1 {
					numberVerse, _ = strconv.Atoi(reIniDig.FindString(value))
					if isInitVerse && numberVerse == 0 {
						isChorus = true
					}
					isInitVerse = false
				} else {
					isInitVerse = true
					isChorus = false
				}

				// Is verse
				if !isInitVerse && !isChorus {

				}
				// is Chorus
				if len(value) > 0 && isChorus {
					// TODO: remover os tab \t
					chorus = append(chorus, s.Trim(value, " "))
				}

			}

			// get index the hymn
			if v, ok := metrics[numberHymn]; ok {
				if v.start == idx { // create
					h = Hymn{Number: numberHymn, Title: title}
				}
				// fmt.Printf("-> %d - %d \n", idx+1, v.end)
				if v.end == idx+1 { // add map
					h.Chorus = chorus
					hymns[numberHymn] = h
				}
			}

		}

		for k, v := range hymns {
			fmt.Printf("id: %d - %v\n", k, v)
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

func delimeted(lines []string) map[int]*position {
	var n, ln, li int // n: number , ln: last number, li: last index
	mts := make(map[int]*position)
	for i, value := range lines {
		if re.MatchString(value) { // title hymn
			m := re.FindStringSubmatch(value)
			n, _ = strconv.Atoi(s.Replace(m[1], ".", "", -1))

			if ln != n {
				// get the last index hymn
				if v, ok := mts[ln]; ok {
					v.start = li
					v.end = i
				}
				// last index (li)
				if li != i {
					li = i
				}
				if _, ok := mts[n]; !ok {
					mts[n] = &position{start: i}
				}
				ln = n
			}
		}

		// last line file
		if len(lines)-1 == i {
			if v, ok := mts[n]; ok {
				v.end = len(lines)
			}
		}
	}
	return mts
}
