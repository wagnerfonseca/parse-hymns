package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const sourcePath = "./raw"

func main() {

	d, _ := os.Open(sourcePath)
	files, _ := d.Readdir(-1)

	for _, fi := range files {
		fmt.Println(fi.Name())

		filePath := sourcePath + "/" + fi.Name()

		lines := Readlines(filePath)

		for _, value := range lines {
			//fmt.Printf("%d : %s\n", idx, value)
			re := regexp.MustCompile(`(?P<Number>[0-9]+\.) (?P<Title>[a-zA-Zà-úÀ-Ú0-9 \,\'])+`)

			if re.MatchString(value) {
				m := re.FindStringSubmatch(value)

				//Captura o titulo
				fmt.Println(m[1])

				fmt.Printf("%q\n", re.FindString(value))
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

// type struct Hymn {
// 	Number int
// 	Title string
// }
