package tools

import "sync"

var (
	statsMu       sync.Mutex
	toolCallStats = map[string]int{}
	sessionCalls  int
)

func recordToolCall(name string) {
	statsMu.Lock()
	defer statsMu.Unlock()
	toolCallStats[name]++
	sessionCalls++
}

func sessionStatsText() string {
	statsMu.Lock()
	defer statsMu.Unlock()
	lines := []string{"Tool calls this session:"}
	for name, n := range toolCallStats {
		lines = append(lines, "- "+name+" tool: called "+itoa(n)+" times")
	}
	return joinLines(lines)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

func joinLines(lines []string) string {
	out := ""
	for i, l := range lines {
		if i > 0 {
			out += "\n"
		}
		out += l
	}
	return out
}
