package main

import (
	"fmt"
	"net/http"

	"appengine/datastore"
	bigquery "code.google.com/p/google-api-go-client/bigquery/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func GetHistoricArchive(mc *MyContext) {
	bigQueryService := getBigQueryService(mc)
	queryResponse := queryHistoricArchive(mc, bigQueryService)
	updateHistoricArchive(mc, queryResponse)
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
func queryHistoricArchive(mc *MyContext, bigQueryService *bigquery.Service) *bigquery.QueryResponse {
	logins := loadMemberLogins(mc)

	query := fmt.Sprintf(HistoricArchiveQuery, logins, logins)
	mc.Infof("the query %s", query)

	queryRequest := &bigquery.QueryRequest{
		Query: query,
		Kind:  "igquery#queryRequest",
	}

	jobsService := bigquery.NewJobsService(bigQueryService)
	jobsQueryCall := jobsService.Query(appengine.AppID(mc.StdContext), queryRequest)
	queryResponse, err := jobsQueryCall.Do()

	if err != nil {
		mc.Infof("failed to query historic archive %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	return queryResponse
}

func loadMemberLogins(mc *MyContext) string {
	key := datastore.NewKey(mc.Context, "Configuration", "organization", 0, nil)

	var o Organization
	err := datastore.Get(mc.Context, key, &o)

	if err != nil {
		mc.Infof("could not load members %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	logins := ""

	for i, m := range o.Members {

		if i != 0 {
			logins += ","
		}

		if m.Login != "" {
			logins += "\"" + m.Login + "\""
		}
	}

	return logins

}

func updateHistoricArchive(mc *MyContext, queryResponse *bigquery.QueryResponse) {

	mc.Infof("Job completed? %#v", queryResponse.JobComplete)
	mc.Infof("Total rows %d", queryResponse.TotalRows)

	var keys []*datastore.Key
	var historicArchives []HistoricArchive

	for _, row := range queryResponse.Rows {
		keys = append(keys, datastore.NewKey(mc.Context, HistoricArchiveEntityKind, row.F[0].V.(string), 0, nil))
		historicArchives = append(historicArchives, HistoricArchive{
			PayloadAction:                   row.F[1].V.(string),
			PayloadPullRequestMerged:        row.F[2].V.(string),
			PayloadPullRequestTitle:         row.F[3].V.(string),
			PayloadPullRequestUrl:           row.F[4].V.(string),
			PayloadPullRequestUserLogin:     row.F[5].V.(string),
			PayloadPullRequestMergedByLogin: row.F[6].V.(string),
			PayloadPullRequestMergedAt:      row.F[1].V.(string),
		})
	}

	_, err := datastore.PutMulti(mc.Context, keys, historicArchives)

	if err != nil {
		mc.Infof("could not save historic archive %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

}
