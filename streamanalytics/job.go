package streamanalytics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cbsinteractive/azurermfutures-go/internal"
)

type StreamAnalyticsJobService interface {
	Get(subscriptionId string, resourceGroupName string, jobName string) (StreamAnalyticsJob, error)
}

type StreamAnalyticsJob struct {
	ID string `json:"id"`
}

func NewStreamAnalyticsJobService(c internal.InternalClient) StreamAnalyticsJobServiceImpl {
	return StreamAnalyticsJobServiceImpl{
		c: c,
	}
}

type StreamAnalyticsJobServiceImpl struct {
	c internal.InternalClient
}

func (i StreamAnalyticsJobServiceImpl) Get(subscriptionId string, resourceGroupName string, jobName string) (StreamAnalyticsJob, error) {
	path := fmt.Sprintf(
		"subscriptions/%s/resourcegroups/%s/providers/Microsoft.StreamAnalytics/streamingjobs/%s?api-version=2017-04-01-preview",
		subscriptionId,
		resourceGroupName,
		jobName,
	)
	resp, err := i.c.Call(path)
	if err != nil {
		return StreamAnalyticsJob{}, err
	}
	defer resp.Response.Body.Close()
	if resp.Response.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Response.Body)
		if err != nil {
			fmt.Println(err)
			return StreamAnalyticsJob{}, err
		}
		return StreamAnalyticsJob{}, fmt.Errorf("failed to get job: %s", body)
	}
	job := StreamAnalyticsJob{}
	err = json.NewDecoder(resp.Response.Body).Decode(&job)
	if err != nil {
		return StreamAnalyticsJob{}, err
	}
	return job, nil
}
