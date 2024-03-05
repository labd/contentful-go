package cda

import (
	"github.com/flaconi/contentful-go/service/common"
)

type SpaceIdClientBuilder interface {
	common.RestClient
	WithSpaceId(spaceId string) SpaceIdClient
}

type SpaceIdClient interface {
	common.RestClient
	WithEnvironment(environment string) EnvironmentClient
}

type EnvironmentClient interface {
	common.EnvironmentClient
	Sync() Sync
}
