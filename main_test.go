package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestMainWithYAML(t *testing.T) {
	os.Args = []string{
		"kenv",
		"-v",
		"fixtures/vars.env",
		"fixtures/deployment.yaml",
	}

	main()
}

func TestMainWithJSON(t *testing.T) {
	os.Args = []string{
		"kenv",
		"-v",
		"fixtures/vars.env",
		"fixtures/deployment.json",
	}

	main()
}

func TestFlagSliceString(t *testing.T) {
	fs := FlagSlice{"foo", "bar"}
	if fs.String() != "foo,bar" {
		t.Fatalf("strings not equal")
	}
}

func TestFlagSliceSet(t *testing.T) {
	want := FlagSlice{"foo", "bar"}

	fs := FlagSlice{"foo"}
	fs.Set("bar")

	if !reflect.DeepEqual(want, fs) {
		t.Fatalf("flagslice not equal")
	}
}

func TestGetDocKing(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/deployment.json")
	if err != nil {
		t.Fatal(err)
	}

	kind, err := getDocKind(data)
	if err != nil {
		t.Fatal(err)
	}

	if kind != "Deployment" {
		t.Fatalf("kind not equal")
	}
}
