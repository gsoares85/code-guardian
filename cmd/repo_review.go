package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gsoares85/code-guardian/internal/github_internal"
	"github.com/gsoares85/code-guardian/internal/openai"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
	"time"
)

var repoReviewCmd = &cobra.Command{
	Use:   "repo-review [owner] [repo] [flags]",
	Short: "Analyse an entire github repository",
	Long: `This command scans all source files recursively in a repository
           to train the AI about the application and provide:
           - A summary of what the application does
           - Key use cases
           - A high-level code quality review (only critical issues)
           - A high-level security review (only critical issues)
           - Key improvement areas`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		owner, repo := args[0], args[1]
		saveOutput, err := cmd.Flags().GetBool("output")
		if err != nil {
			color.Red("Error: %v", err)
			return
		}

		color.Blue("\n🔍 Fetching all source code files recursively...\n")
		files, err := github_internal.GetRepositoryFilesRecursive(owner, repo)
		if err != nil {
			color.Red("❌ ERROR: Fetching repository files: %s\n", err)
			return
		}

		color.Green("📂 Repository contains %d files\n", len(files))

		sourceCode := fetchAllSourceCode(owner, repo, files)
		if len(sourceCode) == 0 {
			color.Red("❌ ERROR: No valid source code found for analysis")
			return
		}

		color.Blue("\n🤖 Training AI with full source code...\n")
		summary, useCases, codeReview, securityReview, improvements := analyzeRepositoryWithAI(sourceCode)

		displayRepoAnalysis(repo, summary, useCases, codeReview, securityReview, improvements)

		if saveOutput {
			saveRepoAnalysis(repo, owner, summary, useCases, codeReview, securityReview, improvements)
		}
	},
}

func fetchAllSourceCode(owner, repo string, files []string) string {
	var allCode strings.Builder

	for _, file := range files {
		if strings.HasSuffix(file, ".md") || strings.Contains(file, "LICENSE") {
			continue
		}

		color.Cyan("\n📄 Reading file: %s", file)

		content, err := github_internal.GetFileContent(owner, repo, file)
		if err != nil {
			color.Red("❌ ERROR fetching file content: %s\n", err)
			continue
		}

		allCode.WriteString(fmt.Sprintf("\n// File: %s\n%s\n", file, content))
	}

	return allCode.String()
}

func analyzeRepositoryWithAI(sourceCode string) (string, string, string, string, string) {
	summaryPrompt := "Analyze this entire source codebase and provide a concise summary of what the application does."
	useCasesPrompt := "Extract the most important use cases from the source code."
	codeQualityPrompt := "Identify the most critical code quality issues found in the source code. Provide a brief list."
	securityPrompt := "Identify the most critical security vulnerabilities in the source code. Provide a brief list."
	improvementPrompt := "Suggest the most important areas to improve in the application."

	summary, _ := openai.AnalyzeCodeWithAI(sourceCode, summaryPrompt)
	useCases, _ := openai.AnalyzeCodeWithAI(sourceCode, useCasesPrompt)
	codeReview, _ := openai.AnalyzeCodeWithAI(sourceCode, codeQualityPrompt)
	securityReview, _ := openai.AnalyzeCodeWithAI(sourceCode, securityPrompt)
	improvements, _ := openai.AnalyzeCodeWithAI(sourceCode, improvementPrompt)

	return summary, useCases, codeReview, securityReview, improvements
}

func displayRepoAnalysis(repo, summary, useCases, codeReview, securityReview, improvements string) {
	color.Magenta("\n📌 Repository Analysis Summary for %s\n", repo)
	color.Cyan("\n📖 Application Summary:\n")
	fmt.Println(summary)

	color.Green("\n✅ Key Use Cases:\n")
	fmt.Println(useCases)

	color.Red("\n🚨 Code Quality Issues (Critical Only):\n")
	fmt.Println(codeReview)

	color.Yellow("\n🔒 Security Issues (Critical Only):\n")
	fmt.Println(securityReview)

	color.Blue("\n📈 Key Areas for Improvement:\n")
	fmt.Println(improvements)
}

func saveRepoAnalysis(repo, owner, summary, useCases, codeReview, securityReview, improvements string) {
	timestamp := time.Now().Format("20060102-150405")
	outputFile := fmt.Sprintf("%s-%s_%s.md", timestamp, repo, owner)
	outputPath := filepath.Join("reports", "repo", outputFile)

	content := generateRepoMarkdown(repo, summary, useCases, codeReview, securityReview, improvements)

	if err := saveAnalysisToFile(outputPath, content); err != nil {
		color.Red("❌ ERROR saving analysis: %s\n", err)
		return
	}
	color.Green("✅ Repository analysis saved to file: %s\n", outputPath)
}

func generateRepoMarkdown(repo, summary, useCases, codeReview, securityReview, improvements string) string {
	return fmt.Sprintf(`# Repository Analysis Report

## 📂 Repository: %s

## 📖 Application Summary:
%s

## ✅ Key Use Cases:
%s

## 🚨 Code Quality Issues (Critical Only):
%s

## 🔒 Security Issues (Critical Only):
%s

## 📈 Key Areas for Improvement:
%s
`, repo, summary, useCases, codeReview, securityReview, improvements)
}

func init() {
	rootCmd.AddCommand(repoReviewCmd)
	repoReviewCmd.Flags().BoolP("output", "o", false, "Save analysis to a Markdown file in ./reports/repo/")
}
