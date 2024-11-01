//go:build wireinject
// +build wireinject

package logger

import "github.com/google/wire"

var Wired = wire.NewSet(
	New,
)
