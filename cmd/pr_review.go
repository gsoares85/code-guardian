/*
Copyright © 2025 Guilherme Soares <GUILHERMELUZSOARES@GMAIL.COM>
*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v49/github"
	"github.com/gsoares85/code-guardian/internal/github_internal"
	"github.com/gsoares85/code-guardian/internal/openai"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// prReviewCmd represents the prReview command
var prReviewCmd = &cobra.Command{
	Use:   "pr-review [owner] [repo] [pr-number] [flags]",
	Short: "Analysis of a pull request",
	Long:  `Use this command to perform a complete code review of a pull request.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		owner, repo := args[0], args[1]
		prNumber, err := strconv.Atoi(args[2])
		if err != nil {
			color.Red("Error: Invalid PR number")
			return
		}
		saveOutput, _ := cmd.Flags().GetBool("output")

		// Fetch PR Data
		pr, files, prDiff, aiFeedback, err := analyzePullRequest(owner, repo, prNumber)
		if err != nil {
			color.Red("Error analyzing PR: %s\n", err)
			return
		}

		// Display results
		displayAnalysis(pr, files, prDiff, aiFeedback)

		// Save to file if `--output` flag is enabled
		if saveOutput {
			saveReport(pr, repo, owner, prNumber, files, prDiff, aiFeedback)
		}
	},
}

func analyzePullRequest(owner, repo string, prNumber int) (*github.PullRequest, []string, string, string, error) {
	pr, err := github_internal.GetPullRequest(owner, repo, prNumber)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("error fetching PR: %w", err)
	}

	files, err := github_internal.GetPullRequestFiles(owner, repo, prNumber)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("error fetching PR files: %w", err)
	}

	prDiff, err := github_internal.GetPullRequestDiff(owner, repo, prNumber)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("error fetching PR diff: %w", err)
	}

	color.Blue("Analyzing PR...")
	aiFeedback, err := openai.AnalyzePRWithAI(prDiff)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("error analyzing PR with AI: %w", err)
	}

	return pr, files, prDiff, aiFeedback, nil
}

func displayAnalysis(pr *github.PullRequest, files []string, prDiff string, aiFeedback string) {
	color.Cyan("\n🔍 PR %d - %s\n", pr.GetNumber(), pr.GetTitle())
	color.Green("👤 Author: %s\n", pr.GetUser().GetLogin())
	color.Yellow("📅 Created at: %s\n", pr.GetCreatedAt())
	color.Blue("🔗 PR Link: %s\n", pr.GetHTMLURL())

	color.Magenta("\n📂 Files changed (%d):\n", len(files))
	for _, file := range files {
		fmt.Println("  -", file)
	}

	color.Magenta("\n📌 Code changes (%d characters):\n", len(prDiff))
	color.Green("\n📢 AI Feedback:\n")
	fmt.Println(aiFeedback)
}

func generateMarkdownReport(pr *github.PullRequest, files []string, prDiff string, aiFeedback string) string {
	return fmt.Sprintf(`# Pull Request Analysis Report

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
`,
		pr.GetNumber(), pr.GetTitle(), pr.GetUser().GetLogin(), pr.GetCreatedAt(),
		pr.GetHTMLURL(), pr.GetHTMLURL(),
		formatFileList(files),
		prDiff,
		aiFeedback,
	)
}

func saveReport(pr *github.PullRequest, repo, owner string, prNumber int, files []string, prDiff string, aiFeedback string) {
	timestamp := time.Now().Format("20060102-150405")
	outputFile := fmt.Sprintf("%s-%s_%s_%d.md", timestamp, repo, owner, prNumber)
	outputPath := filepath.Join("reports", "pr", outputFile)

	content := generateMarkdownReport(pr, files, prDiff, aiFeedback)

	if err := saveAnalysisToFile(outputPath, content); err != nil {
		color.Red("❌ ERROR saving analysis: %s\n", err)
		return
	}
	color.Green("✅ Analysis saved to file: %s\n", outputPath)
}

func formatFileList(files []string) string {
	formated := ""
	for _, file := range files {
		formated += fmt.Sprintf("- %s\n", file)
	}
	return formated
}

func saveAnalysisToFile(fileName, content string) error {
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(fileName)
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
