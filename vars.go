package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
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

	// dirty hack for sorting
	keys := []string{}
	for k, _ := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vars = append(vars, Var{
			Key:   k,
			Value: config[k],
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
		key, err := validateKey(v.Key, convert)
		if err != nil {
			return envVars, &v1.ConfigMap{}, err
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

func (vars Vars) toSecret(name string, namespace string, convert bool) ([]v1.EnvVar, *v1.Secret, error) {
	envVars := []v1.EnvVar{}
	data := make(map[string][]byte)

	for _, v := range vars {
		key, err := validateKey(v.Key, convert)
		if err != nil {
			return envVars, &v1.Secret{}, err
		}

		data[key] = []byte(v.Value)

		envVars = append(envVars, v1.EnvVar{
			Name: v.Key,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: name,
					},
					Key: key,
				},
			},
		})
	}

	secret := &v1.Secret{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	return envVars, secret, nil
}

func validateKey(key string, convert bool) (string, error) {
	if convert {
		convertedKey := strings.ToLower(strings.Replace(key, "_", "-", -1))

		errs := validation.IsDNS1123Subdomain(convertedKey)
		if len(errs) > 0 {
			err := fmt.Errorf("%s is not a valid ConfigMap key: %s", key, strings.Join(errs, ", "))
			return "", err
		}

		return convertedKey, nil
	}

	errs := validation.IsConfigMapKey(key)
	if len(errs) > 0 {
		err := fmt.Errorf("%s is not a valid ConfigMap key: %s", key, strings.Join(errs, ", "))
		return "", err
	}

	return key, nil
}
