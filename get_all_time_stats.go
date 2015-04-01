package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
)

func GetAllTimeStats(mc *MyContext) {
	historicArchive, err := LoadHistoricArchive(mc.Context)

	if err != nil {
		mc.Infof("could not load historic archive %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	mc.Infof("Loaded %d historic archive entries", len(historicArchive))

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

func GetAllTimeStatsNoCrowd(mc *MyContext) {
	historicArchive, err := LoadHistoricArchive(mc.Context)

	if err != nil {
		mc.Infof("could not load historic archive %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	mc.Infof("Loaded %d historic archive entries", len(historicArchive))

	users := map[string]*RankedUser{}
	rankedUsers := []*RankedUser{}

	for _, ha := range historicArchive {
		matched, err := regexp.MatchString("magma|crowdint", ha.PayloadPullRequestUrl)

		if err != nil {
			mc.Infof("Error matching string: %#v", err)
			http.Error(mc.W, err.Error(), http.StatusInternalServerError)
		}

		if !matched {

			if _, ok := users[ha.PayloadPullRequestUserLogin]; !ok {
				rankedUser := &RankedUser{Name: ha.PayloadPullRequestUserLogin}
				users[ha.PayloadPullRequestUserLogin] = rankedUser
				rankedUsers = append(rankedUsers, rankedUser)
			}

			users[ha.PayloadPullRequestUserLogin].TotalPRs++
			users[ha.PayloadPullRequestUserLogin].PRs =
				append(users[ha.PayloadPullRequestUserLogin].PRs, ha.PayloadPullRequestUrl)
		}
	}

	sort.Sort(RankedUsers(rankedUsers))

	bytes, err := json.MarshalIndent(rankedUsers, "", "\t")

	_, err = mc.W.Write(bytes)

	if err != nil {
		mc.Infof("could not write response %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

}
