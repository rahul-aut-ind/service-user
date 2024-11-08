//go:build wireinject
// +build wireinject

package middlewares

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
