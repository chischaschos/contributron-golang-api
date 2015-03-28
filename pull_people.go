package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Member struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"url"`
}

//func GetPublicMembersList(w http.ResponseWriter, r *http.Request) {
func GetPublicMembersList(mc *MyContext) {
	fmt.Println(callPublicMembersListEndpoint(mc))

}

func callPublicMembersListEndpoint(mc *MyContext) []Member {
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/crowdint/public_members", nil)

	if err != nil {
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := mc.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	var members []Member

	err = json.Unmarshal(body, &members)

	if err != nil {
		http.Error(mc.W, err.Error(), http.StatusInternalServerError)
	}

	return members
}
