package cmd

import _ "embed"

//go:embed help/root.md
var helpRoot string

//go:embed help/topics.md
var helpTopics string

//go:embed help/version.md
var helpVersion string

//go:embed help/completion.md
var helpCompletion string

//go:embed help/systems.md
var helpSystems string

func GetHelp(topic string) string {
	switch topic {
	case "root":
		return renderMarkdown(helpRoot)
	case "topics":
		return renderMarkdown(helpTopics)
	case "version":
		return renderMarkdown(helpVersion)
	case "completion":
		return renderMarkdown(helpCompletion)
	case "systems":
		return renderMarkdown(helpSystems)
	default:
		return ""
	}
}
