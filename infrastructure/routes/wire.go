//go:build wireinject
// +build wireinject

package routes

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
