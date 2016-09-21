package main

import (
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
)

func TestReadVarsFileKV(t *testing.T) {
	want := Vars{
		Var{
			Key:   "KEY1",
			Value: "VALUE1",
		},
		Var{
			Key:   "key2",
			Value: "value2",
		},
	}

	vars, err := ReadVarsFile("fixtures/vars.env")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestReadVarsFileYAML(t *testing.T) {
	want := Vars{
		Var{
			Key:   "KEY1",
			Value: "VALUE1",
		},
	}

	vars, err := ReadVarsFile("fixtures/vars.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestReadKVFile(t *testing.T) {
	want := Vars{
		Var{
			Key:   "KEY1",
			Value: "VALUE1",
		},
		Var{
			Key:   "key2",
			Value: "value2",
		},
	}

	vars, err := ReadKVFile("fixtures/vars.env")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestReadYAMLFile(t *testing.T) {
	want := Vars{
		Var{
			Key:   "KEY1",
			Value: "VALUE1",
		},
	}

	vars, err := ReadYAMLFile("fixtures/vars.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestToStringMap(t *testing.T) {
	want := map[string]string{
		"KEY1": "VALUE1",
		"key2": "value2",
	}

	v := Vars{
		Var{
			Key:   "KEY1",
			Value: "VALUE1",
		},
		Var{
			Key:   "key2",
			Value: "value2",
		},
	}

	vars := v.toStringMap()

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestToEnvVar(t *testing.T) {
	want := []v1.EnvVar{
		v1.EnvVar{
			Name:  "KEY1",
			Value: "VALUE1",
		},
		v1.EnvVar{
			Name:  "key2",
			Value: "value2",
		},
	}

	v := Vars{
		Var{
			Key:   "KEY1",
			Value: "VALUE1",
		},
		Var{
			Key:   "key2",
			Value: "value2",
		},
	}

	vars := v.toEnvVar()

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestToConfigMap(t *testing.T) {
	wantEnvVars := []v1.EnvVar{
		v1.EnvVar{
			Name: "KEY1",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "foo",
					},
					Key: "KEY1",
				},
			},
		},
		v1.EnvVar{
			Name: "key2",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "foo",
					},
					Key: "key2",
				},
			},
		},
	}

	wantConfigMap := v1.ConfigMap{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
		Data: map[string]string{
			"KEY1": "VALUE1",
			"key2": "value2",
		},
	}

	v := Vars{
		Var{
			Key:   "KEY1",
			Value: "VALUE1",
		},
		Var{
			Key:   "key2",
			Value: "value2",
		},
	}

	envVars, configMap := v.toConfigMap("foo", "bar")

	if !reflect.DeepEqual(wantEnvVars, envVars) {
		t.Fatalf("EnvVars not equal, wanted: %+v, got: %+v", wantEnvVars, envVars)
	}

	if !reflect.DeepEqual(wantConfigMap, configMap) {
		t.Fatalf("ConfigMap not equal, wanted: %+v, got: %+v", wantConfigMap, configMap)
	}
}
