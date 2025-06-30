//go:build tools
// +build tools

// Package tools imports various tools used in the development process.
// This file is used to track tool dependencies using Go modules.
// These imports ensure that `go mod tidy` does not remove these dependencies.
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/securego/gosec/v2/cmd/gosec"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/vuln/cmd/govulncheck"
)