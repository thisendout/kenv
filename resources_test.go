package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/v1"
)

func TestParseDocs(t *testing.T) {
	file, err := os.Open("fixtures/deployment-service.yml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := ParseDocs(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(resources) != 2 {
		t.Fatalf("Resource count is not 2: %+v", resources)
	}
}

func TestInjectVarsDeployment(t *testing.T) {
	file, err := os.Open("fixtures/deployment.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := ParseDocs(file)
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

	deployment, err := resources[0].InjectVarsDeployment(envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range deployment.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}

func TestInjectVarsDaemonSet(t *testing.T) {
	file, err := os.Open("fixtures/deployment.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := ParseDocs(file)
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

	daemonSet, err := resources[0].InjectVarsDaemonSet(envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range daemonSet.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}

func TestInjectVarsReplicaSet(t *testing.T) {
	file, err := os.Open("fixtures/deployment.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := ParseDocs(file)
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

	replicaSet, err := resources[0].InjectVarsReplicaSet(envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range replicaSet.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}

func TestInjectVarsRC(t *testing.T) {
	file, err := os.Open("fixtures/deployment.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := ParseDocs(file)
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

	replicationController, err := resources[0].InjectVarsRC(envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range replicationController.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}

func TestUnmarshalGeneric(t *testing.T) {
	file, err := os.Open("fixtures/deployment.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := ParseDocs(file)
	if err != nil {
		t.Fatal(err)
	}

	_, err = resources[0].UnmarshalGeneric()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetResourceKind(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/deployment.json")
	if err != nil {
		t.Fatal(err)
	}

	kind, err := getResourceKind(data)
	if err != nil {
		t.Fatal(err)
	}

	if kind != "Deployment" {
		t.Fatalf("kind not equal")
	}
}

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
