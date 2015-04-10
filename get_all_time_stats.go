package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
)

func GetAllTimeStats(mc *MyContext) {
	events, err := LoadEvents(mc.Context)

	if err != nil {
		mc.Infof("could not load historic archive %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	mc.Infof("Loaded %d historic archive entries", len(events))

	users := map[string]*RankedUser{}
	rankedUsers := []*RankedUser{}

	for _, event := range events {
		userLogin := event.PullRequest.User.Login

		if _, ok := users[userLogin]; !ok {
			rankedUser := &RankedUser{Name: userLogin}
			users[userLogin] = rankedUser
			rankedUsers = append(rankedUsers, rankedUser)
		}

		users[userLogin].TotalPRs++
		users[userLogin].PRs =
			append(users[userLogin].PRs, event.PullRequest.URL)
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
	events, err := LoadEvents(mc.Context)

	if err != nil {
		mc.Infof("could not load historic archive %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	mc.Infof("Loaded %d historic archive entries", len(events))

	users := map[string]*RankedUser{}
	rankedUsers := []*RankedUser{}

	for _, event := range events {
		matched, err := regexp.MatchString("magma|crowdint", event.PullRequest.URL)

		if err != nil {
			mc.Infof("Error matching string: %#v", err)
			http.Error(mc.W, err.Error(), http.StatusInternalServerError)
		}

		if !matched {
			userLogin := event.PullRequest.User.Login

			if _, ok := users[userLogin]; !ok {
				rankedUser := &RankedUser{Name: userLogin}
				users[userLogin] = rankedUser
				rankedUsers = append(rankedUsers, rankedUser)
			}

			users[userLogin].TotalPRs++
			users[userLogin].PRs =
				append(users[userLogin].PRs, event.PullRequest.URL)
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
