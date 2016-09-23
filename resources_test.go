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

	deployment, err := injectVarsDeployment(data, envVars)
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

	daemonSet, err := injectVarsDaemonSet(data, envVars)
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

	replicaSet, err := injectVarsReplicaSet(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range replicaSet.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
	}
}

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

	replicationController, err := injectVarsReplicationController(data, envVars)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range replicationController.Spec.Template.Spec.Containers {
		if !reflect.DeepEqual(c.Env, envVars) {
			t.Fatalf("container env vars not equal")
		}
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
