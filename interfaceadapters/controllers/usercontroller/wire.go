//go:build wireinject
// +build wireinject

package usercontroller

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
