v0.5.0 (2020-10-22)
===================

New resources
-------------
* `+` uploads entity
* `+` app-definitions entity
* `+` app-installations entity
* `+` getting usage statistics of an organization
* `+` webhook entity
* `+` extensions entity
* `+` editor interfaces entity
* `+` scheduledactions entity
* `+` entrytasks entity
* `+` accesstokens entity
* `+` snapshots entity
* `+` assets entity
* `+` memberships entity
* `+` roles entity
* `+` entries entity
* `+` environments entity
* `+` environment-alias entity
* `+` getting organization of authenticated user
* `+` getting the authenticated user
* `+` Webhook calls and webhook health support

Other improvements
------------------
* `~` Added missing unit tests
* `~` Improved Coverage
* `~` Code cleanup
* `~` More asserts in unit tests, test more specific per entity


v0.4.0 (2019-04-29)
===================
* `+` Added support for creating/uploading resources.
* `+` Adding test for contentful.SetClient
* `~` Allowing the HTTP client to be set by the consumer. Useful for testing or implementing more robust HTTP clients.
* `+` Add Go Modules
* `+` added field type Symbol
* `+` Created an asset Alias for custom UnmarshalJSON method and fixed typo
* `+` Add MIT License
* `~` Capture validation errors
* `+` New content types can specify their ID.
* `~` Rename methods Activate/Deactivate to Publish/Unpublish
* `+` Extend EntryService by Delete, Publish and Unpublish actions


v0.3.1 (2017-11-28)
===================
* `~` sdk version header format fixed.
  

v0.3.0 (2017-11-11)
===================
* `~` use codecov as coverage service
* `+` `golint` is added to the CI process
* `~` `dep` is updated to the latest version
* `x` `vendor` folder is not under version control
* `~` `makefile` simplifications
* `+` testing and linting is now handled by scripts under `tools` folder
* `~` cosmetic changes in codebase to make linter happy


v0.2.0 (2017-04-12)
===================
* `~` Godoc style examples
* `+` Query.go tests
* `+` Locale resource tests
* `+` ContentType tests
* `+` Missing space resource tests
* `+` User-Agent for api requests


v0.1.1 (2017-03-31)
===================

* `+` Rate-limited api requests
* `~` Locale model
* `+` Content type field unmarshaling

v0.1.0 (2017-03-26)
===================

### Introducing resource services
Every entity resource now has its own service definition for handling api communication. With this release, we don't store `contentful client` and `space` objects inside entities anymore. Resource services now get `spaceID` as a string parameter when it is neccessary.

With the old versions, in order to create a new `ContentType`, for example, you first need to observe `Space` object. That is no longer required. The problem with the old method was that you had to make an extra api request to observe the `Space` in order to interact with rest of the resources. The following example demonstras the difference between old and new version.

```go
// prior to v0.1.0
space, err := cma.GetSpace("space-id") // this call was making an extra api call
contentTypes, err := space.ContentTypes()

// after v0.1.0
contentType := &contentful.ContentType{ ... }
spaceID := "space-id"
cma.ContentTypes.Upsert(spaceID, contentType) // now we are passing spaceID as string
```

You can access available resources as follows:

```go
cma := contentful.NewCMA("token")
cma.Spaces
cma.APIKeys
cma.Assets
cma.ContentTypes
cma.Entries
cma.Locales
cma.Webhooks
```

Every resource service exposes the following interface:

`List(spaceID string) *Collection`

`Get(spaceID, contentTypeID string)(*ContentType, error)`

`Upsert(spaceID string, ct *ContentType) error`

`Delete(spaceID string, ct *ContentType) error`

### Create resource instancas directly from their model definitions

All `New{ResourceName}`, such as `NewContentType`, `NewSpace`, functions are removed from the SDK. It turned out that it wasn't a good practice in golang. Instead, you can directly initiate resource entities directly from their models, such as:

```go
contentType := &contentful.ContentType{
    Name: "name",
    ... other fields
}
```


v0.0.3 (2017-03-22)
===================
* `+` PredefinedValues validation
* `+` Range validation greater/less than equal to support.
* `+` Size validation for content type field.
* `+` Packages are vendored with `godep`.
* `+` `version.go`.
* `+` `entity/content_type`: regex validation for content type field.
* `+` Validation data structures added: `MinMax`, `Regex`
* `+` `LinkType` support for `Field` struct
* `+` New validations: `MimeType`, `Dimension`, `FileSize`


v0.0.2 (2017-03-21)
===================
* `~` `entity/webhook`: add tests for webhook entity.
* `~` `entity/space`: add tests for space entity.
* `~` `errors`: add tests for error handler.
* `~` `entity/content_type`: add test for content type entity.
* `~` `entity/content_type`: Field validations added for link type
* `~` `entity/content_type`: field validations added: Range, PredefinedValues, Unique


v0.0.1 (2017-03-20)
===================
* `~` `sdk`: first implementation.
* `~` `collection`: first implementation.
* `~` `entity/content_type`: first implementation.
* `~` `entity/entry`: first implementation.
* `~` `entity/query`: first implementation.
* `~` `entity/asset`: first implementation.
* `~` `entity/locale`: first implementation.
* `~` `entity/space`: first implementation.
* `~` `entity/webhook`: first implementation.
* `~` `entity/api_key`: first implementation.
* `~` `sdk`: basic documentation.
* `~` `examples`: some examples for entities
