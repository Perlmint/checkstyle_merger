package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"path/filepath"
	"github.com/blang/semver"
)

var (
	DEFAULT_VER = "3.0.0"
)

type Error struct {
	Line     int32  `xml:"line,attr"`
	Column   int32  `xml:"column,attr"`
	Severity string `xml:"severity,attr"`
	Message  string `xml:"message,attr"`
	Source   string `xml:"source,attr"`
}

type File struct {
	Name   string  `xml:"name,attr"`
	Errors []Error `xml:"error"`
}

type CheckStyle struct {
	XMLName xml.Name `xml:"checkstyle"`
	Version string   `xml:"version,attr"`
	Files   []File   `xml:"file"`
}

func check(e error, i string) {
	if e != nil {
		fmt.Println(i)
		panic(e)
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

	flag.Parse()

	inputs := flag.Args()
	if len(inputs) == 0 {
		fmt.Fprintf(os.Stderr, "checkstyle_merger needs at least one input")
		os.Exit(1)
	}

	var vers []semver.Version

	for _, input := range inputs {
		data, err := ioutil.ReadFile(input)
		check(err, input)
		parsed := CheckStyle{}
		err = xml.Unmarshal(data, &parsed)
		check(err, input)
		ver, err := semver.Make(parsed.Version)
		if err == nil {
			vers = append(vers, ver)
		}
		for _, file := range parsed.Files {
			if relativeBase != "" && filepath.IsAbs(file.Name) {
				file.Name, _ = filepath.Rel(relativeBase, file.Name)
			}
			result.Files = append(result.Files, file)
		}
	}

	sort.Slice(vers, func(i, j int) bool {
		return vers[i].GT(vers[j])
	})

	var ver semver.Version
	if len(vers) == 0 {
		fmt.Fprintf(os.Stderr, "There is no valid version. use checkstyl_merger default value %q", DEFAULT_VER)
		ver, _ = semver.Make(DEFAULT_VER)
	} else {
		ver = vers[0]
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
