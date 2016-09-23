package main

import (
	"encoding/json"
	"fmt"
	"log"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/runtime"
)

func InjectVars(data []byte, envVars []v1.EnvVar) (runtime.Object, error) {
	kind, err := getDocKind(data)
	if err != nil {
		log.Fatal(err)
	}

	switch kind {
	case "Deployment":
		return injectVarsDeployment(data, envVars)
	case "DaemonSet":
		return injectVarsDaemonSet(data, envVars)
	case "ReplicaSet":
		return injectVarsReplicaSet(data, envVars)
	case "ReplicationController":
		return injectVarsReplicationController(data, envVars)
	default:
		return &v1beta1.Deployment{}, fmt.Errorf("Kind %s not supported\n", kind)
	}
}

// injectVarsDeployment inserts EnvVars into a deployment doc
func injectVarsDeployment(data []byte, envVars []v1.EnvVar) (*v1beta1.Deployment, error) {
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

// injectVarsDaemonSet inserts EnvVars into a daemonSet doc
func injectVarsDaemonSet(data []byte, envVars []v1.EnvVar) (*v1beta1.DaemonSet, error) {
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

// injectVarsReplicaSet inserts EnvVars into a replicaSet doc
func injectVarsReplicaSet(data []byte, envVars []v1.EnvVar) (*v1beta1.ReplicaSet, error) {
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

// injectVarsReplicationController inserts EnvVars into a replicationController doc
func injectVarsReplicationController(data []byte, envVars []v1.EnvVar) (*v1.ReplicationController, error) {
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

// getDocKind unmarshalls a file and returns the kind of resource doc
func getDocKind(data []byte) (string, error) {
	typeMeta := unversioned.TypeMeta{}
	if err := json.Unmarshal(data, &typeMeta); err != nil {
		return "", err
	}
	return typeMeta.Kind, nil
}

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
