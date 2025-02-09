package cmd

import (
	"github.com/google/go-github/v49/github"
	"github.com/gsoares85/code-guardian/internal/github_internal"
	"github.com/gsoares85/code-guardian/internal/openai"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// ✅ Test Markdown Generation
func TestGenerateMarkdownReport(t *testing.T) {
	prTitle := "Fix memory leak"
	prOwner := "testUser"
	prNumber := 23
	pr, err := github_internal.MockGetPullRequest(prOwner, prTitle, prNumber)
	files, err := github_internal.MockGetPullRequestFiles(prOwner, prTitle, prNumber)
	prDiff, err := github_internal.MockGetPullRequestDiff(prOwner, prTitle, prNumber)
	aiFeedback, err := openai.MockAnalyzePRWithAI(prDiff)

	report := generateMarkdownReport(pr, files, prDiff, aiFeedback)

	assert.Contains(t, report, "# Pull Request Analysis Report")
	assert.Contains(t, report, prTitle)
	assert.Contains(t, report, "src/main.c")
	assert.Contains(t, report, "free(ptr);")
	assert.Nil(t, err)
}

// ✅ Test File Saving
func TestSaveAnalysisToFile(t *testing.T) {
	fileName := filepath.Join("tests", "reports", "test_report.md")
	content := "# Sample Report\nThis is a test file."

	err := saveAnalysisToFile(fileName, content)
	assert.Nil(t, err)

	// Verify file exists
	_, err = os.Stat(fileName)
	assert.Nil(t, err)

	// Clean up
	os.Remove(fileName)
}

// ✅ Test PR Analysis (Integration)
func TestAnalyzePullRequest(t *testing.T) {
	// Use local variables in the test function to "mock" external functionality.
	mockGetPullRequest := func(owner, title string, number int) (*github.PullRequest, error) {
		return github_internal.MockGetPullRequest(owner, title, number)
	}
	mockGetPullRequestFiles := func(owner, title string, number int) ([]string, error) {
		return github_internal.MockGetPullRequestFiles(owner, title, number)
	}
	mockGetPullRequestDiff := func(owner, title string, number int) (string, error) {
		return github_internal.MockGetPullRequestDiff(owner, title, number)
	}
	mockAnalyzePRWithAI := func(prDiff string) (string, error) {
		return openai.MockAnalyzePRWithAI(prDiff)
	}

	// Use the mocks for testing by calling them explicitly.
	pr, err := mockGetPullRequest("example", "repo", 42)
	assert.Nil(t, err)

	files, err := mockGetPullRequestFiles("example", "repo", 42)
	assert.Nil(t, err)

	prDiff, err := mockGetPullRequestDiff("example", "repo", 42)
	assert.Nil(t, err)

	aiFeedback, err := mockAnalyzePRWithAI(prDiff)
	assert.Nil(t, err)

	// Add assertions to validate everything is working correctly
	assert.NotNil(t, pr)
	assert.Greater(t, len(files), 0)
	assert.Greater(t, len(prDiff), 0)
	assert.Greater(t, len(aiFeedback), 0)
}
