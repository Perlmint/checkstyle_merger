//+build !test

package main

import (
 	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/blang/semver"
)

func check(e error, i string) {
	if e != nil {
		fmt.Println(i)
		panic(e)
	}
}

func parseInput(path string) (CheckStyle, *semver.Version) {
	data, err := ioutil.ReadFile(path)
	check(err, path)
	parsed := CheckStyle{}
	err = xml.Unmarshal(data, &parsed)
	check(err, path)
	ver, err := semver.Make(parsed.Version)

	if err == nil {
		return parsed, &ver
	} else {
		return parsed, nil
	}
}

func main() {
	result := CheckStyle{}
	outPath := flag.String("o", "-", "Output path to combined checkstyle xml.")
	basePath := flag.String("b", "", "Relative path base")
	relativeBase := ""
	if *basePath != "" {
		relativeBase, _ = filepath.Abs(*basePath)
	} else {
		relativeBase, _ = os.Getwd()
	}

	var modifier func(file *File)

	if relativeBase != "" {
		modifier = makeRelativeModifier(relativeBase)
	} else {
		modifier = emptyModifier
	}

	flag.Parse()

	inputs := flag.Args()
	if len(inputs) == 0 {
		fmt.Fprintf(os.Stderr, "checkstyle_merger needs at least one input")
		os.Exit(1)
	}

	var vers []semver.Version

	for _, input := range inputs {
		parsed, ver := parseInput(input)
		if ver != nil {
			vers = append(vers, *ver)
		}

		mergeData(parsed, &result, modifier)
	}

	sortCheckStyle(&result)

	ver, w := getProperVersion(vers)
	if w != nil {
		fmt.Fprintf(os.Stderr, "%q", w)
	}

	result.Version = ver.String()

	out, err := xml.Marshal(result)
	check(err, "out")
	header := []byte(xml.Header)
	out = append(header, out...)
	if *outPath != "-" {
		ioutil.WriteFile(*outPath, out, 0644)
	} else {
		fmt.Printf("%s", out)
	}
}
