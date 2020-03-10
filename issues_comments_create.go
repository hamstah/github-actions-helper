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
	Comment        *string
	Markdown       *bool
	MarkdownSyntax *string
}

func IssuesCommentsCreateFlags(cmd *kingpin.CmdClause) IssuesCommentsCreate {
	return IssuesCommentsCreate{
		Issue:          IssueFlags(cmd),
		Comment:        cmd.Flag("comment", "Comment to add, use file:// prefix to load a file").String(),
		Markdown:       cmd.Flag("markdown", "Format the comment with markdown").Default("false").Bool(),
		MarkdownSyntax: cmd.Flag("markdown-syntax", "Syntax highlighting to use in markdown").String(),
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

	if *issuesCommentsCreateFlags.Markdown || *issuesCommentsCreateFlags.MarkdownSyntax != "" {
		comment = fmt.Sprintf("```%s\n%s\n```", *issuesCommentsCreateFlags.MarkdownSyntax, comment)
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
