package contentful

import "fmt"

// UsagesService service
type UsagesService service

// Usage model
type Usage struct {
	Sys           *Sys              `json:"sys"`
	UnitOfMeasure string            `json:"unitOfMeasure"`
	Metric        string            `json:"metric"`
	DateRange     DateRange         `json:"dateRange"`
	TotalUsage    int               `json:"usage"`
	UsagePerDay   map[string]string `json:"usagePerDay"`
}

// DateRange model
type DateRange struct {
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
}

// GetOrganizationUsage returns the usage of the specified organization
func (service *UsagesService) GetOrganizationUsage(organizationID, orderBy, metric, startAt, endAt string) *Collection {
	path := fmt.Sprintf(
		"/organizations/%s/organization_periodic_usages?order=%s&metric[in]=%s&dateRange.startAt=%s&dateRange.endAt=%s",
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

// GetSpaceUsage returns the organization usage by space
func (service *UsagesService) GetSpaceUsage(organizationID, orderBy, metric, startAt, endAt string) *Collection {
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
