package main

import (
	"encoding/base64"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
)

func TestReadVarsFromFiles(t *testing.T) {
	want := Vars{
		Var{
			Key:   "KVKey1",
			Value: "KVValue1",
		},
		Var{
			Key:   "kvkey2",
			Value: "kvvalue2",
		},
		Var{
			Key:   "YAMLKey1",
			Value: "YAMLValue1",
		},
		Var{
			Key:   "yamlkey2",
			Value: "yamlvalue2",
		},
	}

	vars, err := newVarsFromFiles([]string{
		"fixtures/vars.env",
		"fixtures/vars.yaml",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestReadYAMLVars(t *testing.T) {
	want := Vars{
		Var{
			Key:   "YAMLKey1",
			Value: "YAMLValue1",
		},
		Var{
			Key:   "yamlkey2",
			Value: "yamlvalue2",
		},
	}

	vars, err := readYAMLVars("fixtures/vars.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestReadKVVars(t *testing.T) {
	want := Vars{
		Var{
			Key:   "KVKey1",
			Value: "KVValue1",
		},
		Var{
			Key:   "kvkey2",
			Value: "kvvalue2",
		},
	}

	vars, err := readKVVars("fixtures/vars.env")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, vars) {
		t.Fatalf("not equal, wanted: %+v, got: %+v", want, vars)
	}
}

func TestToEnvVar(t *testing.T) {
	want := []v1.EnvVar{
		v1.EnvVar{
			Name:  "KVKey1",
			Value: "KVValue1",
		},
		v1.EnvVar{
			Name:  "kvkey2",
			Value: "kvvalue2",
		},
	}

	v := Vars{
		Var{
			Key:   "KVKey1",
			Value: "KVValue1",
		},
		Var{
			Key:   "kvkey2",
			Value: "kvvalue2",
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
			Name: "KVKey1",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "foo",
					},
					Key: "KVKey1",
				},
			},
		},
		v1.EnvVar{
			Name: "kvkey2",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "foo",
					},
					Key: "kvkey2",
				},
			},
		},
	}

	wantConfigMap := &v1.ConfigMap{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
		Data: map[string]string{
			"KVKey1": "KVValue1",
			"kvkey2": "kvvalue2",
		},
	}

	v := Vars{
		Var{
			Key:   "KVKey1",
			Value: "KVValue1",
		},
		Var{
			Key:   "kvkey2",
			Value: "kvvalue2",
		},
	}

	envVars, configMap, err := v.toConfigMap("foo", "bar", false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantEnvVars, envVars) {
		t.Fatalf("EnvVars not equal, wanted: %+v, got: %+v", wantEnvVars, envVars)
	}

	if !reflect.DeepEqual(wantConfigMap, configMap) {
		t.Fatalf("ConfigMap not equal, wanted: %+v, got: %+v", wantConfigMap, configMap)
	}
}

func TestToSecret(t *testing.T) {
	wantEnvVars := []v1.EnvVar{
		v1.EnvVar{
			Name: "KVKey1",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "foo",
					},
					Key: "KVKey1",
				},
			},
		},
		v1.EnvVar{
			Name: "kvkey2",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "foo",
					},
					Key: "kvkey2",
				},
			},
		},
	}

	wantSecret := &v1.Secret{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
		Data: map[string][]byte{
			"KVKey1": []byte(base64.StdEncoding.EncodeToString([]byte("KVValue1"))),
			"kvkey2": []byte(base64.StdEncoding.EncodeToString([]byte("kvvalue2"))),
		},
	}

	v := Vars{
		Var{
			Key:   "KVKey1",
			Value: "KVValue1",
		},
		Var{
			Key:   "kvkey2",
			Value: "kvvalue2",
		},
	}

	envVars, secret, err := v.toSecret("foo", "bar", false)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantEnvVars, envVars) {
		t.Fatalf("EnvVars not equal, wanted: %+v, got: %+v", wantEnvVars, envVars)
	}

	if !reflect.DeepEqual(wantSecret, secret) {
		t.Fatalf("Secret not equal, wanted: %+v, got: %+v", wantSecret, secret)
	}
}

func TestValidateKey(t *testing.T) {
	_, err := validateKey("foo", false)
	if err != nil {
		t.Fatal(err)
	}

	// test valid key for kubernetes >= 1.4
	_, err = validateKey("FOO.BAR_BANANA", false)
	if err != nil {
		t.Fatal(err)
	}

	// test valid key for kubernetes < 1.4 with conversion
	_, err = validateKey("FOO.BAR_BANANA", true)
	if err != nil {
		t.Fatal(err)
	}

	// test invalid key without conversion
	_, err = validateKey("@@@@@@@@", false)
	if err == nil {
		t.Fatal("expecting error")
	}

	// test invalid key with conversion
	_, err = validateKey("@@@@@@@@", true)
	if err == nil {
		t.Fatal("expecting error")
	}
}
