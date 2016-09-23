package main

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

// InjectVarsDeployment inserts EnvVars into a deployment doc
func InjectVarsDeployment(data []byte, envVars []v1.EnvVar) (*v1beta1.Deployment, error) {
	deployment := &v1beta1.Deployment{}
	if err := json.Unmarshal(data, deployment); err != nil {
		return deployment, err
	}

	podSpec := injectPodSpecEnvVars(
		deployment.Spec.Template.Spec,
		envVars,
	)

	deployment.Spec.Template.Spec = podSpec
	return deployment, nil
}
