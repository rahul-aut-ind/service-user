//go:build wireinject
// +build wireinject

package s3repo

import "github.com/google/wire"

var Wired = wire.NewSet(
	New,
)
