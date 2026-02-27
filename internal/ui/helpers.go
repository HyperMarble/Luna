package ui

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func takeLast(items []string, n int) []string {
	if n <= 0 || len(items) == 0 {
		return []string{}
	}
	if len(items) <= n {
		return items
	}
	return items[len(items)-n:]
}
