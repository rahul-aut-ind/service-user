//go:build wireinject
// +build wireinject

package imageservice

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
