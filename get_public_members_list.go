package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	req, err := http.NewRequest("GET", fmt.Sprintf(PublicMembersEndpoint, org.Name), nil)

	if err != nil {
		mc.Infof("Could not create request: %#v", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := mc.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		mc.Infof("Could not make request: %#v", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		mc.Infof("Could not read response body: %#v", err)
		return nil, err
	}

	var members []Member

	err = json.Unmarshal(body, &members)

	if err != nil {
		mc.Infof("Could not unmarshal body: %#v", err)
		return nil, err
	}

	return members, nil
}
