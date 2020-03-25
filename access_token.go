package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// AccessTokensService service
type AccessTokensService service

// AccessToken model
type AccessToken struct {
	Sys       *Sys     `json:"sys,omitempty"`
	Name      string   `json:"name,omitempty"`
	RevokedAt string   `json:"description,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
}

// GetVersion returns entity version
func (accessToken *AccessToken) GetVersion() int {
	version := 1
	if accessToken.Sys != nil {
		version = accessToken.Sys.Version
	}

	return version
}

// List returns an access tokens collection
func (service *AccessTokensService) List() *Collection {
	path := fmt.Sprint("/users/me/access_tokens")
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single access token
func (service *AccessTokensService) Get(accessTokenID string) (*AccessToken, error) {
	path := fmt.Sprintf("/users/me/access_tokens/%s", accessTokenID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &AccessToken{}, err
	}

	var accessToken AccessToken
	if ok := service.c.do(req, &accessToken); ok != nil {
		return nil, err
	}

	return &accessToken, err
}
