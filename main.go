package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	for _, input := range flag.Args() {
		data, err := ioutil.ReadFile(input)
		check(err, input)
		parsed := CheckStyle{}
		err = xml.Unmarshal(data, &parsed)
		check(err, input)
		for _, file := range parsed.Files {
			if relativeBase != "" && filepath.IsAbs(file.Name) {
				file.Name, _ = filepath.Rel(relativeBase, file.Name)
			}
			result.Files = append(result.Files, file)
		}
	}

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
