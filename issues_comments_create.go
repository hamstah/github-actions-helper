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

type Repo struct {
	OwnerName *string
	RepoName  *string
}

func RepoFlags(cmd *kingpin.CmdClause) Repo {
	return Repo{
		OwnerName: cmd.Flag("owner", "Owner name").String(),
		RepoName:  cmd.Flag("repo", "Repo name").String(),
	}
}

type Issue struct {
	Repo
	ID *int
}

func IssueFlags(cmd *kingpin.CmdClause) Issue {
	return Issue{
		Repo: RepoFlags(cmd),
		ID:   cmd.Flag("id", "PR id").Int(),
	}
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

	var owner, repo string
	var id int

	switch event.(type) {
	case *github.PullRequestEvent:
		pr := event.(*github.PullRequestEvent)
		fullName := *pr.GetPullRequest().GetBase().Repo.FullName
		parts := strings.Split(fullName, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid repo full_name %s", fullName)
		}
		owner = parts[0]
		repo = parts[1]
		id = *pr.Number
	default:
		owner = *issuesCommentsCreateFlags.OwnerName
		repo = *issuesCommentsCreateFlags.RepoName
		id = *issuesCommentsCreateFlags.ID
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
		owner,
		repo,
		id,
		comment,
	)
}

func HandleIssuesCommentsCreate(client *github.Client, owner, repo string, id int, comment string) (interface{}, error) {
	if owner == "" {
		return nil, fmt.Errorf("owner can't be empty")
	}

	if repo == "" {
		return nil, fmt.Errorf("repo can't be empty")
	}

	if id == 0 {
		return nil, fmt.Errorf("id can't be empty")
	}

	input := &github.IssueComment{Body: github.String(comment)}

	c, _, err := client.Issues.CreateComment(context.Background(), owner, repo, id, input)
	if err != nil {
		return nil, fmt.Errorf("Issues.CreateComment returned error: %v", err)
	}

	return c, nil
}
