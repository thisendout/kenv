package main

import (
	"testing"

	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

func TestPrintResource(t *testing.T) {
	if err := printResource(v1beta1.Deployment{}, false); err != nil {
		t.Fatal(err)
	}

	if err := printResource(v1beta1.Deployment{}, true); err != nil {
		t.Fatal(err)
	}
}
