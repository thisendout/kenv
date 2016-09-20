# kenv

[![Build Status](https://travis-ci.org/thisendout/kenv.svg?branch=master)](https://travis-ci.org/thisendout/kenv)

Environment file preprocessor for Kubernetes deployments

kenv injects variables into Kubernetes resource documents by loading a list of files containing variables, and modifying the resource document to include those values. This way, you can dynamically set environment variables without having to template your resource documents.

## Getting Started

[Download](https://github.com/thisendout/kenv/releases/tag/v0.1.0) kenv and run by passing a resource doc (YAML or JSON) to:

STDIN with a pipe:

```
cat fixtures/deployment.yaml | ./kenv -v fixtures/vars.env
```

or as a CLI Arg:

```
./kenv -v fixtures/vars.env fixtures/deployment.yaml
```

## Variables

Variables are stored in files either as simple `key=value` pairs, or YAML pairs (`key: value`). Multiple variable files can be specified by repeating the `-v` flag:

```
./kenv -v fixtures/vars.env -v fixtures/complex.env fixtures/deployment.yaml
```

### YAML

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

## Building

```
govendor sync
go test -v .
```

## TODO

* Support [ConfigMaps](http://kubernetes.io/docs/user-guide/configmap/)
* Support [Secrets](http://kubernetes.io/docs/user-guide/secrets/walkthrough/)
* Support [DaemonSets](http://kubernetes.io/docs/admin/daemons/)
* Support [ReplicationControllers](http://kubernetes.io/docs/user-guide/replication-controller/)
* Support [ReplicaSets](http://kubernetes.io/docs/user-guide/replicasets/)
