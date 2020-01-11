// +build tools

package tools

import (
	"github.com/alecthomas/gometalinter"
	"github.com/client9/misspell/cmd/misspell" // lint
	"github.com/gordonklaus/ineffassign"       // lint
	"github.com/tsenart/deadcode"              // lint
	"golang.org/x/lint/golint"                 // lint
)
