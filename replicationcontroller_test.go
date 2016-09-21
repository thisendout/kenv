package main

import (
	"encoding/json"
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

	s, err := InjectVarsReplicationController(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	if s == "" {
		t.Fatalf("got empty string")
	}
}

func TestUpdateContainersReplicationController(t *testing.T) {
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

	replicationController := v1.ReplicationController{}
	if err := json.Unmarshal(data, &replicationController); err != nil {
		t.Fatal(err)
	}

	updateContainersReplicationController(&replicationController, envVars)

	for _, c := range replicationController.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}
