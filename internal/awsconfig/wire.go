//go:build wireinject
// +build wireinject

package awsconfig

import "github.com/google/wire"

var Wired = wire.NewSet(
	NewAWSConfig,
)
