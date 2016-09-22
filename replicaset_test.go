package main

import (
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/v1"
)

func TestInjectVarsReplicaSet(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/replicaset.json")
	if err != nil {
		t.Fatal(err)
	}

	envVars := []v1.EnvVar{
		v1.EnvVar{
			Name:  "key1",
			Value: "value1",
		},
		v1.EnvVar{
			Name:  "key2",
			Value: "value2",
		},
	}

	replicaSet, err := InjectVarsReplicaSet(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range replicaSet.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}
