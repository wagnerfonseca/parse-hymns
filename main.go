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
	hymns    map[int]*Hymn
	indxs    map[int]*indexes
	re       = regexp.MustCompile(`(?P<Number>[0-9]+\.) (?P<Title>[a-zA-Zà-úÀ-Ú0-9 \,\']+)`)
	reIniDig = regexp.MustCompile(`^[1-9]`)
)

type (
	// Hymn type
	Hymn struct {
		Number int
		Title  string
		Verse  []Verse
		Chorus string
	}

	// Verse strophes of the hymnal
	Verse struct {
		Number int
		Verse  string
	}

	indexes struct {
		init int
		end  int
	}
)

func main() {

	d, _ := os.Open(sourcePath)
	files, _ := d.Readdir(-1)

	for _, fi := range files {
		fmt.Println(fi.Name())

		filePath := sourcePath + "/" + fi.Name()

		lines := Readlines(filePath)

		hymns = make(map[int]*Hymn)

		var isInitVerse, isChorus bool
		var n, u, init, numberVerse, cChorus int
		chorus := make([]string, 30)

		// capturar indice de inicio e fim hino

		var nTitle int
		for i, value := range lines {
			isTitle := re.MatchString(value)
			if isTitle {
				m := re.FindStringSubmatch(value)
				nTitle, _ = strconv.Atoi(s.Replace(m[1], ".", "", -1))

				fmt.Printf("i. %d | n %d | u %d \n", i, nTitle, u)
				indxs = make(map[int]*indexes)
				if u != nTitle {
					indxs[nTitle] = &indexes{init: i}
					u = nTitle
				}
			}

		}

		for idx, value := range lines {

			isTitle := re.MatchString(value)
			if isTitle {
				init = idx
				isChorus = false
				cChorus = 0
				chorus = make([]string, 30)

				m := re.FindStringSubmatch(value)
				n, _ = strconv.Atoi(s.Replace(m[1], ".", "", -1))
			} else {
				// defines if is verse or chorus
				// TODO: create function return if is verse or chorus
				if len(value) > 1 {
					numberVerse, _ = strconv.Atoi(reIniDig.FindString(value))
					if isInitVerse && numberVerse == 0 {
						isChorus = true
						cChorus++
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
					chorus = append(chorus, s.Trim(value, " "))
				}

			}

			if init == idx {
				fmt.Println("-----------------------------")
			}
			i1 := len(value) == 0
			fmt.Printf("-%d | %d | N.v %d | chorus %t | qt c %d | verse %t | title %t | eol %t\n", idx, n, numberVerse, isChorus, cChorus, isInitVerse, isTitle, i1)
			// h := new(Hymn)
			// m := re.FindStringSubmatch(value)
			// n, _ = strconv.Atoi(s.Replace(m[1], ".", "", -1))
			// h.Number = n
			// h.Title = s.Trim(m[2], " ")

			// hymns[h.Number] = h

			if n != u {
				u = n
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
