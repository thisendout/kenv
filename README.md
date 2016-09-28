# kenv

[![Build Status](https://travis-ci.org/thisendout/kenv.svg?branch=master)](https://travis-ci.org/thisendout/kenv)

kenv is an environment file injector for Kubernetes resources.

kenv injects variables into Kubernetes resource documents by loading a list of files containing variables and modifying the resource document to include those values. This way, you can dynamically set environment variables without having to template your resource documents.

![kenv Example](example.gif)

kenv supports referencing environment variables as:

* [Plaintext K/V Pairs](http://kubernetes.io/docs/user-guide/configuring-containers/#environment-variables-and-variable-expansion)
* [ConfigMaps](http://kubernetes.io/docs/user-guide/configmap/)
* [Secrets](http://kubernetes.io/docs/user-guide/secrets/)

## Getting Started

[Download](https://github.com/thisendout/kenv/releases/tag/v0.3.0) kenv and run by passing a resource doc (YAML or JSON) to:

STDIN:

```
cat fixtures/deployment.yaml | ./kenv -v fixtures/vars.env
./kenv -v fixtures/vars.env < fixtures/deployment.yaml
```

or as a CLI Arg:

```
./kenv -v fixtures/vars.env fixtures/deployment.yaml
```

## Usage

```
Usage: kenv [options] file

Examples:

  kenv -v fixtures/vars.env fixtures/deployment.yaml
  kenv -name nginx -v fixtures/vars.env -s fixtures/secrets.yml fixtures/deployment.yaml
  cat fixtures/deployment.yaml | kenv -v fixtures/vars.env

Options:
  -c value
    	Files containing variables to inject as ConfigMaps (repeatable)
  -convert-keys
    	Convert ConfigMap keys to support k8s version < 1.4
  -name string
    	Name to give the ConfigMap and Secret resources
  -namespace string
    	Namespace to create the ConfigMap in (default "default")
  -s value
    	Files containing variables to inject as Secrets (repeatable)
  -v value
    	Files containing variables to inject as environment variables (repeatable)
  -yaml
    	Output as YAML
```

## Variables

Variables are stored in files either as simple `key=value` pairs, or YAML pairs (`key: value`). Multiple variable files can be specified by repeating the `-v` (inject directly as plaintext variables), `-c` (create and map as ConfigMaps), and `-s` (create and map as Secrets) flags. These flags are not mutually exclusive, meaning you can specify plaintext, ConfigMaps, and Secrets within the same command:

```
./kenv \
  -v fixtures/vars.env \
  -c fixtures/configmap.env \
  -s fixtures/secrets.yml \
  fixtures/deployment.yaml
```

### YAML Format

YAML files must be in the following format (nested data types are not currently supported):

```
key1: value1
key2: value2
```

### KV Format

A `key=value` format is also supported, for example:

```
key1=value1
key2=value2
```

### Injection

Variables are injected into the resource doc specified by the user as either plaintext environment variables, [ConfigMaps](http://kubernetes.io/docs/user-guide/configmap/), or [Secrets](http://kubernetes.io/docs/user-guide/secrets/). When specifying ConfigMaps and/or Secrets, you must also set a `-name` for the ConfigMap/Secret resource being created.

kenv injects the variables into the PodSpec for the following resources:

 * `DaemonSet`
 * `Deployment`
 * `ReplicaSet`
 * `ReplicationController`

### Conversion and Support for K8S < 1.4

When using ConfigMap and/or Secret resources in Kubernetes version < 1.4, keys must adhere to the following regex:

```
\.?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*
```

As a convenience, you can pass `-convert-keys`, which will replace underscores with dashes and convert uppercase strings to lower.

## Examples

### Plaintext K/V Injection

Variables can be injected as plaintext, directly tied to the containers in the resource document by using the `-v` flag.

Given the variable file below (fixtures/plaintext.env):

```
ptkey1=ptvalue1
ptkey2=ptvalue2
```

Running the following command injecting as plaintext would result in:

```
$ kenv -yaml -v fixtures/plaintext.env fixtures/deployment.yaml
---
apiVersion: extensions/v1beta1
kind: Deployment
... snip ...
    spec:
      containers:
      - env:
        - name: ptkey1
          value: ptvalue1
        - name: pykey2
          value: ptvalue2
        image: nginx:latest
... snip ...
```

### ConfigMap Injection

Variables can be injected as ConfigMaps references by specifying variable files with the flag `-c`. This will both print a ConfigMap resource to be created and modify the specified resource to use the ConfigMap.

Given the variable file below (fixtures/configmap.env):

```
cmkey1=cmvalue1
cmkey2=cmvalue2
```

Running the following command injecting as a ConfigMap would result in:

```
$ kenv -yaml -name nginx -c fixtures/plaintext.env fixtures/deployment.yaml
---
apiVersion: v1
data:
  cmkey1: cmvalue1
  cmkey2: cmvalue2
kind: ConfigMap
metadata:
  name: nginx
  namespace: default
---
apiVersion: extensions/v1beta1
kind: Deployment
... snip ...
    spec:
      containers:
      - env:
        - name: cmkey1
          valueFrom:
            configMapKeyRef:
              key: cmkey1
              name: nginx
        - name: cmkey2
          valueFrom:
            configMapKeyRef:
              key: cmkey2
              name: nginx
        image: nginx:latest
... snip ...
```

### Secret Injection

Variables can be injected as Secret references by specifying variable files with the flag `-s`. This will both print a Secret resource to be created and modify the specified resource to use the Secret.

Given the variable file below (fixtures/secrets.yml):

```
secretkey1: secretvalue1
secretkey2: secretvalue2
```

Running the following command injecting as a ConfigMap would result in:

```
$ kenv -yaml -name nginx -s fixtures/secrets.yml fixtures/deployment.yaml
---
apiVersion: v1
data:
  secretkey1: YzJWamNtVjBkbUZzZFdVeA==
  secretkey2: YzJWamNtVjBkbUZzZFdVeQ==
kind: Secret
metadata:
  creationTimestamp: null
  name: nginx
  namespace: default
---
apiVersion: extensions/v1beta1
kind: Deployment
... snip ...
    spec:
      containers:
      - env:
				- name: secretkey1
          valueFrom:
            secretKeyRef:
              key: secretkey1
              name: nginx
        - name: secretkey2
          valueFrom:
            secretKeyRef:
              key: secretkey2
              name: nginx
        image: nginx:latest
... snip ...
```

### Plaintext, ConfigMaps, and Secrets

Combining the examples above into one command, you would get the following output:

```
$ kenv -yaml -name nginx -v fixtures/plaintext.env -c fixtures/configmap.env -s fixtures/secrets.yml fixtures/deployment.yaml
---
apiVersion: v1
data:
  secretkey1: YzJWamNtVjBkbUZzZFdVeA==
  secretkey2: YzJWamNtVjBkbUZzZFdVeQ==
kind: Secret
metadata:
  name: nginx
  namespace: default
---
apiVersion: v1
data:
  cmkey1: cmvalue1
  cmkey2: cmvalue2
kind: ConfigMap
metadata:
  name: nginx
  namespace: default
---
apiVersion: extensions/v1beta1
kind: Deployment
... snip ...
    spec:
      containers:
      - env:
        - name: ptkey1
          value: ptvalue1
        - name: pykey2
          value: ptvalue2
        - name: secretkey1
          valueFrom:
            secretKeyRef:
              key: secretkey1
              name: nginx
        - name: secretkey2
          valueFrom:
            secretKeyRef:
              key: secretkey2
              name: nginx
        - name: cmkey1
          valueFrom:
            configMapKeyRef:
              key: cmkey1
              name: nginx
        - name: cmkey2
          valueFrom:
            configMapKeyRef:
              key: cmkey2
              name: nginx
        image: nginx:latest
... snip ...
```

## Building

```
govendor sync
go test -v .
```
