package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Var represents a basic key/value variable
type Var struct {
	Key   string
	Value string
}

// Vars is a Var slice
type Vars []Var

// ReadVarsFile determines the filetype and
// reads the contents into Vars
func ReadVarsFile(filename string) (Vars, error) {
	if path.Ext(filename) == ".yml" || path.Ext(filename) == ".yaml" {
		return ReadYAMLFile(filename)
	}

	return ReadKVFile(filename)
}

// ReadKVFile reads files in "key=value" format and returns Vars
func ReadKVFile(filename string) (Vars, error) {
	vars := Vars{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return vars, err
	}

	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		if l == "" {
			continue
		}

		lSplit := strings.Split(l, "=")

		if len(lSplit) < 2 {
			fmt.Printf("Skipping %s; not in key=value format\n", l)
			continue
		}

		vars = append(vars, Var{
			Key:   lSplit[0],
			Value: strings.Join(lSplit[1:], "="),
		})
	}

	return vars, nil
}

// ReadYAMLFile reads files in "key: value" format and returns Vars
func ReadYAMLFile(filename string) (Vars, error) {
	vars := Vars{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return vars, err
	}

	config := make(map[string]string)

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return vars, err
	}

	for k, v := range config {
		vars = append(vars, Var{
			Key:   k,
			Value: v,
		})
	}

	return vars, nil
}

func (vars Vars) toStringMap() map[string]string {
	m := make(map[string]string)
	for _, v := range vars {
		m[v.Key] = v.Value
	}

	return m
}

func (vars Vars) toEnvVar() []v1.EnvVar {
	envVars := []v1.EnvVar{}
	for _, v := range vars {
		envVars = append(envVars, v1.EnvVar{
			Name:  v.Key,
			Value: v.Value,
		})
	}

	return envVars
}

func (vars Vars) toConfigMap(name string, namespace string) ([]v1.EnvVar, v1.ConfigMap) {
	envVars := []v1.EnvVar{}

	configMap := v1.ConfigMap{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: vars.toStringMap(),
	}

	for _, v := range vars {
		envVars = append(envVars, v1.EnvVar{
			Name: v.Key,
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: name,
					},
					Key: v.Key,
				},
			},
		})
	}

	return envVars, configMap
}
