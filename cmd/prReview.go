/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gsoares85/code-guardian/internal/github"
	"github.com/gsoares85/code-guardian/internal/openai"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
)

// prReviewCmd represents the prReview command
var prReviewCmd = &cobra.Command{
	Use:   "pr-review [owner] [repo] [pr-number] [flags]",
	Short: "Analysis of a pull request",
	Long:  `Use this command to perform a complete code review of a pull request.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		owner := args[0]
		repo := args[1]
		prNumber, err := strconv.Atoi(args[2])
		if err != nil {
			color.Red("Error: Invalid PR number")
			return
		}

		saveOutput, _ := cmd.Flags().GetBool("output")

		pr, err := github.GetPullRequest(owner, repo, prNumber)
		if err != nil {
			color.Red("Error getting PR: %s\n", err)
			return
		}

		color.Cyan("PR %d - %s\n", prNumber, pr.GetTitle())
		color.Green("Author: %s\n", pr.GetUser().GetLogin())
		color.Yellow("Created at: %s\n", pr.GetCreatedAt())
		color.Blue("Link: %s\n", pr.GetHTMLURL())

		files, err := github.GetPullRequestFiles(owner, repo, prNumber)
		if err != nil {
			color.Red("Error getting PR files: %s\n", err)
			return
		}

		color.Magenta("Files changed: %d\n", len(files))
		for _, file := range files {
			fmt.Println("  -", file)
		}

		prDiff, err := github.GetPullRequestDiff(owner, repo, prNumber)
		if err != nil {
			color.Red("Error getting PR diff: %s\n", err)
			return
		}

		color.Magenta("Code changes")
		fmt.Println(len(prDiff))

		color.Blue("Sending for AI analysis")
		aiFeedback, err := openai.AnalyzePRWithAI(prDiff)
		if err != nil {
			color.Red("Error analyzing PR: %s\n", err)
			return
		}
		color.Green("AI feedback:")
		fmt.Println(aiFeedback)

		outputContent := fmt.Sprintf(`# Pull Request analysis report
## PR %d - %s
- **Author:** %s
- **Created at:** %s
- **PR Link:** [%s](%s)

## 📂 Files changed:
%s

## 🔍 Code changes:
'''diff
%s
'''
## 🤖 AI Suggestions:
%s
		`, prNumber, pr.GetTitle(), pr.GetUser().GetLogin(), pr.GetCreatedAt(), pr.GetHTMLURL(), pr.GetHTMLURL(),
			formatFileList(files),
			prDiff,
			aiFeedback)
		if saveOutput {
			timestamp := pr.GetCreatedAt().Format("20060102-150405")
			outputFile := fmt.Sprintf("%s-%s_%s_%d.md", timestamp, repo, owner, prNumber)
			outputPath := filepath.Join("reports", "pr", outputFile)

			err := saveAnalysisToFile(outputPath, outputContent)
			if err != nil {
				color.Red("Error saving analysis to file: %s\n", err)
				return
			}
			color.Green("Analysis saved to file: %s\n", outputPath)
		}
	},
}

func formatFileList(files []string) string {
	formated := ""
	for _, file := range files {
		formated += fmt.Sprintf("- %s\n", file)
	}
	return formated
}

func saveAnalysisToFile(fineName, content string) error {
	dir := filepath.Dir(fineName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(fineName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func init() {
	rootCmd.AddCommand(prReviewCmd)
	prReviewCmd.Flags().BoolP("output", "o", false, "Save analysis to a Markdown file in ./reports/pr/")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prReviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prReviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
