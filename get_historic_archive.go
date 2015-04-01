package main

import (
	"fmt"
	"net/http"

	bigquery "code.google.com/p/google-api-go-client/bigquery/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func GetHistoricArchive(mc *MyContext) {
	bigQueryService := getBigQueryService(mc)

	queryResponse, err := queryHistoricArchive(mc, bigQueryService)

	if err != nil {
		mc.Infof("Could not query historic archive: %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	mc.Infof("Job completed?: %#v", queryResponse.JobComplete)
	mc.Infof("Total rows: %d", queryResponse.TotalRows)

	err = UpdateHistoricArchive(mc, queryResponse)

	if err != nil {
		mc.Infof("Could not update historic archive: %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

}

func getBigQueryService(mc *MyContext) *bigquery.Service {

	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(mc.StdContext, "https://www.googleapis.com/auth/bigquery"),
			Base:   &urlfetch.Transport{Context: mc.StdContext},
		},
	}

	service, err := bigquery.New(hc)

	if err != nil {
		mc.Infof("could not crate big query service %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	return service
}

// TODO: What id the query result size is bigger than 1K? learn how to do query
// pages result sets in the API
func queryHistoricArchive(mc *MyContext, bigQueryService *bigquery.Service) (*bigquery.QueryResponse, error) {
	organization, err := LoadOrganization(mc)

	if err != nil {
		mc.Infof("Could not load organization: %#v", err)
		return nil, err
	}

	logins := extractMemberLogins(organization)

	query := fmt.Sprintf(HistoricArchiveQuery, logins, logins)
	mc.Infof("The query: %s", query)

	queryRequest := &bigquery.QueryRequest{
		Query: query,
		Kind:  "igquery#queryRequest",
	}

	jobsService := bigquery.NewJobsService(bigQueryService)
	jobsQueryCall := jobsService.Query(appengine.AppID(mc.StdContext), queryRequest)
	queryResponse, err := jobsQueryCall.Do()

	if err != nil {
		mc.Infof("Could not call query job: %#v", err)
		return nil, err
	}

	return queryResponse, nil
}

func extractMemberLogins(org *Organization) string {
	logins := ""

	for i, m := range org.Members {

		if i != 0 {
			logins += ","
		}

		if m.Login != "" {
			logins += "\"" + m.Login + "\""
		}
	}

	return logins
}
