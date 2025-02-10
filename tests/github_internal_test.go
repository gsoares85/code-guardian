package tests

import (
	"github.com/gsoares85/code-guardian/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRepositoryFilesRecursive(t *testing.T) {
	files, err := mocks.MockGetRepositoryFilesRecursive("example", "repo")

	assert.Nil(t, err)
	assert.Greater(t, len(files), 0, "Expected to fetch files")
	assert.Contains(t, files, "src/main.go")
}

func TestGetFileContent(t *testing.T) {
	content, err := mocks.MockGetFileContent("example", "repo", "src/main.go")

	assert.Nil(t, err)
	assert.Contains(t, content, "package main", "Expected Go package definition")
}
