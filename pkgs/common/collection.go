package common

import (
	"bytes"
	"encoding/json"

	"github.com/labd/contentful-go/pkgs/model"
)

// CollectionOptions holds init options
type CollectionOptions struct {
	Limit uint16
}

type BaseCollection struct {
	*Query
	Sys         *model.BaseSys `json:"sys"`
	Total       int            `json:"total"`
	Skip        int            `json:"skip"`
	Limit       int            `json:"limit"`
	NextPageUrl string         `json:"nextPageUrl,omitempty"`
	NextSyncUrl string         `json:"NextSyncUrl,omitempty"`
}

type Collection[Items any, Includes any] struct {
	BaseCollection
	Items    []Items  `json:"items"`
	Includes Includes `json:"includes"`
}

type InterfaceCollection struct {
	BaseCollection
	Items    []any `json:"items"`
	Includes any   `json:"includes"`
}

// ToAsset cast Items to Asset model
func (col *InterfaceCollection) ToAsset() []*model.Asset {
	var assets []*model.Asset

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&assets)

	return assets
}

// ToEntry cast Items to Entry model
func (col *InterfaceCollection) ToEntry() []*model.Entry {
	var entries []*model.Entry

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&entries)

	return entries
}

// NewCollection initializes a new collection
func NewCollection[Items any, Includes any](options *CollectionOptions) *Collection[Items, Includes] {
	query := NewQuery()
	query.Order("sys.createdAt", true)

	if options.Limit > 0 {
		query.Limit(options.Limit)
	}

	return &Collection[Items, Includes]{
		BaseCollection: BaseCollection{
			Query: query,
		},
	}
}

// NewInterfaceCollection initializes a new collection
func NewInterfaceCollection(options *CollectionOptions) *InterfaceCollection {
	query := NewQuery()
	query.Order("sys.createdAt", true)

	if options.Limit > 0 {
		query.Limit(options.Limit)
	}

	return &InterfaceCollection{
		BaseCollection: BaseCollection{
			Query: query,
		},
	}
}
