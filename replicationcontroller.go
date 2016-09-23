package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/v1"
)

// InjectVarsReplicationController inserts EnvVars into a replicationController doc
func InjectVarsReplicationController(data []byte, envVars []v1.EnvVar) (*v1.ReplicationController, error) {
	replicationController := &v1.ReplicationController{}
	if err := json.Unmarshal(data, replicationController); err != nil {
		return replicationController, err
	}

	podSpec := injectPodSpecEnvVars(
		replicationController.Spec.Template.Spec,
		envVars,
	)

	replicationController.Spec.Template.Spec = podSpec
	return replicationController, nil
}
