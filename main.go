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
		Verse  []string
	}

	position struct {
		start int
		end   int
	}
)

var (
	hymns    map[int]Hymn
	limHymn  map[int]*position
	limVerse map[int]map[int]*position
	re       = regexp.MustCompile(`(?P<Number>[0-9]+\.) (?P<Title>[a-zA-Zà-úÀ-Ú0-9 \,\']+)`)
	reIniDig = regexp.MustCompile(`^[1-9]`)

	//reInsideWhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
)

func main() {

	d, _ := os.Open(sourcePath)
	files, _ := d.Readdir(-1)

	for _, fi := range files {
		fmt.Println(fi.Name())

		filePath := sourcePath + "/" + fi.Name()

		lines := Readlines(filePath)

		// Delimited the start and end each of hymn
		limHymn = delimetedHymn(lines)
		limVerse = delimetedVerse(lines)
		// for k, v := range limHymn {
		// 	fmt.Printf("N: %3d - s: %5d e: %5d\n", k, v.start, v.end)
		// }
		// for k, v := range limVerse {
		// 	fmt.Printf("nh: %3d ", k)
		// 	for k, v := range v {
		// 		fmt.Printf("ns: %3d | s: %3d e: %3d\n", k, v.start, v.end)
		// 	}
		// }

		//
		var h Hymn
		var vs Verse
		var numberHymn, numberVerse int
		var isEndVerse, isChorus bool
		var title string
		var chorus, strophe []string
		// var strophes []Verse

		hymns = make(map[int]Hymn)

		for idx, value := range lines {

			isTitle := re.MatchString(value)
			if isTitle {
				isChorus = false
				chorus = nil
				numberHymn, title = getNumberTitleHymn(value)
				//title = s.Trim(m[2], " ")
			} else {
				// defines if is verse or chorus
				// TODO: create function return if is verse or chorus
				if len(value) > 1 {
					numberVerse = getNumberVerse(value)
					if isEndVerse && numberVerse == 0 {
						isChorus = true
					}
					isEndVerse = false
				} else {
					isEndVerse = true
					isChorus = false
				}

				// Is Chorus
				if len(value) > 0 && isChorus {
					// TODO: remover os tab \t
					chorus = append(chorus, s.Trim(value, " "))
				}

			}

			// get index the hymn
			if v, ok := limHymn[numberHymn]; ok {
				if v.start == idx { // create
					h = Hymn{Number: numberHymn, Title: title}
					// fmt.Printf("n-> %d ----------------------\n", numberHymn) *****
					vs = Verse{Number: 0, Verse: nil}
				}

				if v.start <= idx+1 && v.end >= idx+1 {
					if !isChorus && !isTitle {
						if len(value) > 1 {
							strophe = append(strophe, s.Trim(value, " "))
						} else {
							if len(strophe) > 0 {
								vs.Verse = strophe
								strophe = nil
							}
						}
						if numberVerse > 0 {
							vs.Number = numberVerse
						}
					}
				}

				if v.end == idx+1 { // add map
					h.Chorus = chorus
					hymns[numberHymn] = h
				}
			}

		}

		// for k, v := range hymns {
		// 	fmt.Printf("id: %d - %v\n", k, v)
		// }

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

func getNumberTitleHymn(value string) (int, string) {
	m := re.FindStringSubmatch(value)
	n, _ := strconv.Atoi(s.Replace(m[1], ".", "", -1))
	t := s.Trim(m[2], " ")
	return n, t
}

func getNumberVerse(value string) int {
	n, _ := strconv.Atoi(reIniDig.FindString(value))
	return n
}

func delimetedHymn(lines []string) map[int]*position {
	var n, ln, li int // n: number , ln: last number, li: last index
	mts := make(map[int]*position)
	for i, value := range lines {
		if re.MatchString(value) { // hymn title
			n, _ = getNumberTitleHymn(value)

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

func delimetedVerse(lines []string) map[int]map[int]*position {
	var nhymn, lnhymn, nstrf, lnstrf, li, iblank, liblank, c int
	lim := make(map[int]map[int]*position)

	for i, value := range lines {
		if re.MatchString(value) {
			nhymn, _ = getNumberTitleHymn(value)
			c = 0
		} else { // somewhere strophes
			if len(value) > 1 {
				nstrf = getNumberVerse(value)
				if nstrf > 0 {
					if lnstrf != nstrf {
						//fmt.Printf(" nh %3d | lnh %3d | nstrf %d | lnstrf %d | i %3d | li %3d | ib %d | lb %d | c %d \n", nhymn, lnhymn, nstrf, lnstrf, i, li, iblank, liblank, c)

						// get the last strophe of hymn previous
						if m, ok := lim[lnhymn][lnstrf]; ok {
							if nstrf == 1 {
								m.end = li + 1
							} else {
								m.end = iblank
							}
						}

						// Add new delimeted inde strophe
						add(lim, nhymn, nstrf, i+1, 0)

						// avoid count index of chorus in first strophe
						if lnstrf == 1 {
							if m, ok := lim[nhymn][1]; ok {
								if liblank > 1 {
									m.end = liblank
								}
							}
						}

						// strophe number
						lnstrf = nstrf

						// last hymn number
						if lnhymn != nhymn {
							lnhymn = nhymn
						}
					}
				} else {
					if liblank != iblank {
						c++
						liblank = iblank
					}
					// last index verse number (liv)
					if li != i {
						li = i
					}
				}
			} else {
				iblank = i
			}
		}

		// last line file from strophe previous
		if len(lines)-1 == i {
			if v, ok := lim[nhymn][lnstrf]; ok {
				v.end = i
			}
		}
	}

	for k, v := range lim {
		fmt.Printf("NH: %3d \n", k)
		for k, v := range v {
			fmt.Printf(" ns: %3d | s: %3d e: %3d\n", k, v.start, v.end)
		}
	}

	return lim
}

func add(m map[int]map[int]*position, hymn, verse, start, end int) {
	mm, ok := m[hymn]
	if !ok {
		mm = make(map[int]*position)
		m[hymn] = mm
	}
	mm[verse] = &position{start, end}
}
