package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

func TestInjectVarsDaemonSet(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/daemonset.json")
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

	s, err := InjectVarsDaemonSet(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	if s == "" {
		t.Fatalf("got empty string")
	}
}

func TestUpdateContainersDaemonSet(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/daemonset.json")
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

	daemonSet := v1beta1.DaemonSet{}
	if err := json.Unmarshal(data, &daemonSet); err != nil {
		t.Fatal(err)
	}

	updateContainersDaemonSet(&daemonSet, envVars)

	for _, c := range daemonSet.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}
