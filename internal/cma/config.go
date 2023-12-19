package cma

import "github.com/flaconi/contentful-go/service/common"

type ClientConfig struct {
	URL        string
	HTTPClient common.HttpClient
	Debug      bool
	UserAgent  string
	Token      string
}
