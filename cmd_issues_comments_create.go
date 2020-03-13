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

func HandleIssuesCommentsCreateCmd(ctx context.Context, client *github.Client, event interface{}) (interface{}, error) {

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

	if *issuesCommentsCreateFlags.Markdown {
		comment = FormatComment(comment)
	}
	return HandleIssuesCommentsCreate(
		ctx,
		client,
		issue,
		comment,
	)
}

func HandleIssuesCommentsCreate(ctx context.Context, client *github.Client, issue *Issue, comment string) (interface{}, error) {
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

	c, _, err := client.Issues.CreateComment(ctx, *issue.OwnerName, *issue.RepoName, *issue.ID, input)
	return c, err
}
