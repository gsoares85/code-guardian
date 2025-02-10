package mocks

func MockAnalyzePRWithAI(diff string) (string, error) {
	return "- ✅ Correctly added memory release (`free(ptr);`).\n- ⚠️ Consider checking if `malloc` returned a valid pointer before using it.", nil
}

func MockAnalyzeCodeWithAI(code string, prompt string) (string, error) {
	return "This application is a simple web server that handles HTTP requests.", nil
}
