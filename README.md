[![codecov](https://codecov.io/gh/labd/contentful-go/branch/master/graph/badge.svg)](https://codecov.io/gh/labd/contentful-go)
[![Godoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](https://godoc.org/github.com/labd/contentful-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/labd/contentful-go.svg?branch=master)](https://travis-ci.org/labd/contentful-go)


# contentful-go

GoLang SDK for [Contentful's](https://www.contentful.com) Content Delivery, Preview and Management API's.

# About

[Contentful](https://www.contentful.com) provides a content infrastructure for digital teams to power content in 
websites, apps, and devices. Unlike a CMS, Contentful was built to integrate with the modern software stack. 
It offers a central hub for structured content, powerful management and delivery APIs, and a customizable web app 
that enables developers and content creators to ship digital products faster.

[Go](https://golang.org) is an open source programming language that makes it easy to build simple, reliable, and 
efficient software.

# Install

`go get github.com/labd/contentful-go`

# Getting started

Import the SDK into your Go project or library

```go
import (
	contentful "github.com/labd/contentful-go"
)
```

Create an API client in order to interact with the Contentful's API endpoints.

```go
token := "your-cma-token" // Observe your CMA token from Contentful's web page
cma := contentful.NewCMA(token)
```

#### Organization

If your Contentful account is part of an organization, you can setup your API client as such. When you set your 
organization id for the SDK client, every API request will have the `X-Contentful-Organization: <your-organization-id>` 
header automatically.

```go
cma.SetOrganization("your-organization-id")
```

#### Debug mode

When debug mode is activated, the SDK client starts to work in verbose mode and tries to print as much information as 
possible. In debug mode, all outgoing HTTP requests are printed nicely in the form of `curl` commands so that you 
can easily drop into your command line to debug specific requests.

```go
cma.Debug = true
```

#### Dependencies

`contentful-go` stores its dependencies under the `vendor` folder and uses [`dep`](https://github.com/golang/dep) to 
manage dependency resolutions. Dependencies in the `vendor` folder will be loaded automatically by 
[Go 1.6+](https://golang.org/cmd/go/#hdr-Vendor_Directories). To install the dependencies, run `dep ensure`, for more 
options and documentation please visit [`dep`](https://github.com/golang/dep).

# Using the SDK

## Working with resource services

Currently, the SDK exposes the following services:

* Access Tokens
* API Keys
* App Definitions
* App Installations
* Assets
* Content Types
* Editor Interfaces
* Entries
* Environments
* Extensions
* Locales
* Memberships
* Organizations
* Resources/Uploads
* Roles
* Scheduled Actions
* Snapshots
* Spaces
* Usage
* Users
* Webhooks
* Webhook Calls and Health

Every service has multiple interface, for example:

```go
List() *Collection
Get(spaceID, resourceID string) <Resource>, error
Upsert(spaceID string, resourceID *Resource) error
Delete(spaceID string, resourceID *Resource) error
```

To read the interfaces of all services, visit the [Contentful GoDoc](https://godoc.org/github.com/labd/contentful-go).

#### Examples

In the example below, the Get interface of the Spaces service is called. This interface returns an object of the type 
Space. This object could be easily read later by calling the properties of the interface, for example: `Space.Name`

```go
space, err := cma.Spaces.Get("space-id")
if err != nil {
  log.Fatal(err)
}
```

In the following example, the List interface of the Spaces service is called. This interface returns an array of Space 
objects. Working with these collections is explained below.
```go
collection := cma.ContentTypes.List(space.Sys.ID)
collection, err = collection.Next()
if err != nil {
  log.Fatal(err)
}

for _, contentType := range collection.ToContentType() {
  fmt.Println(contentType.Name, contentType.Description)
}
```

## Working with collections

All the endpoints which return an array of objects are wrapped with the `Collection` struct. The main features of the 
`Collection` struct are pagination and type assertion.

### Type assertion

The `Collection` struct exposes the necessary converters (type assertions) such as `ToSpace()`. The following example 
gets all spaces for the given account:

### Example

```go
collection := cma.Spaces.List() // returns a collection
collection, err := collection.Next() // makes the actual api call
if err != nil {
  log.Fatal(err)
}

spaces := collection.ToSpace() // make the type assertion
for _, space := range spaces {
  fmt.Println(space.Name)
  fmt.Println(space.Sys.ID)
}

// In order to access collection metadata
fmt.Println(col.Total)
fmt.Println(col.Skip)
fmt.Println(col.Limit)
```

## Testing

```shell
$> go test
```

To enable higher verbose mode

```shell
$> go test -v
```

## Documentation/References

### Contentful
[Content Delivery API](https://www.contentful.com/developers/docs/references/content-delivery-api/)
[Content Management API](https://www.contentful.com/developers/docs/references/content-management-api/)
[Content Preview API](https://www.contentful.com/developers/docs/references/content-preview-api/)

### GoLang
[Effective Go](https://golang.org/doc/effective_go.html)

## Support

This is a project created for demo purposes and not officially supported, so if you find issues or have questions you can let us know via the [issue](https://github.com/labd/contentful-go/issues/new) page but don't expect a quick and prompt response.

## Contributing

[WIP]

## License

MIT
