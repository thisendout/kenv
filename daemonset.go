package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// InjectVarsDaemonSet inserts EnvVars into a daemonSet doc
func InjectVarsDaemonSet(data []byte, envVars []v1.EnvVar) (*v1beta1.DaemonSet, error) {
	daemonSet := &v1beta1.DaemonSet{}

	if err := json.Unmarshal(data, daemonSet); err != nil {
		return daemonSet, err
	}

	podSpec := injectPodSpecEnvVars(
		daemonSet.Spec.Template.Spec,
		envVars,
	)
	daemonSet.Spec.Template.Spec = podSpec
	return daemonSet, nil
}
