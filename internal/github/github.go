package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v49/github"
	"github.com/gsoares85/code-guardian/config"
	"golang.org/x/oauth2"
)

func NewGithubClient() *github.Client {
	token := config.GetEnv("GITHUB_TOKEN")
	if token == "" {
		fmt.Print("Error: GITHUB_TOKEN not found")
		return nil
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return github.NewClient(tc)
}

func GetPullRequest(owner string, repo string, prNumber int) (*github.PullRequest, error) {
	client := NewGithubClient()
	if client == nil {
		return nil, fmt.Errorf("error creating Github client")
	}

	pr, _, err := client.PullRequests.Get(context.Background(), owner, repo, prNumber)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func GetPullRequestDiff(owner string, repo string, prNumber int) (string, error) {
	client := NewGithubClient()
	diff, _, err := client.PullRequests.GetRaw(context.Background(), owner, repo, prNumber, github.RawOptions{Type: github.Diff})
	if err != nil {
		return "", err
	}
	return diff, nil
}
