package cma

import "github.com/flaconi/contentful-go/pkgs/common"

type NextableCollection[Items any, Includes any] interface {
	Next() (*common.Collection[Items, Includes], error)
}
