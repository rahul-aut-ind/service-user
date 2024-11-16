//go:build wireinject
// +build wireinject

package dynamorepo

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
