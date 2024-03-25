package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var matchStyle = lipgloss.NewStyle().
	Foreground(log.ErrorLevelStyle.GetForeground())

func PathMatch(path string, rules []Regexp) (*string, int) {
	for i, rule := range rules {
		match := rule.FindStringIndex(path)
		if match == nil {
			continue
		}

		left := path[:match[0]]
		middle := matchStyle.Render(path[match[0]:match[1]])
		right := path[match[1]:]

		builder := new(strings.Builder)
		builder.Grow(len(left) + len(middle) + len(right))
		_, err := builder.WriteString(left)
		check(err)
		_, err = builder.WriteString(middle)
		check(err)
		_, err = builder.WriteString(right)
		check(err)

		renderedString := builder.String()
		return &renderedString, i
	}

	return nil, -1
}
