package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// InjectVarsDaemonSet inserts EnvVars into a daemonSet doc
func InjectVarsDaemonSet(data []byte, envVars []v1.EnvVar) (string, error) {
	daemonSet := v1beta1.DaemonSet{}

	if err := json.Unmarshal(data, &daemonSet); err != nil {
		return "", err
	}

	updateContainersDaemonSet(&daemonSet, envVars)

	data, err := json.MarshalIndent(&daemonSet, "", "  ")

	return string(data), err
}

func updateContainersDaemonSet(daemonSet *v1beta1.DaemonSet, envVars []v1.EnvVar) {
	podSpec := daemonSet.Spec.Template.Spec
	containers := []v1.Container{}

	for _, c := range podSpec.Containers {
		c.Env = append(c.Env, envVars...)
		containers = append(containers, c)
	}

	podSpec.Containers = containers
	daemonSet.Spec.Template.Spec = podSpec
}
