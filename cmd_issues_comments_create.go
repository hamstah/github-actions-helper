package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/v29/github"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	issuesCommentsCreate      = kingpin.Command("issues-comments-create", "Create a comment on a PR")
	issuesCommentsCreateFlags = IssuesCommentsCreateFlags(issuesCommentsCreate)
)

type Section struct {
	Title     string
	Content   []string
	Collapsed bool
}

type IssuesCommentsCreate struct {
	Issue
	Comment  *string
	Markdown *bool
}

func IssuesCommentsCreateFlags(cmd *kingpin.CmdClause) IssuesCommentsCreate {
	return IssuesCommentsCreate{
		Issue:    IssueFlags(cmd),
		Comment:  cmd.Flag("comment", "Comment to add, use file:// prefix to load a file").String(),
		Markdown: cmd.Flag("markdown", "Format the comment with markdown").Default("false").Bool(),
	}
}

func HandleIssuesCommentsCreateCmd(client *github.Client, event interface{}) (interface{}, error) {

	issue, err := RepoFromEvent(event)
	if err != nil {
		return nil, err
	}

	if issue == nil {
		issue = &issuesCommentsCreateFlags.Issue
	}

	commentArg := *issuesCommentsCreateFlags.Comment
	var comment string
	if commentArg == "" {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		comment = string(data)
		if len(comment) == 0 {
			return nil, fmt.Errorf("could not read comment from stdin")
		}
	} else if strings.HasPrefix(commentArg, "file://") {
		filename := commentArg[7:]
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		comment = string(data)
	} else {
		comment = commentArg
	}

	sections := []Section{Section{}}

	if *issuesCommentsCreateFlags.Markdown {
		lines := strings.Split(comment, "\n")

		for _, line := range lines {
			if strings.HasPrefix(line, "::") {
				if len(sections[len(sections)-1].Content) != 0 {
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
		final := make([]string, len(sections))
		for index, section := range sections {
			content := strings.Join(section.Content, "\n")
			if section.Collapsed {
				content = fmt.Sprintf("<details><summary>%s</summary>\n\n```\n%s```\n</details>\n", section.Title, content)
			} else {
				content = fmt.Sprintf("%s\n\n```\n%s```\n", section.Title, content)
			}
			final[index] = content
			comment = strings.Join(final, "\n")
		}
	}
	return HandleIssuesCommentsCreate(
		client,
		issue,
		comment,
	)
}

func HandleIssuesCommentsCreate(client *github.Client, issue *Issue, comment string) (interface{}, error) {
	if *issue.OwnerName == "" {
		return nil, fmt.Errorf("owner can't be empty")
	}

	if *issue.RepoName == "" {
		return nil, fmt.Errorf("repo can't be empty")
	}

	if *issue.ID == 0 {
		return nil, fmt.Errorf("id can't be empty")
	}

	input := &github.IssueComment{Body: github.String(comment)}

	c, _, err := client.Issues.CreateComment(context.Background(), *issue.OwnerName, *issue.RepoName, *issue.ID, input)
	return c, err
}
