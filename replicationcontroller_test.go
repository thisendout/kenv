package main

import (
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/v1"
)

func TestInjectVarsReplicationController(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/replicationcontroller.json")
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

	replicationController, err := InjectVarsReplicationController(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range replicationController.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}
