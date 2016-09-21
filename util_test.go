package main

import (
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
