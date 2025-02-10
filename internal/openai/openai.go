package openai

import (
	"context"
	"fmt"
	"github.com/gsoares85/code-guardian/config"
	"github.com/sashabaranov/go-openai"
	"strings"
)

func AnalyzePRWithAI(diff string) (string, error) {
	fmt.Println("Starting AI analysis")
	apiKey := config.GetEnv("OPENAI_API_KEY")

	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not found")
	}

	client := openai.NewClient(apiKey)

	prompt := fmt.Sprintf("Do a complete review of the following PR diff. Suggeting me comments, improviments and code quality:\n\n%s", diff)
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are an assistant specialized in PR code review. Your goal is to provide a complete review of the PR diff. You should suggest improvements and code quality. If you don't know what to say, just say.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 500,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}
	fmt.Println("AI analysis finished")

	return resp.Choices[0].Message.Content, nil
}

func AnalyzeCodeWithAI(code string, prompt string) (string, error) {
	apiKey := config.GetEnv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing OpenAI API key")
	}

	client := openai.NewClient(apiKey)

	codeChunks := SplitLargeCode(code, 3000)

	var fullResponse strings.Builder

	for _, chunk := range codeChunks {
		requestPrompt := fmt.Sprintf("%s\n\n%s", prompt, chunk)

		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are a senior software engineer reviewing code."},
				{Role: "user", Content: requestPrompt},
			},
		})

		if err != nil {
			return "", err
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("empty response from OpenAI")
		}

		fullResponse.WriteString(resp.Choices[0].Message.Content + "\n\n")
	}

	return fullResponse.String(), nil
}

func SplitLargeCode(code string, maxTokens int) []string {
	words := strings.Fields(code)
	var chunks []string

	for i := 0; i < len(words); i += maxTokens {
		end := i + maxTokens
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}

	return chunks
}
