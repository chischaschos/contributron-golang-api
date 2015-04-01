package main

import (
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
