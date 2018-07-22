package main

import (
	"testing"
	"github.com/blang/semver"
	"github.com/go-test/deep"
)

func TestVersionSelectionBasic(t *testing.T) {
	vers := []semver.Version{}
	v, _ := semver.Make("1.0.0")
	vers = append(vers, v)
	v, _ = semver.Make("1.1.0")
	vers = append(vers, v)

	ver, _ := getProperVersion(vers)

	correct_ver, _ := semver.Make("1.1.0")

	if ver.Compare(correct_ver) != 0 {
		t.Errorf("getProperVersion should return maximum version, but returned %q", ver)
	}
}

func TestVerionSelectionFallback(t *testing.T) {
	vers := []semver.Version{}

	ver, _ := getProperVersion(vers)

	correct_ver, _ := semver.Make("3.0.0")

	if ver.Compare(correct_ver) != 0 {
		t.Errorf("getProperVersion should return fallback version, but returned %q", ver)
	}
}

func TestMerge(t *testing.T) {
	inputs := []CheckStyle{
		{
			Version: "1.0.0",
			Files: []File{
				{
					Name: "src.go",
					Errors: []Error{
						{
							Line: 1,
							Column: 1,
							Severity: "warning",
							Message: "error1",
							Source: "import \"asdf\"",
						},
					},
				},
			},
		},
		{
			Version: "2.0.0",
			Files: []File{
				{
					Name: "src.go",
					Errors: []Error{
						{
							Line: 2,
							Column: 20,
							Severity: "warning",
							Message: "error2",
							Source: "func aaaa",
						},
					},
				},
				{
					Name: "src2.go",
					Errors: []Error{
						{
							Line: 3,
							Column: 1,
							Severity: "warning",
							Message: "error1",
							Source: "impot \"qqwe\"",
						},
						{
							Line: 1,
							Column: 2,
							Severity: "warning",
							Message: "eee",
							Source: "aaaa",
						},
					},
				},
			},
		},
	}

	correct := CheckStyle{
		Version: "",
		Files: []File{
			{
				Name: "src.go",
				Errors: []Error{
					{
						Line: 1,
						Column: 1,
						Severity: "warning",
						Message: "error1",
						Source: "import \"asdf\"",
					},
					{
						Line: 2,
						Column: 20,
						Severity: "warning",
						Message: "error2",
						Source: "func aaaa",
					},
				},
			},
			{
				Name: "src2.go",
				Errors: []Error{
					{
						Line: 1,
						Column: 2,
						Severity: "warning",
						Message: "eee",
						Source: "aaaa",
					},
					{
						Line: 3,
						Column: 1,
						Severity: "warning",
						Message: "error1",
						Source: "impot \"qqwe\"",
					},
				},
			},
		},
	}

	ret := CheckStyle{}

	for _, input := range inputs {
		mergeData(input, &ret, emptyModifier)
	}

	sortCheckStyle(&correct)
	sortCheckStyle(&ret)

	if diff := deep.Equal(ret, correct); diff != nil {
		t.Error(diff)
	}
}

func TestModifier(t *testing.T) {
	file := File{
		Name: "/foo/bar.go",
	}

	emptyModifier(&file)

	if file.Name != "/foo/bar.go" {
		t.Error("emptyModifier does not works correctly")
	}

	modifier := makeRelativeModifier("/foo")
	modifier(&file)

	if file.Name != "bar.go" {
		t.Errorf("relativeModifier does not works correctly, actual value - %q", file.Name)
	}
}
