package main

import (
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
