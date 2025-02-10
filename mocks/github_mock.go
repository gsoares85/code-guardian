package mocks

import "github.com/google/go-github/v49/github"

func MockGetPullRequest(owner string, repo string, prNumber int) (*github.PullRequest, error) {
	return &github.PullRequest{
		Number:  github.Int(prNumber),
		Title:   github.String(repo),
		User:    &github.User{Login: github.String(owner)},
		HTMLURL: github.String("https://github.com/example/repo/pull/42"),
	}, nil
}

func MockGetPullRequestFiles(_ string, _ string, _ int) ([]string, error) {
	return []string{"src/main.c", "include/utils.h"}, nil
}

func MockGetPullRequestDiff(_ string, _ string, _ int) (string, error) {
	return `@@ -23,6 +23,7 @@
 void fixMemory() {
     int* ptr = malloc(100);
     ptr[0] = 1;
+    free(ptr);
 }`, nil
}

func MockGetRepositoryFilesRecursive(owner, repo string) ([]string, error) {
	return []string{
		"src/main.go",
		"src/utils/helper.go",
		"config/config.yaml",
		"README.md",
	}, nil
}

func MockGetFileContent(owner, repo, path string) (string, error) {
	if path == "src/main.go" {
		return `package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}`, nil
	}

	return "", nil
}
