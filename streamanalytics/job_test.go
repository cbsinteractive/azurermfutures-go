package streamanalytics_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	azurermfutures "github.com/cbsinteractive/azurermfutures-go"
	"github.com/cbsinteractive/azurermfutures-go/internal"
	"github.com/cbsinteractive/azurermfutures-go/streamanalytics"
)

func TestGetJobs(t *testing.T) {
	const subscriptionIDTemplate = "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingjobs/%s"
	tests := []struct {
		description       string
		id                string
		subscriptionId    string
		tenantId          string
		clientId          string
		clientSecret      string
		resourceGroupName string
		jobName           string
	}{
		{
			"FindJobOkay",
			fmt.Sprintf(subscriptionIDTemplate, "foosubscriptionid", "foorg", "foojobname"),
			"foosubscriptionid",
			"footenantid",
			"fooclientid",
			"fooclientsecret",
			"foorg",
			"foojobname",
		},
	}
	type handler func(rw http.ResponseWriter, req *http.Request) bool
	var handleAuthRequest handler = func(rw http.ResponseWriter, req *http.Request) bool {
		if strings.Contains(req.URL.Path, "oauth2/v2.0/token") {
			rw.WriteHeader(http.StatusOK)
			d, _ := json.Marshal(internal.AuthData{
				TokenType:    "Bearer",
				ExpiresIn:    3599,
				ExtExpiresIn: 3599,
				AccessToken:  "some token",
			})
			rw.Write(d)
			return true
		}
		return false
	}
	var handlerGetJobRequest handler = func(rw http.ResponseWriter, req *http.Request) bool {
		if req.Method != "GET" {
			return false
		}
		re := regexp.MustCompile(`(?i)^(/subscriptions/([-a-z0-9]+)/resourcegroups/([-a-z]+)/providers/microsoft\.streamanalytics/streamingjobs/([-a-z0-9]+))\?api-version=2017-04-01-preview`)
		m := re.FindSubmatch([]byte(req.URL.Path))
		if m != nil {
			rw.WriteHeader(http.StatusOK)
			d, _ := json.Marshal(streamanalytics.StreamAnalyticsJob{
				ID: fmt.Sprintf(subscriptionIDTemplate, m[2], m[3], m[4]),
			})
			rw.Write(d)
			return true
		}
		return false
	}

	server := httptest.NewTLSServer(
		http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				var result bool
				if result = handleAuthRequest(rw, req); result {
					return
				}
				if result = handlerGetJobRequest(rw, req); result {
					return
				}
				t.Fatalf("Unexpected URL in handler: %s", req.URL.String())
			},
		),
	)
	defer server.Close()
	for _, d := range tests {
		t.Run(d.description, func(t *testing.T) {
			u, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			c, err := azurermfutures.New(
				server.Client(),
				azurermfutures.SetHost(u.Host),
				azurermfutures.SetAuthHost(u.Host),
				azurermfutures.SetTenantID(d.tenantId),
				azurermfutures.SetClientID(d.clientId),
				azurermfutures.SetClientSecret(d.clientSecret),
			)
			if err != nil {
				t.Fatal(err)
			}
			job, err := c.StreamAnalyticsJob.Get(d.subscriptionId, d.resourceGroupName, d.jobName)
			if err != nil {
				t.Fatal(err)
			}
			if job.ID != d.id {
				t.Errorf("Unexpected job ID. Expected %s. Got %s", d.id, job.ID)
			}
		})
	}
}
