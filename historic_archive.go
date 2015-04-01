package main

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
