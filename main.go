package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/util/yaml"
)

var (
	varsFiles      FlagSlice
	secretFiles    FlagSlice
	configMapFiles FlagSlice
	name           string
	namespace      string
	convertKeys    bool
	toYAML         bool
	flagSet        *flag.FlagSet
)

func init() {
	// workaround to avoid inheriting vendor flags
	flagSet = flag.NewFlagSet("kenv", flag.ExitOnError)
	flagSet.StringVar(&name, "name", "", "Name to give the ConfigMap and Secret resources")
	flagSet.StringVar(&namespace, "namespace", "default", "Namespace to create the ConfigMap in")
	flagSet.BoolVar(&convertKeys, "convert-keys", false, "Convert ConfigMap keys to support k8s version < 1.4")
	flagSet.BoolVar(&toYAML, "yaml", false, "Output as YAML")
	flagSet.Var(&varsFiles, "v", "Files containing variables to inject as environment variables (repeatable)")
	flagSet.Var(&secretFiles, "s", "Files containing variables to inject as Secrets (repeatable)")
	flagSet.Var(&configMapFiles, "c", "Files containing variables to inject as ConfigMaps (repeatable)")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, `
Examples:

  kenv -v fixtures/vars.env fixtures/deployment.yaml
	kenv -v fixtures/vars.env -s fixtures/vars.yaml fixtures/deployment.yaml
  cat fixtures/deployment.yaml | kenv -v fixtures/vars.env

Options:
`)
		flagSet.PrintDefaults()
	}
}

func main() {
	var in *os.File
	var err error

	if err = flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

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

	envVars := []v1.EnvVar{}

	if len(varsFiles) > 0 {
		v, err := newVarsFromFiles(varsFiles)
		if err != nil {
			log.Fatal(err)
		}

		e := v.toEnvVar()
		if err != nil {
			log.Fatal(err)
		}

		envVars = append(envVars, e...)
	}

	if len(secretFiles) > 0 {
		if name == "" {
			log.Fatal("A name must be set for the Secret resource")
		}

		v, err := newVarsFromFiles(secretFiles)
		if err != nil {
			log.Fatal(err)
		}

		e, secret, err := v.toSecret(name, namespace, convertKeys)
		if err != nil {
			log.Fatal(err)
		}

		if err = printResource(secret, toYAML); err != nil {
			log.Fatal(err)
		}

		envVars = append(envVars, e...)
	}

	if len(configMapFiles) > 0 {
		if name == "" {
			log.Fatal("A name must be set for the ConfigMap resource")
		}

		v, err := newVarsFromFiles(configMapFiles)
		if err != nil {
			log.Fatal(err)
		}

		e, configMap, err := v.toConfigMap(name, namespace, convertKeys)
		if err != nil {
			log.Fatal(err)
		}

		if err = printResource(configMap, toYAML); err != nil {
			log.Fatal(err)
		}

		envVars = append(envVars, e...)
	}

	// inject environment variables into the supplied resource doc
	// and print the result to STDOUT
	object, err := InjectVars(data, envVars)
	if err = printResource(object, toYAML); err != nil {
		log.Fatal(err)
	}
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
