package openai

// Mock function to simulate AI-generated feedback
func MockAnalyzePRWithAI(diff string) (string, error) {
	return "- ✅ Correctly added memory release (`free(ptr);`).\n- ⚠️ Consider checking if `malloc` returned a valid pointer before using it.", nil
}
