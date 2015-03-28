package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"code.google.com/p/appengine-go/appengine/aetest"
	"github.com/nicholasf/fakepoint"
)

func TestCallPublicMembersListEndpoint(t *testing.T) {
	maker := fakepoint.NewFakepointMaker()
	maker.NewGet("https://api.github.com/orgs/crowdint/public_members", 200).
		SetResponse(`
[
  {
    "login": "octocat",
    "id": 1,
    "avatar_url": "https://github.com/images/error/octocat_happy.gif",
    "gravatar_id": "",
    "url": "https://api.github.com/users/octocat",
    "html_url": "https://github.com/octocat",
    "followers_url": "https://api.github.com/users/octocat/followers",
    "following_url": "https://api.github.com/users/octocat/following{/other_user}",
    "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
    "organizations_url": "https://api.github.com/users/octocat/orgs",
    "repos_url": "https://api.github.com/users/octocat/repos",
    "events_url": "https://api.github.com/users/octocat/events{/privacy}",
    "received_events_url": "https://api.github.com/users/octocat/received_events",
    "type": "User",
    "site_admin": false
  }
]`).
		SetHeader("Content-Type", "application/json")

	c, err := aetest.NewContext(nil)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	mc := &MyContext{Context: c, Client: maker.Client()}
	wrapee := Wrap(GetPublicMembersList, mc)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/pull-people", nil)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	wrapee(w, req)

	t.Log(w.Body.String())

}
