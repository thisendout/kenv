package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// InjectVarsDeployment inserts EnvVars into a deployment doc
func InjectVarsDeployment(data []byte, envVars []v1.EnvVar) (string, error) {
	deployment := v1beta1.Deployment{}

	if err := json.Unmarshal(data, &deployment); err != nil {
		return "", err
	}

	updateContainersDeployment(&deployment, envVars)

	data, err := json.MarshalIndent(&deployment, "", "  ")

	return string(data), err
}

func updateContainersDeployment(deployment *v1beta1.Deployment, envVars []v1.EnvVar) {
	podSpec := deployment.Spec.Template.Spec
	containers := []v1.Container{}

	for _, c := range podSpec.Containers {
		c.Env = mergeEnvVars(c.Env, envVars)
		containers = append(containers, c)
	}

	podSpec.Containers = containers
	deployment.Spec.Template.Spec = podSpec
}
