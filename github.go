package main
import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	*github.Client
	Owner string
	Repo string
}

func NewClient(token, owner, repo string) *GithubClient {
	client :=  &GithubClient{
		Client: getClient(token),
		Owner: owner,
		Repo: repo,
	}
	return client
}

func getClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}

func (client GithubClient) CreateIssue(title,body string) (err error) {
	service := client.Issues
	issue := github.IssueRequest{
		Title: title,
		Body: body,
		Labels: []string{UmengLabel},
	}
	_, err = service.Create(client.Owner, client.Repo, &issue)
	return
}