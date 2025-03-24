package cda

import (
	"context"

	"github.com/labd/contentful-go/service/cma"
)

type SyncType int

const (
	Asset SyncType = iota
	Entry
	All
	OnlyDeletion
	DeletedAsset
	ADeletedEntry
)

func (s SyncType) String() string {
	return [...]string{"Asset", "Entry", "all", "Deletion", "DeletedAsset", "DeletedEntry"}[s]
}

type Sync interface {
	Init(ctx context.Context, syncType SyncType, contentType *string) cma.SyncCollection
	GetFromSyncUrl(ctx context.Context, syncUrl string) cma.SyncCollection
}
