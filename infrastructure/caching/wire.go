//go:build wireinject
// +build wireinject

package caching

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
