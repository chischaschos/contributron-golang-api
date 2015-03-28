package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	bigquery "code.google.com/p/google-api-go-client/bigquery/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	http.HandleFunc("/", doit)

	http.HandleFunc("/pull-people", Wrap(GetPublicMembersList, &MyContext{}))
}

// https://developer.github.com/v3/activity/events/types/#pullrequestevent
// https://github.com/google/google-api-go-client/blob/master/examples/bigquery.go
// https://github.com/google/google-api-go-client/blob/master/bigquery/v2/bigquery-gen.go
// https://cloud.google.com/bigquery/query-reference#where
// https://cloud.google.com/bigquery/docs/reference/v2/jobs/query
// https://cloud.google.com/appengine/docs/go/googlecloudstorageclient/getstarted
// https://godoc.org/golang.org/x/oauth2/google#AppEngineTokenSource
func doit(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, "https://www.googleapis.com/auth/bigquery"),
			Base:   &urlfetch.Transport{Context: c},
		},
	}

	bigqueryService, err := bigquery.New(hc)

	if err != nil {
		panic(err)
	}

	queryRequest := &bigquery.QueryRequest{
		Query: `
	SELECT type, JSON_EXTRACT(payload, '$.action') as action, JSON_EXTRACT(payload, '$.merged') as merged
  FROM TABLE_DATE_RANGE(githubarchive:day.events_,
                    TIMESTAMP('2015-03-01'),
                    TIMESTAMP('2015-03-27'))
  WHERE actor.login ='dhh' AND
    type = 'PullRequestEvent';
    `,
		Kind: "igquery#queryRequest",
	}

	jobsService := bigquery.NewJobsService(bigqueryService)
	jobsQueryCall := jobsService.Query(appengine.AppID(c), queryRequest)
	queryResponse, err := jobsQueryCall.Do()

	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, queryResponse.JobComplete)
	fmt.Fprintln(w, queryResponse.TotalRows)

	for _, row := range queryResponse.Rows {
		fmt.Fprintf(w, "%s %s %s \n", row.F[0], row.F[1], row.F[2])
	}
	fmt.Fprintln(w, queryResponse.TotalRows)
}
