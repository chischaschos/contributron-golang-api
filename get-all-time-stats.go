package main

import (
	"encoding/json"
	"net/http"
	"sort"

	"appengine/datastore"
)

type RankedUser struct {
	Name     string
	PRs      []string
	TotalPRs int
}

type RankedUsers []*RankedUser

func (a RankedUsers) Len() int           { return len(a) }
func (a RankedUsers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RankedUsers) Less(i, j int) bool { return a[i].TotalPRs < a[j].TotalPRs }

func GetAllTimeStats(mc *MyContext) {
	historicArchive := loadHistoricArchive(mc)

	mc.Infof("Loaded historic archive", len(historicArchive))

	users := map[string]*RankedUser{}
	rankedUsers := []*RankedUser{}

	for _, ha := range historicArchive {
		if _, ok := users[ha.PayloadPullRequestUserLogin]; !ok {
			rankedUser := &RankedUser{Name: ha.PayloadPullRequestUserLogin}
			users[ha.PayloadPullRequestUserLogin] = rankedUser
			rankedUsers = append(rankedUsers, rankedUser)
		}

		users[ha.PayloadPullRequestUserLogin].TotalPRs++
		users[ha.PayloadPullRequestUserLogin].PRs =
			append(users[ha.PayloadPullRequestUserLogin].PRs, ha.PayloadPullRequestUrl)
	}

	sort.Sort(RankedUsers(rankedUsers))

	bytes, err := json.MarshalIndent(rankedUsers, "", "\t")

	_, err = mc.W.Write(bytes)

	if err != nil {
		mc.Infof("could not write response %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

}

func loadHistoricArchive(mc *MyContext) []HistoricArchive {
	var historicArchive []HistoricArchive
	q := datastore.NewQuery(HistoricArchiveEntityKind)
	_, err := q.GetAll(mc.Context, &historicArchive)

	if err != nil {
		mc.Infof("could not load historic archive %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	return historicArchive
}
