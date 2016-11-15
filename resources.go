package main

import (
	"encoding/json"
	"io"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/yaml"
)

// KubeResource represents a resource kind and raw data to be used
// later for injecting EnvVars
type KubeResource struct {
	Kind string
	Data []byte
}

// ParseDocs iterates through YAML or JSON docs and discovers
// their type returning a list of KubeResources
func ParseDocs(reader io.Reader) ([]KubeResource, error) {
	resources := []KubeResource{}
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 4096)

	for {
		rawExtension := runtime.RawExtension{}
		err := decoder.Decode(&rawExtension)
		if err == io.EOF {
			break
		} else if err != nil {
			return resources, err
		}

		kind, err := getResourceKind(rawExtension.Raw)
		if err != nil {
			return resources, err
		}

		resources = append(resources, KubeResource{
			Kind: kind,
			Data: rawExtension.Raw,
		})
	}

	return resources, nil
}

// InjectVarsDeployment inserts EnvVars into a Deployment doc
func (k *KubeResource) InjectVarsDeployment(envVars []v1.EnvVar) (*v1beta1.Deployment, error) {
	deployment := &v1beta1.Deployment{}

	if err := json.Unmarshal(k.Data, &deployment); err != nil {
		return deployment, err
	}

	podSpec := injectPodSpecEnvVars(
		deployment.Spec.Template.Spec,
		envVars,
	)
	deployment.Spec.Template.Spec = podSpec
	return deployment, nil
}

// InjectVarsDaemonSet inserts EnvVars into a daemonSet doc
func (k *KubeResource) InjectVarsDaemonSet(envVars []v1.EnvVar) (*v1beta1.DaemonSet, error) {
	daemonSet := &v1beta1.DaemonSet{}

	if err := json.Unmarshal(k.Data, daemonSet); err != nil {
		return daemonSet, err
	}

	podSpec := injectPodSpecEnvVars(
		daemonSet.Spec.Template.Spec,
		envVars,
	)
	daemonSet.Spec.Template.Spec = podSpec
	return daemonSet, nil
}

// InjectVarsReplicaSet inserts EnvVars into a replicaSet doc
func (k *KubeResource) InjectVarsReplicaSet(envVars []v1.EnvVar) (*v1beta1.ReplicaSet, error) {
	replicaSet := &v1beta1.ReplicaSet{}
	if err := json.Unmarshal(k.Data, replicaSet); err != nil {
		return replicaSet, err
	}

	podSpec := injectPodSpecEnvVars(
		replicaSet.Spec.Template.Spec,
		envVars,
	)

	replicaSet.Spec.Template.Spec = podSpec
	return replicaSet, nil
}

// InjectVarsRC inserts EnvVars into a replicationController doc
func (k *KubeResource) InjectVarsRC(envVars []v1.EnvVar) (*v1.ReplicationController, error) {
	replicationController := &v1.ReplicationController{}
	if err := json.Unmarshal(k.Data, replicationController); err != nil {
		return replicationController, err
	}

	podSpec := injectPodSpecEnvVars(
		replicationController.Spec.Template.Spec,
		envVars,
	)

	replicationController.Spec.Template.Spec = podSpec
	return replicationController, nil
}

// UnmarshalGeneric does not attempt to unmarshal to a known type,
// instead returns a generic interface object for displaying to the user
func (k *KubeResource) UnmarshalGeneric() (interface{}, error) {
	var generic interface{}
	if err := json.Unmarshal(k.Data, &generic); err != nil {
		return generic, err
	}

	return generic, nil
}

// getResourceKind unmarshalls a file and returns the kind of resource doc
func getResourceKind(data []byte) (string, error) {
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
