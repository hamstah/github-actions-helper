package main

import (
	"fmt"
	"strings"
)

type Section struct {
	Title     string
	Content   []string
	Collapsed bool
}

func ParseComment(comment string) []Section {
	sections := []Section{Section{}}

	lines := strings.Split(strings.TrimSpace(comment), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "::") {
			previousSection := sections[len(sections)-1]
			if len(previousSection.Content) != 0 || previousSection.Title != "" {
				sections = append(sections, Section{})
			}
			line = line[2:]

			section := &sections[len(sections)-1]
			if strings.HasPrefix(line, "-") {
				section.Collapsed = true
				line = line[1:]
			}
			section.Title = strings.TrimSpace(line)
			continue
		}

		section := &sections[len(sections)-1]
		section.Content = append(section.Content, line)
	}
	return sections
}

func FormatComment(comment string) string {
	sections := ParseComment(comment)
	final := make([]string, len(sections))
	for index, section := range sections {
		content := strings.Join(section.Content, "\n")
		if section.Collapsed {
			content = fmt.Sprintf("<details><summary>%s</summary>\n\n```\n%s```\n</details>\n", section.Title, content)
		} else {
			content = fmt.Sprintf("%s\n\n```\n%s```\n", section.Title, content)
		}
		final[index] = content
	}
	return strings.Join(final, "\n")
}
