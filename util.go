package main

import "k8s.io/kubernetes/pkg/api/v1"

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
