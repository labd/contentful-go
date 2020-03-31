package contentful

import (
	"fmt"
)

// OrganizationsService service
type OrganizationsService service

// Organization model
type Organization struct {
	Sys  *Sys   `json:"sys"`
	Name string `json:"name"`
}

// GetVersion returns entity version
func (o *Organization) GetVersion() int {
	version := 1
	if o.Sys != nil {
		version = o.Sys.Version
	}

	return version
}

// List returns an organizations collection
func (service *OrganizationsService) List() *Collection {
	path := fmt.Sprintf("/organizations")
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// GetUsage returns the usage of the specified organization
func (service *OrganizationsService) GetUsage(organizationID, orderBy, metric, startAt, endAt string) *Collection {
	path := fmt.Sprintf(
		"/organizations/%s/space_periodic_usages?order=%s&metric[in]=%s&dateRange.startAt=%s&dateRange.endAt=%s",
		organizationID,
		orderBy,
		metric,
		startAt,
		endAt,
	)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}
