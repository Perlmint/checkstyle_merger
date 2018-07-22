package main

import (
 	"encoding/xml"
	"fmt"
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

func sortErrors(file *File) {
	sort.Slice(file.Errors, func(i, j int) bool {
		e1, e2 := file.Errors[i], file.Errors[j]
		if e1.Line == e2.Line {
			return e1.Column < e2.Column
		}
		return e1.Line < e2.Line
	})
}

func mergeData(src CheckStyle, dst *CheckStyle, modifier func(*File)) {
	for _, file := range src.Files {
		modifier(&file)

		found := false
		for idx, _ := range dst.Files {
			if dst.Files[idx].Name == file.Name {
				dst.Files[idx].Errors = append(dst.Files[idx].Errors, file.Errors...)
				found = true
				break
			}
		}

		if found == false {
			dst.Files = append(dst.Files, file)
		}
	}
}

func sortCheckStyle(checkstyle *CheckStyle) {
	sortFiles(checkstyle)

	for _, file := range checkstyle.Files {
		sortErrors(&file)
	}
}

func sortFiles(checkstyle *CheckStyle) {
	sort.Slice(checkstyle.Files, func (i, j int) bool {
		return checkstyle.Files[i].Name < checkstyle.Files[j].Name
	})
}

func makeRelativeModifier(relativeBase string) func(file *File) {
	return func (file *File) {
		if filepath.IsAbs(file.Name) {
			file.Name, _ = filepath.Rel(relativeBase, file.Name)
		}
	}
}

func emptyModifier(file *File) {
	// do nothing
}

func getProperVersion(vers []semver.Version) (semver.Version, *error) {
	if len(vers) == 0 {
		ver, _ := semver.Make(DEFAULT_VER)
		e := fmt.Errorf("There is no valid version. use checkstyl_merger default value %q", DEFAULT_VER)

		return ver, &e
	} else {
		sort.Slice(vers, func(i, j int) bool {
			return vers[i].GT(vers[j])
		})

		return  vers[0], nil
	}
}
