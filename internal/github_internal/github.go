package github_internal

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

func GetPullRequestFiles(owner string, repo string, prNumber int) ([]string, error) {
	client := NewGithubClient()
	files, _, err := client.PullRequests.ListFiles(context.Background(), owner, repo, prNumber, nil)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.GetFilename())
	}
	return fileNames, nil
}

func GetRepositoryFilesRecursive(owner, repo string) ([]string, error) {
	client := NewGithubClient()
	var files []string

	err := fetchFilesRecursive(client, owner, repo, "", &files)
	if err != nil {
		return nil, fmt.Errorf("error fetching repository files: %w", err)
	}

	return files, nil
}

func fetchFilesRecursive(client *github.Client, owner, repo, path string, files *[]string) error {
	contents, dirContents, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return err
	}

	if contents != nil {
		*files = append(*files, contents.GetPath())
		return nil
	}

	for _, item := range dirContents {
		if item.GetType() == "file" {
			*files = append(*files, item.GetPath())
		} else if item.GetType() == "dir" {
			fetchFilesRecursive(client, owner, repo, item.GetPath(), files)
		}
	}

	return nil
}

func GetFileContent(owner, repo, path string) (string, error) {
	client := NewGithubClient()
	fileContent, _, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return "", fmt.Errorf("error fetching file content: %w", err)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("error decoding file content: %w", err)
	}

	return content, nil
}
