//go:build wireinject
// +build wireinject

package controllers

import (
	"github.com/google/wire"
)

var Wired = wire.NewSet(
	New,
)
