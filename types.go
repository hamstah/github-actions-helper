package main

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v29/github"
)

type Repo struct {
	OwnerName *string
	RepoName  *string
}

type Issue struct {
	Repo
	ID *int
}

func ParseRepoName(repo *github.Repository) (*Repo, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository is nil")
	}
	parts := strings.Split(*repo.FullName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repo full_name %s", *repo.FullName)
	}
	return &Repo{
		OwnerName: &parts[0],
		RepoName:  &parts[1],
	}, nil
}

func RepoFromEvent(event interface{}) (*Issue, error) {
	switch event.(type) {
	case *github.IssueCommentEvent:
		comment := event.(*github.IssueCommentEvent)
		issue := comment.GetIssue()
		repo, err := ParseRepoName(comment.GetRepo())
		if err != nil {
			return nil, err
		}

		return &Issue{
			ID:   issue.Number,
			Repo: *repo,
		}, nil
	case *github.PullRequestEvent:
		pr := event.(*github.PullRequestEvent)

		repo, err := ParseRepoName(pr.GetPullRequest().GetBase().Repo)
		if err != nil {
			return nil, err
		}

		return &Issue{
			Repo: *repo,
			ID:   pr.Number,
		}, nil
	default:
		return nil, nil
	}
}
