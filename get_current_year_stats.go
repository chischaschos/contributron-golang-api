package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.google.com/p/google-api-go-client/bigquery/v2"
	"google.golang.org/appengine"
)

func GetCurrentYearArchive(mc *MyContext) {
	bigQueryService := getBigQueryService(mc)

	queryResponse, err := queryCurrentYearArchive(mc, bigQueryService)

	if err != nil {
		mc.Infof("Could not query currenty year archive: %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	mc.Infof("Job completed?: %#v", queryResponse.JobComplete)
	mc.Infof("Total rows: %d", queryResponse.TotalRows)
	mc.Infof("Token: %s", queryResponse.PageToken)

	events := []Event{}

	for _, row := range queryResponse.Rows {
		event := Event{
			ID:   row.F[0].V.(string),
			Type: row.F[1].V.(string),
		}

		if event.Type == "PullRequestEvent" {
			var payload Payload
			data := row.F[2].V.(string)

			err = json.Unmarshal([]byte(data), &payload)

			if err != nil {
				mc.Infof("Unmarshaling error", err)
				http.Error(mc.W, err.Error(), http.StatusInternalServerError)
			}
			events = append(events, event)
		}
	}

	err = UpdateEvents(mc.Context, events)

	if err != nil {
		mc.Infof("Could not update events: %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

}

func queryCurrentYearArchive(mc *MyContext, bigQueryService *bigquery.Service) (*bigquery.QueryResponse, error) {
	organization, err := LoadOrganization(mc)

	if err != nil {
		mc.Infof("Could not load organization: %#v", err)
		return nil, err
	}

	logins := extractMemberLogins(organization)

	query := fmt.Sprintf(CurrentYearQuery, logins)
	mc.Infof("The query: %s", query)

	queryRequest := &bigquery.QueryRequest{
		Query: query,
		Kind:  "igquery#queryRequest",
	}

	projectID := appengine.AppID(mc.StdContext)
	jobsService := bigquery.NewJobsService(bigQueryService)
	jobsQueryCall := jobsService.Query(projectID, queryRequest)
	queryResponse, err := jobsQueryCall.Do()

	if err != nil {
		mc.Infof("Could not call query job: %#v", err)
		return nil, err
	}

	return queryResponse, nil
}
