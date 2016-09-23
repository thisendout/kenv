package main

import (
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/v1"
)

func TestMergeEnvVars(t *testing.T) {
	merged := mergeEnvVars([]v1.EnvVar{
		v1.EnvVar{
			Name:  "dup",
			Value: "base",
		},
		v1.EnvVar{
			Name:  "key1",
			Value: "value1",
		},
	}, []v1.EnvVar{
		v1.EnvVar{
			Name:  "dup",
			Value: "overwrite",
		},
		v1.EnvVar{
			Name:  "key2",
			Value: "value2",
		},
	})

	want := []v1.EnvVar{
		v1.EnvVar{
			Name:  "dup",
			Value: "overwrite",
		},
		v1.EnvVar{
			Name:  "key2",
			Value: "value2",
		},
		v1.EnvVar{
			Name:  "key1",
			Value: "value1",
		},
	}

	if !reflect.DeepEqual(want, merged) {
		t.Fatalf("slices not equal; want: %+v, got: %+v", want, merged)
	}
}

func TestIsDuplicateEnvVar(t *testing.T) {
	d := isDuplicateEnvVar(v1.EnvVar{
		Name:  "dup",
		Value: "base",
	}, []v1.EnvVar{
		v1.EnvVar{
			Name:  "dup",
			Value: "overwrite",
		},
	})

	if !d {
		t.Fatalf("should be duplicate")
	}
}

func TestGetDocKind(t *testing.T) {
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
