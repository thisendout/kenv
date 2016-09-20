package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
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

	s, err := InjectVarsDeployment(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	if s == "" {
		t.Fatalf("got empty string")
	}
}

func TestUpdateContainersDeployment(t *testing.T) {
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

	deployment := v1beta1.Deployment{}
	if err := json.Unmarshal(data, &deployment); err != nil {
		t.Fatal(err)
	}

	updateContainersDeployment(&deployment, envVars)

	for _, c := range deployment.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}
