package cma

import (
	"github.com/flaconi/contentful-go/service/common"
)

type SpaceIdClientBuilder interface {
	common.RestClient
	WithSpaceId(spaceId string) SpaceIdClient
	WithOrganizationId(organizationId string) OrganizationIdClient
}

type SpaceIdClient interface {
	common.RestClient
	WithEnvironment(environment string) EnvironmentClient
	ApiKeys() ApiKeys
	PreviewApiKeys() PreviewApiKeys
	EnvironmentAliases() EnvironmentAliases
	Environments() Environments
}

type EnvironmentClient interface {
	common.EnvironmentClient
	AppInstallations() AppInstallations
	Entries() Entries
	Assets() Assets
	ContentTypes() ContentTypes
	Locales() Locales
}

type OrganizationIdClient interface {
	common.RestClient
	AppDefinitions() AppDefinitions
}
