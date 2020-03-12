package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/go-github/v29/github"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	issuesCommentsFilter      = kingpin.Command("issues-comments-filter", "Filter a comment on a PR")
	issuesCommentsFilterFlags = IssuesCommentsFilterFlags(issuesCommentsFilter)
)

type IssuesCommentsFilter struct {
	Issue
	Regex            *string
	State            *[]string
	Action           *[]string
	RequireMergeable *bool
}

func IssuesCommentsFilterFlags(cmd *kingpin.CmdClause) IssuesCommentsFilter {
	return IssuesCommentsFilter{
		Issue:            IssueFlags(cmd),
		Regex:            cmd.Flag("regex", "Regex to use for filtering").Required().String(),
		State:            cmd.Flag("state", "State of the issue").Default("open").Strings(),
		Action:           cmd.Flag("action", "Action done on the comment").Default("created").Strings(),
		RequireMergeable: cmd.Flag("require-mergeable", "Require the PR to be mergeable").Default("true").Bool(),
	}
}

func Contains(hay *[]string, needle *string) bool {
	if hay == nil {
		return false
	}

	for _, e := range *hay {
		if e == *needle {
			return true
		}
	}
	return false
}

func HandleIssuesCommentsFilterCmd(ctx context.Context, client *github.Client, event interface{}) (interface{}, error) {
	issue, err := RepoFromEvent(event)
	if err != nil {
		return nil, err
	}

	if issue == nil {
		issue = &issuesCommentsFilterFlags.Issue
	}

	issueComment, ok := event.(*github.IssueCommentEvent)
	if !ok {
		return nil, fmt.Errorf("event is not issue_comment")
	}

	if !Contains(issuesCommentsFilterFlags.Action, issueComment.Action) {
		return nil, fmt.Errorf("comment action is not valid: %s", *issueComment.Action)
	}

	issueState := issueComment.GetIssue().State
	if !Contains(issuesCommentsFilterFlags.State, issueState) {
		return nil, fmt.Errorf("issue state is not valid: %s", *issueState)
	}

	if *issuesCommentsFilterFlags.RequireMergeable {
		links := issueComment.GetIssue().GetPullRequestLinks()
		if links != nil {
			pr, _, err := client.PullRequests.Get(ctx, *issue.OwnerName, *issue.RepoName, *issue.ID)
			if err != nil {
				return nil, err
			}

			if pr.GetMergeableState() != "clean" {
				return nil, fmt.Errorf("pull request is not mergeable")
			}
		}
	}

	regex, err := regexp.Compile(*issuesCommentsFilterFlags.Regex)
	if err != nil {
		return nil, err
	}

	matches := regex.FindStringSubmatch(*issueComment.GetComment().Body)

	if len(matches) == 0 {
		return nil, fmt.Errorf("comment does not match the regex")
	}

	result := map[string]string{}
	names := regex.SubexpNames()

	for index, match := range matches {
		if len(names) > index && names[index] != "" {
			result[names[index]] = match
		} else {
			result[strconv.Itoa(index)] = match
		}
	}

	return result, nil
}
