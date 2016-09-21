package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/v1"
)

// InjectVarsReplicationController inserts EnvVars into a replicationController doc
func InjectVarsReplicationController(data []byte, envVars []v1.EnvVar) (string, error) {
	replicationController := v1.ReplicationController{}

	if err := json.Unmarshal(data, &replicationController); err != nil {
		return "", err
	}

	updateContainersReplicationController(&replicationController, envVars)

	data, err := json.MarshalIndent(&replicationController, "", "  ")

	return string(data), err
}

func updateContainersReplicationController(replicationController *v1.ReplicationController, envVars []v1.EnvVar) {
	podSpec := replicationController.Spec.Template.Spec
	containers := []v1.Container{}

	for _, c := range podSpec.Containers {
		c.Env = append(c.Env, envVars...)
		containers = append(containers, c)
	}

	podSpec.Containers = containers
	replicationController.Spec.Template.Spec = podSpec
}
