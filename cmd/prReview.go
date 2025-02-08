/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/gsoares85/code-guardian/internal/github"
	"github.com/gsoares85/code-guardian/internal/openai"
	"github.com/spf13/cobra"
	"strconv"
)

// prReviewCmd represents the prReview command
var prReviewCmd = &cobra.Command{
	Use:   "pr-review [owner] [repo] [pr-number]",
	Short: "Analysis of a pull request",
	Long:  `Use this command to perform a complete code review of a pull request.`,
	Run: func(cmd *cobra.Command, args []string) {
		owner := args[0]
		repo := args[1]
		prNumber, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Print("Error: Invalid PR number")
			return
		}

		pr, err := github.GetPullRequest(owner, repo, prNumber)
		if err != nil {
			fmt.Printf("Error getting PR: %s\n", err)
			return
		}
		prDiff, err := github.GetPullRequestDiff(owner, repo, prNumber)
		if err != nil {
			fmt.Printf("Error getting PR diff: %s\n", err)
		}

		aiFeedback, err := openai.AnalyzePRWithAI(prDiff)
		if err != nil {
			fmt.Printf("Error analyzing PR: %s\n", err)
		}

		fmt.Printf("PR %d - %s\n", prNumber, pr.GetTitle())
		fmt.Printf("Author: %s\n", pr.GetUser().GetLogin())
		fmt.Printf("Created at: %s\n", pr.GetCreatedAt())
		fmt.Printf("Link: %s\n", pr.GetHTMLURL())
		fmt.Println("Code changes")
		fmt.Println(len(prDiff))
		fmt.Println("AI feedback:")
		fmt.Println(aiFeedback)
	},
}

func init() {
	rootCmd.AddCommand(prReviewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prReviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prReviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
