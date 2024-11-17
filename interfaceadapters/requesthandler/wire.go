//go:build wireinject
// +build wireinject

package requesthandler

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
