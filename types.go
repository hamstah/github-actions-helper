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

func RepoFromEvent(event interface{}) (*Issue, error) {
	switch event.(type) {
	case *github.PullRequestEvent:
		pr := event.(*github.PullRequestEvent)
		fullName := *pr.GetPullRequest().GetBase().Repo.FullName
		parts := strings.Split(fullName, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid repo full_name %s", fullName)
		}
		return &Issue{
			Repo: Repo{
				OwnerName: &parts[0],
				RepoName:  &parts[1],
			},
			ID: pr.Number,
		}, nil
	default:
		return nil, nil
	}
}
