package cma

import (
	"github.com/flaconi/contentful-go/pkgs/common"
)

type NextableCollection[Items any, Includes any] interface {
	Next() (*common.Collection[Items, Includes], error)
	GetQuery() *common.Query
}

type SyncCollection interface {
	Next() (*common.InterfaceCollection, error)
	GetQuery() *common.Query
}
