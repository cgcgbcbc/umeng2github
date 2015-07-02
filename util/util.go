package util

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubClient github.Client

func GetGithubClient(token string) GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return GithubClient(*github.NewClient(tc))
}

func (client GithubClient) CheckCurrentUser(users []string) (user *github.User, result bool) {
	if len(users) == 0 {
		result = true
		return
	}
	user, _, err := client.Users.Get("")
	handleError(err)

	for _, u := range users {
		if u == *user.Login {
			result = true
			return
		}
	}
	result = false
	return
}

func (client GithubClient) CheckCurrentOrg(orgs []string) bool {
	if len(orgs) == 0 {
		return true
	}
	option := &github.ListOrgMembershipsOptions{State: "active"}
	memberships, _, err := client.Organizations.ListOrgMemberships(option)
	handleError(err)

	orgsMap := make(map[string]struct{}, len(orgs))
	for _, org := range orgs {
		orgsMap[org] = struct{}{}
	}

	for _, m := range memberships {
		_, ok := orgsMap[*m.Organization.Login]
		if ok {
			return ok
		}
	}
	return false
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
