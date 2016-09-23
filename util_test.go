package main

import (
	"testing"

	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
)

func TestPrintJSON(t *testing.T) {
	if err := printJSON(v1beta1.Deployment{}); err != nil {
		t.Fatal(err)
	}
}
