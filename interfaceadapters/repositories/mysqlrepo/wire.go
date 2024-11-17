//go:build wireinject
// +build wireinject

package mysqlrepo

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
