package tests

import (
	"github.com/gsoares85/code-guardian/internal/openai"
	"github.com/gsoares85/code-guardian/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyzeCodeWithAI(t *testing.T) {
	response, err := mocks.MockAnalyzeCodeWithAI("package main\nfunc main() {}", "Analyze this code")

	assert.Nil(t, err)
	assert.NotEmpty(t, response, "Expected AI analysis response")
}

func TestSplitLargeCode(t *testing.T) {
	longCode := "word " + "word " + "word " + "word " // Simulating large input
	chunks := openai.SplitLargeCode(longCode, 2)

	assert.Greater(t, len(chunks), 1, "Expected multiple chunks")
}
