# contributron-golang-api

Hopefully one day will become an open source contribution dashboard API for a
group of people based on https://www.githubarchive.org/

## Development

Run with:

```bash
export GOROOT=~/YOURPATHTO/go_appengine/goroot
➜  gcloud preview app run . --appidentity-email-address <your_app_email_address>@developer.gserviceaccount.com --appidentity-private-key-path pem_file.pem
```

## Tests

```bash
export GOROOT=~/YOURPATHTO/go_appengine/goroot
goapp test
```

## References
- https://developer.github.com/v3/activity/events/types/#pullrequestevent
- https://github.com/google/google-api-go-client/blob/master/examples/bigquery.go
- https://github.com/google/google-api-go-client/blob/master/bigquery/v2/bigquery-gen.go
- https://cloud.google.com/bigquery/query-reference#where
- https://cloud.google.com/bigquery/docs/reference/v2/jobs/query
- https://cloud.google.com/appengine/docs/go/googlecloudstorageclient/getstarted
- https://godoc.org/golang.org/x/oauth2/google#AppEngineTokenSource

## Examples
- https://github.com/rails/jbuilder/pull/59
- https://api.github.com/repos/rails/jbuilder/pulls/59
