package github_internal

import "github.com/google/go-github/v49/github"

// Mock function to simulate fetching a PR
func MockGetPullRequest(owner string, repo string, prNumber int) (*github.PullRequest, error) {
	return &github.PullRequest{
		Number:  github.Int(prNumber),
		Title:   github.String(repo),
		User:    &github.User{Login: github.String(owner)},
		HTMLURL: github.String("https://github.com/example/repo/pull/42"),
	}, nil
}

// Mock function to simulate fetching changed files
func MockGetPullRequestFiles(_ string, _ string, _ int) ([]string, error) {
	return []string{"src/main.c", "include/utils.h"}, nil
}

// Mock function to simulate fetching PR diff
func MockGetPullRequestDiff(_ string, _ string, _ int) (string, error) {
	return `@@ -23,6 +23,7 @@
 void fixMemory() {
     int* ptr = malloc(100);
     ptr[0] = 1;
+    free(ptr);
 }`, nil
}
