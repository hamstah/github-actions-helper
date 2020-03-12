package main

import (
	"github.com/google/go-github/v29/github"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	pullsGet      = kingpin.Command("pulls-get", "Get a PR")
	pullsGetFlags = PullsGetFlags(pullsGet)
)

type PullsGet struct {
	Issue
}

func PullsGetFlags(cmd *kingpin.CmdClause) PullsGet {
	return PullsGet{
		Issue: IssueFlags(cmd),
	}
}

func HandlePullsGetCmd(ctx context.Context, client *github.Client, event interface{}) (interface{}, error) {
	issue, err := RepoFromEvent(event)
	if err != nil {
		return nil, err
	}

	if issue == nil {
		issue = &pullsGetFlags.Issue
	}

	var pr *github.PullRequest
	switch event.(type) {
	case *github.PullRequestEvent:
		pr = event.(*github.PullRequestEvent).GetPullRequest()
	default:
		pr, _, err = client.PullRequests.Get(ctx, *issue.OwnerName, *issue.RepoName, *issue.ID)
		if err != nil {
			return nil, err
		}
	}

	return pr, err
}
