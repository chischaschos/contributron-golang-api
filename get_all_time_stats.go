package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
)

var (
	ReposToIgnoreRegExp = regexp.MustCompile("magma|crowdint")
)

func GetAllTimeStats(mc *MyContext) {
	events, err := LoadEvents(mc.Context)

	if err != nil {
		mc.Infof("could not load historic archive %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	mc.Infof("Loaded %d historic archive entries", len(events))

	rankedUsers, err := analyzeEvents(mc, events)

	if err != nil {
		mc.Infof("could not analize events %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	for i, r := range rankedUsers {
		r.Rank = i + 1
	}

	bytes, err := json.MarshalIndent(rankedUsers, "", "\t")

	mc.W.Header().Add("Content-Type", "application/json")
	_, err = mc.W.Write(bytes)

	if err != nil {
		mc.Infof("could not write response %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}
}

func analyzeEvents(mc *MyContext, events []Event) ([]*RankedUser, error) {
	users := map[string]*RankedUser{}
	rankedUsers := []*RankedUser{}

	org, _ := LoadOrganization(mc)

	cm := map[string]bool{}

	for _, m := range org.Members {
		cm[m.Login] = true
	}

	for _, event := range events {
		userLogin := event.PullRequest.User.Login
		mergedByLogin := event.PullRequest.MergedBy.Login
		pr := PR{URL: event.URL}

		if _, ok := cm[userLogin]; !ok {
			continue
		}

		// Initialize this user structure
		if _, ok := users[userLogin]; !ok {
			rankedUser := &RankedUser{Name: userLogin}
			users[userLogin] = rankedUser
			rankedUsers = append(rankedUsers, rankedUser)
		}

		if ReposToIgnoreRegExp.MatchString(event.URL) {
			pr.Notes = append(pr.Notes, "Ignored repo "+ReposToIgnoreRegExp.String())

		} else if userLogin == mergedByLogin {
			pr.Notes = append(pr.Notes, "Ignored self merge")

		} else {
			users[userLogin].TotalPRs++
			pr.Notes = append(pr.Notes, "External collaboration")
		}

		users[userLogin].PRs = append(users[userLogin].PRs, pr)
	}

	sort.Sort(RankedUsers(rankedUsers))

	return rankedUsers, nil

}
