package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/util/validation"
)

// Var represents a basic key/value variable
type Var struct {
	Key   string
	Value string
}

// Vars is a Var slice
type Vars []Var

// NewVarsFromFiles takes a slice of files and returns a Vars struct
func newVarsFromFiles(files []string) (Vars, error) {
	vars := Vars{}

	// read in vars files
	for _, filename := range files {
		var v Vars
		var err error

		if path.Ext(filename) == ".yml" || path.Ext(filename) == ".yaml" {
			v, err = readYAMLVars(filename)
		} else {
			v, err = readKVVars(filename)
		}

		if err != nil {
			return vars, err
		}

		vars = append(vars, v...)
	}

	return vars, nil
}

// readKVVars reads files in "key=value" format and returns Vars
func readKVVars(filename string) (Vars, error) {
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

// readYAMLVars reads files in "key: value" format and returns Vars
func readYAMLVars(filename string) (Vars, error) {
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

func (vars Vars) toConfigMap(name string, namespace string, convert bool) ([]v1.EnvVar, *v1.ConfigMap, error) {
	envVars := []v1.EnvVar{}
	data := make(map[string]string)

	for _, v := range vars {
		key := v.Key

		if convertKeys {
			key = strings.ToLower(strings.Replace(key, "_", "-", -1))
			errs := validation.IsDNS1123Subdomain(key)
			if len(errs) > 0 {
				err := fmt.Errorf("%s is not a valid ConfigMap key: %s", v.Key, strings.Join(errs, ", "))
				return envVars, &v1.ConfigMap{}, err
			}
		} else {
			errs := validation.IsConfigMapKey(v.Key)
			if len(errs) > 0 {
				err := fmt.Errorf("%s is not a valid ConfigMap key: %s", v.Key, strings.Join(errs, ", "))
				return envVars, &v1.ConfigMap{}, err
			}
		}

		data[key] = v.Value

		envVars = append(envVars, v1.EnvVar{
			Name: v.Key,
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: name,
					},
					Key: key,
				},
			},
		})
	}

	configMap := &v1.ConfigMap{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	return envVars, configMap, nil
}
