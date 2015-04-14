package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lostisland/go-sawyer/hypermedia"
)

var (
	PublicMembersEndpoint = "https://api.github.com/orgs/%s/public_members"
)

func GetPublicMembersList(mc *MyContext) {
	org := &Organization{
		Name: "crowdint",
	}

	members, err := callPublicMembersListEndpoint(mc, org)

	org.Members = members

	if err != nil {
		mc.Infof("Could not call end point: %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	err = UpdateOrganizationMembers(mc.Context, org)

	if err != nil {
		mc.Infof("Could not update organization configuration: %#v", err)
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

}

func callPublicMembersListEndpoint(mc *MyContext, org *Organization) ([]Member, error) {
	url := fmt.Sprintf(PublicMembersEndpoint, org.Name)

	var members []Member
	err := requestMembers(mc, url, &members)

	if err != nil {
		return nil, err
	}

	return members, nil
}

func requestMembers(mc *MyContext, url string, members *[]Member) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		mc.Infof("Could not create request: %#v", err)
		return err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := mc.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		mc.Infof("Could not make request: %#v", err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		mc.Infof("Could not read response body: %#v", err)
		return err
	}

	var newMembers []Member
	err = json.Unmarshal(body, &newMembers)

	if err != nil {
		mc.Infof("Could not unmarshal body: %#v", err)
		return err
	}

	for _, member := range newMembers {
		*members = append(*members, member)
	}

	rels := hypermedia.HyperHeaderRelations(resp.Header, hypermedia.NewRels())
	nextURL, err := rels.Rel("next", nil)
	if err == nil {
		err := requestMembers(mc, nextURL.String(), members)
		if err != nil {
			return err
		}
	}

	return nil
}
