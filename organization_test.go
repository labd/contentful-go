package contentful

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationsServiceList(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/organizations")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("organization.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	_, err = cma.Organizations.List().Next()
	assertions.Nil(err)
}

func TestOrganizationsService_GetUsage(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/organizations/organization_id/space_periodic_usages?order=-usage&metric[in]=cma,cpa,gql&dateRange.startAt=2020-01-01&dateRange.endAt=2020-01-03")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("usage_organization.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	res, err := cma.Organizations.GetUsage("organization_id", "-usage", "cma,cpa,gql", "2020-01-01", "2020-01-03").Next()
	assertions.Nil(err)

	usage := res.ToUsage()
	assertions.Equal(1, len(usage))
	assertions.Equal("<usage_metric_id>", usage[0].Sys.ID)
}
