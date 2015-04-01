package main

import (
	bigquery "code.google.com/p/google-api-go-client/bigquery/v2"

	"appengine"
	"appengine/datastore"
)

const (
	HistoricArchiveEntityKind = "HistoricArchive"
)

type HistoricArchive struct {
	PayloadAction                   string
	PayloadPullRequestMerged        string
	PayloadPullRequestTitle         string
	PayloadPullRequestUrl           string
	PayloadPullRequestUserLogin     string
	PayloadPullRequestMergedByLogin string
	PayloadPullRequestMergedAt      string
}

func LoadHistoricArchive(c appengine.Context) ([]HistoricArchive, error) {
	var ha []HistoricArchive

	q := datastore.NewQuery(HistoricArchiveEntityKind)

	_, err := q.GetAll(c, &ha)

	return ha, err
}

func UpdateHistoricArchive(c appengine.Context, queryResponse *bigquery.QueryResponse) error {
	var kBatch []*datastore.Key
	var haBatch []HistoricArchive
	batchSize := 0

	for _, row := range queryResponse.Rows {
		key := datastore.NewKey(c, HistoricArchiveEntityKind, row.F[0].V.(string), 0, nil)

		ha := HistoricArchive{
			PayloadAction:                   row.F[1].V.(string),
			PayloadPullRequestMerged:        row.F[2].V.(string),
			PayloadPullRequestTitle:         row.F[3].V.(string),
			PayloadPullRequestUrl:           row.F[4].V.(string),
			PayloadPullRequestUserLogin:     row.F[5].V.(string),
			PayloadPullRequestMergedByLogin: row.F[6].V.(string),
			PayloadPullRequestMergedAt:      row.F[1].V.(string),
		}

		if batchSize < 500 {
			kBatch = append(kBatch, key)
			haBatch = append(haBatch, ha)
			batchSize++
		} else {
			err := UpdateHistoricArchiveBatch(c, kBatch, haBatch)

			if err != nil {
				c.Infof("Could not update historic archive: %#v", err)
				return err
			}

			kBatch = []*datastore.Key{}
			haBatch = []HistoricArchive{}
			batchSize = 0
		}
	}

	if batchSize > 0 {
		err := UpdateHistoricArchiveBatch(c, kBatch, haBatch)

		if err != nil {
			c.Infof("Could not update historic archive: %#v", err)
			return err
		}
	}

	return nil

}

func UpdateHistoricArchiveBatch(c appengine.Context, kBatch []*datastore.Key, haBatch []HistoricArchive) error {
	_, err := datastore.PutMulti(c, kBatch, haBatch)

	if err != nil {
		return err
	}

	return nil
}
