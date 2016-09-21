package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/util/yaml"
)

var (
	varsFiles FlagSlice
	flagSet   *flag.FlagSet
)

func init() {
	// workaround to avoid inheriting vendor flags
	flagSet = flag.NewFlagSet("kenv", flag.ExitOnError)
	flagSet.Var(&varsFiles, "v", "File containing environment variables (repeatable)")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  kenv -v fixtures/vars.env fixtures/deployment.yaml\n")
		fmt.Fprintf(os.Stderr, "  cat fixtures/deployment.yaml | kenv -v fixtures/vars.env\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flagSet.PrintDefaults()
	}
}

func main() {
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	// take either a doc as a cli arg or stdin
	var in *os.File
	var err error

	switch name := flagSet.Arg(0); {
	case name == "":
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal(err)
		}
		// Print usage unless we already have STDIN data or incoming pipe
		if fi.Size() == 0 && fi.Mode()&os.ModeNamedPipe == 0 {
			flagSet.Usage()
			return
		}
		in = os.Stdin
	default:
		if in, err = os.Open(name); err != nil {
			log.Fatal(err)
		}
		defer in.Close()
	}

	// read in doc
	data, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal(err)
	}

	// ensure doc is JSON
	data, err = yaml.ToJSON(data)
	if err != nil {
		log.Fatal(err)
	}

	vars := Vars{}

	// read in vars files
	for _, filename := range varsFiles {
		v, err := ReadVarsFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		vars = append(vars, v...)
	}

	envVars := vars.toEnvVar()
	var doc string

	kind, err := getDocKind(data)
	if err != nil {
		log.Fatal(err)
	}

	switch kind {
	case "Deployment":
		doc, err = InjectVarsDeployment(data, envVars)
	case "ReplicationController":
		doc, err = InjectVarsReplicationController(data, envVars)
	case "DaemonSet":
		doc, err = InjectVarsDaemonSet(data, envVars)
	case "ReplicaSet":
		doc, err = InjectVarsReplicaSet(data, envVars)
	default:
		err = fmt.Errorf("Kind %s not supported\n", kind)
	}

	if err != nil {
		log.Fatal(err)
	}

	println(doc)
}

// FlagSlice represents a repeatable string flag
type FlagSlice []string

// String returns a string representation of FlagSlice
func (f *FlagSlice) String() string {
	return strings.Join(*f, ",")
}

// Set appends a string value to FlagSlice
func (f *FlagSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func getDocKind(data []byte) (string, error) {
	typeMeta := unversioned.TypeMeta{}
	if err := json.Unmarshal(data, &typeMeta); err != nil {
		return "", err
	}
	return typeMeta.Kind, nil
}
