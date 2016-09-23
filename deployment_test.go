package main

import (
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/v1"
)

func TestInjectVarsDeployment(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/deployment.json")
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

	deployment, err := InjectVarsDeployment(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range deployment.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}
