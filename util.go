package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
)

// creates a flattened EnvVar slice giving preference to user supplied vars
func mergeEnvVars(docVars []v1.EnvVar, userVars []v1.EnvVar) []v1.EnvVar {
	mergedVars := userVars
	for _, v := range docVars {
		if !isDuplicateEnvVar(v, userVars) {
			mergedVars = append(mergedVars, v)
		}
	}

	return mergedVars
}

// checks whether an EnvVar exists by name in an EnvVar slice
func isDuplicateEnvVar(e v1.EnvVar, envVars []v1.EnvVar) bool {
	for _, envVar := range envVars {
		if e.Name == envVar.Name {
			return true
		}
	}
	return false
}

// injectPodSpecEnvVars injects a slice of EnvVars into each PodSpec container
func injectPodSpecEnvVars(podSpec v1.PodSpec, envVars []v1.EnvVar) v1.PodSpec {
	containers := []v1.Container{}
	for _, c := range podSpec.Containers {
		c.Env = mergeEnvVars(c.Env, envVars)
		containers = append(containers, c)
	}
	podSpec.Containers = containers
	return podSpec
}

// newList creates a new list of resources to create
func newList(objects ...runtime.Object) v1.List {
	rawExtensions := []runtime.RawExtension{}
	for _, o := range objects {
		rawExtensions = append(rawExtensions, runtime.RawExtension{
			Object: o,
		})
	}

	return v1.List{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		Items: rawExtensions,
	}
}

// getDocKind unmarshalls a file and returns the kind of resource doc
func getDocKind(data []byte) (string, error) {
	typeMeta := unversioned.TypeMeta{}
	if err := json.Unmarshal(data, &typeMeta); err != nil {
		return "", err
	}
	return typeMeta.Kind, nil
}
