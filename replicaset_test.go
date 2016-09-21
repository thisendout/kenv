package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
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

	s, err := InjectVarsReplicaSet(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	if s == "" {
		t.Fatalf("got empty string")
	}
}

func TestUpdateContainersReplicaSet(t *testing.T) {
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

	replicaSet := v1beta1.ReplicaSet{}
	if err := json.Unmarshal(data, &replicaSet); err != nil {
		t.Fatal(err)
	}

	updateContainersReplicaSet(&replicaSet, envVars)

	for _, c := range replicaSet.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}
