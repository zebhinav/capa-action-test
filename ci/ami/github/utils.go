package github

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

func GetGithubClientCtx(token string) (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
}

func ListRepos(client *github.Client, ctx context.Context) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.List(ctx, os.Getenv("GITHUB_REPOSITORY_OWNER"), nil)
	return repos, err
}

func CreateIssue(client *github.Client, ctx context.Context) {
	MyIssue := struct {
		Title    string
		Body     string
		Labels   []string
		Assignee string
	}{
		Title:    "Test issue",
		Body:     "Test issue body",
		Labels:   []string{"test-issue-label"},
		Assignee: "",
	}
	testIssue := &github.IssueRequest{
		Title:    &MyIssue.Title,
		Body:     &MyIssue.Body,
		Labels:   &MyIssue.Labels,
		Assignee: &MyIssue.Assignee,
	}

	repoSplit := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	client.Issues.Create(ctx, repoSplit[0], repoSplit[1], testIssue)
}
