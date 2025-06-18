//go:build tools

// This file implements that pattern:
// https://go.dev/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// for etcd. Thanks to this file 'go mod tidy' does not removes dependencies.

package mod

import (
	_ "github.com/elastic/crd-ref-docs"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "sigs.k8s.io/controller-runtime/tools/setup-envtest"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
	_ "sigs.k8s.io/kind"
	_ "sigs.k8s.io/kustomize/kustomize/v5"
)