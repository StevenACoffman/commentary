// +build tools

package tools

import (
	_ "github.com/Khan/genqlient"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/segmentio/golines"
	_ "mvdan.cc/gofumpt"
)

// There is nothing here intentionally
