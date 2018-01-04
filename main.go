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
		for k, v := range limVerse {
			fmt.Printf("nh: %3d ", k)
			for k, v := range v {
				fmt.Printf("ns: %3d | s: %3d e: %3d\n", k, v.start, v.end)
			}
		}

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
		if re.MatchString(value) { // title hymn
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
	var nhymn, nstrf, lnstrf, li int // n: number hymn, nv: verse number, lnv: last number, li: last index from verse
	lim := make(map[int]map[int]*position)

	for i, value := range lines {
		if re.MatchString(value) {
			nhymn, _ = getNumberTitleHymn(value)
			//lim[nhymn] = make(map[int]*position)
		} else { // somewhere strophes
			if len(value) > 1 {
				nstrf = getNumberVerse(value)
				if nstrf > 0 {

					if lnstrf != nstrf {
						fmt.Printf(" nh %3d | nstrf %d | lnstrf %d | i %3d | li %3d\n", nhymn, nstrf, lnstrf, i, li)

						// Novo mapa de estofes
						mms := make(map[int]*position)

						if v, ok := mms[lnstrf]; ok {
							v.end = i
						}

						if _, ok := mms[nstrf]; !ok {
							mms[nstrf] = &position{start: i}
						}
						lim[nhymn] = mms

						// verifico se o hino tem estrofe
						//
						// for k, v := range mm {
						// 	fmt.Printf("%d %d %d \n", k, v.start, v.end)
						// }

						// if !ok { // se não tem estrofe crio uma novo mapa de estrofes
						// 	mm = make(map[int]*position)
						// 	mm[nstrf] = &position{start: i} // crio um mapa com a estrofe atual
						// 	lim[nhymn] = mm
						// }

						// else {
						// 	// o hino ja tem estrofe
						// 	// verifico a estrofe ja existe que esta em `mm`

						// 	fmt.Println(mm)
						// 	if v, ok := mm[lnstrf]; ok {
						// 		v.end = i
						// 	}
						// }

						// last index verse number (liv)
						if li != i {
							li = i
						}

						// strophe number
						lnstrf = nstrf
					}
				}
			}
		}

		// last line file
		// if len(lines)-1 == i {
		// 	if v, ok := lim[n]; ok {
		// 		v.end = len(lines)
		// 	}
		// }
	}
	return lim
}

/*

 nh 113 | nstrf 1 | lnstrf 0 | i   2 | li   0
 nh 113 | nstrf 2 | lnstrf 1 | i  13 | li   2
 nh 113 | nstrf 3 | lnstrf 2 | i  24 | li  13
 nh 114 | nstrf 1 | lnstrf 3 | i  36 | li  24
 nh 114 | nstrf 2 | lnstrf 1 | i  46 | li  36
 nh 114 | nstrf 3 | lnstrf 2 | i  51 | li  46
 nh 114 | nstrf 4 | lnstrf 3 | i  56 | li  51
 nh 115 | nstrf 1 | lnstrf 4 | i  63 | li  56
 nh 115 | nstrf 2 | lnstrf 1 | i  75 | li  63
 nh 115 | nstrf 3 | lnstrf 2 | i  81 | li  75
 nh 116 | nstrf 1 | lnstrf 3 | i  89 | li  81
 nh 116 | nstrf 2 | lnstrf 1 | i 100 | li  89
 nh 116 | nstrf 3 | lnstrf 2 | i 106 | li 100
 nh 117 | nstrf 1 | lnstrf 3 | i 114 | li 106
 nh 117 | nstrf 2 | lnstrf 1 | i 123 | li 114
 nh 117 | nstrf 3 | lnstrf 2 | i 128 | li 123
 nh 117 | nstrf 4 | lnstrf 3 | i 133 | li 128
 nh 168 | nstrf 1 | lnstrf 4 | i 140 | li 133
 nh 168 | nstrf 2 | lnstrf 1 | i 153 | li 140
 nh 168 | nstrf 3 | lnstrf 2 | i 161 | li 153
 nh 169 | nstrf 1 | lnstrf 3 | i 171 | li 161
 nh 169 | nstrf 2 | lnstrf 1 | i 181 | li 171
 nh 169 | nstrf 3 | lnstrf 2 | i 186 | li 181
 nh 169 | nstrf 4 | lnstrf 3 | i 191 | li 186
 nh 170 | nstrf 1 | lnstrf 4 | i 198 | li 191
 nh 170 | nstrf 2 | lnstrf 1 | i 208 | li 198
 nh 170 | nstrf 3 | lnstrf 2 | i 213 | li 208
*/
