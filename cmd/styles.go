package cmd

import (
	"os"

	"github.com/charmbracelet/glamour"
)

func renderMarkdown(text string) string {
	if !isTerminal() {
		return text
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return text
	}

	out, err := renderer.Render(text)
	if err != nil {
		return text
	}

	return out
}

func isTerminal() bool {
	stat, _ := os.Stdout.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}
