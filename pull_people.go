package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"appengine/datastore"
)

func GetPublicMembersList(mc *MyContext) {
	members := callPublicMembersListEndpoint(mc)
	updateMembers(mc, members)
}

func callPublicMembersListEndpoint(mc *MyContext) []Member {
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/crowdint/public_members", nil)

	if err != nil {
		mc.Infof("%#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := mc.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		mc.Infof("%#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		mc.Infof("%#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	var members []Member

	err = json.Unmarshal(body, &members)

	if err != nil {
		mc.Infof("%#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	return members
}

func updateMembers(mc *MyContext, members []Member) {
	organization := &Organization{Name: "crowdint", Members: members}
	mc.Infof("to-create %#v", organization)

	key := datastore.NewKey(mc.Context, "Configuration", "organization", 0, nil)
	_, err := datastore.Put(mc.Context, key, organization)

	if err != nil {
		mc.Infof("%#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	var fo Organization
	err = datastore.Get(mc.Context, key, &fo)

	if err != nil {
		mc.Infof("%#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	mc.Infof("found %#v", fo)
}
