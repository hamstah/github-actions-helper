package main

import (
	"github.com/google/go-github/v29/github"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	pullsMerge      = kingpin.Command("pulls-merge", "Merge a PR")
	pullsMergeFlags = PullsMergeFlags(pullsMerge)
)

type PullsMerge struct {
	Issue
	CommitMessage *string
	CommitTitle   *string
	MergeMethod   *string
	SHA           *string
}

func PullsMergeFlags(cmd *kingpin.CmdClause) PullsMerge {
	return PullsMerge{
		Issue:         IssueFlags(cmd),
		CommitMessage: cmd.Flag("commit-message", "Commit message").Required().String(),
		CommitTitle:   cmd.Flag("commit-title", "Commit title").String(),
		MergeMethod:   cmd.Flag("merge-method", "Merge method").Default("merge").Enum("merge", "rebase", "squash"),
		SHA:           cmd.Flag("sha", "SHA of the commit to merge").String(),
		DeleteBranch:  cmd.Flag("delete-branch", "Delete branch after merging").Default("false").Bool(),
	}
}

func HandlePullsMergeCmd(ctx context.Context, client *github.Client, event interface{}) (interface{}, error) {
	issue, err := RepoFromEvent(event)
	if err != nil {
		return nil, err
	}

	if issue == nil {
		issue = &issuesCommentsCreateFlags.Issue
	}

	options := &github.PullRequestOptions{
		CommitTitle: *pullsMergeFlags.CommitTitle,
		MergeMethod: *pullsMergeFlags.MergeMethod,
		SHA:         *pullsMergeFlags.SHA,
	}

	c, _, err := client.PullRequests.Merge(
		ctx,
		*issue.OwnerName,
		*issue.RepoName,
		*issue.ID,
		*pullsMergeFlags.CommitMessage,
		options,
	)

	return c, err
}
