package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// InjectVarsReplicaSet inserts EnvVars into a replicaSet doc
func InjectVarsReplicaSet(data []byte, envVars []v1.EnvVar) (string, error) {
	replicaSet := v1beta1.ReplicaSet{}

	if err := json.Unmarshal(data, &replicaSet); err != nil {
		return "", err
	}

	updateContainersReplicaSet(&replicaSet, envVars)

	data, err := json.MarshalIndent(&replicaSet, "", "  ")

	return string(data), err
}

func updateContainersReplicaSet(replicaSet *v1beta1.ReplicaSet, envVars []v1.EnvVar) {
	podSpec := replicaSet.Spec.Template.Spec
	containers := []v1.Container{}

	for _, c := range podSpec.Containers {
		c.Env = mergeEnvVars(c.Env, envVars...)
		containers = append(containers, c)
	}

	podSpec.Containers = containers
	replicaSet.Spec.Template.Spec = podSpec
}
