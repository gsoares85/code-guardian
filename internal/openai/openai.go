package openai

import (
	"context"
	"fmt"
	"github.com/gsoares85/code-guardian/config"
	"github.com/sashabaranov/go-openai"
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
