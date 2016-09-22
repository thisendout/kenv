package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// InjectVarsReplicaSet inserts EnvVars into a replicaSet doc
func InjectVarsReplicaSet(data []byte, envVars []v1.EnvVar) (*v1beta1.ReplicaSet, error) {
	replicaSet := &v1beta1.ReplicaSet{}
	if err := json.Unmarshal(data, replicaSet); err != nil {
		return replicaSet, err
	}

	podSpec := injectPodSpecEnvVars(
		replicaSet.Spec.Template.Spec,
		envVars,
	)

	replicaSet.Spec.Template.Spec = podSpec
	return replicaSet, nil
}
