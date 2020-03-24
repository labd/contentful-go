package contentful

import (
	"fmt"
	"net/http"
	"net/url"
)

// SnapshotsService service
type SnapshotsService service

// EntrySnapshot model
type EntrySnapshot struct {
	Sys                 *Sys                `json:"sys"`
	EntrySnapshotDetail EntrySnapshotDetail `json:"snapshot"`
}

// EntrySnapshotDetail model
type EntrySnapshotDetail struct {
	Fields map[string]interface{} `json:"fields"`
	Sys    *Sys                   `json:"sys"`
}

// ContentTypeSnapshot model
type ContentTypeSnapshot struct {
	Sys                       *Sys                      `json:"sys"`
	ContentTypeSnapshotDetail ContentTypeSnapshotDetail `json:"snapshot"`
}

// ContentTypeSnapshotDetail model
type ContentTypeSnapshotDetail struct {
	Name   string              `json:"name"`
	Fields []ContentTypeFields `json:"fields"`
	Sys    *Sys                `json:"sys"`
}

// ContentTypeFields model
type ContentTypeFields struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Required  bool   `json:"required"`
	Localized bool   `json:"localized"`
	Type      string `json:"type"`
}

// ListEntrySnapshots returns snapshot collection
func (service *SnapshotsService) ListEntrySnapshots(spaceID, entryID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments/%s/entries/%s/snapshots", spaceID, service.c.Environment, entryID)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// GetEntrySnapshot returns a single snapshot of an entry
func (service *EntriesService) GetEntrySnapshot(spaceID, entryID, snapshotID string) (*EntrySnapshot, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/entries/%s/snapshots/%s", spaceID, service.c.Environment, entryID, snapshotID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &EntrySnapshot{}, err
	}

	var entrySnapshot EntrySnapshot
	if ok := service.c.do(req, &entrySnapshot); ok != nil {
		return nil, err
	}

	return &entrySnapshot, err
}

// ListContentTypeSnapshots returns snapshot collection
func (service *SnapshotsService) ListContentTypeSnapshots(spaceID, contentTypeID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s/snapshots", spaceID, service.c.Environment, contentTypeID)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// GetContentTypeSnapshots returns a single snapshot of an entry
func (service *EntriesService) GetContentTypeSnapshots(spaceID, contentTypeID, snapshotID string) (*ContentTypeSnapshot, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s/snapshots/%s", spaceID, service.c.Environment, contentTypeID, snapshotID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &ContentTypeSnapshot{}, err
	}

	var contentTypeSnapshot ContentTypeSnapshot
	if ok := service.c.do(req, &contentTypeSnapshot); ok != nil {
		return nil, err
	}

	return &contentTypeSnapshot, err
}
