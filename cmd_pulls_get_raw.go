package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/go-github/v29/github"
	"golang.org/x/net/context"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	pullsGetRaw      = kingpin.Command("pulls-get-raw", "Get a PR in raw format")
	pullsGetRawFlags = PullsGetRawFlags(pullsGetRaw)
)

type PullsGetRaw struct {
	Issue
	Format   *string
	NameOnly *bool
}

func PullsGetRawFlags(cmd *kingpin.CmdClause) PullsGetRaw {
	return PullsGetRaw{
		Issue:    IssueFlags(cmd),
		Format:   cmd.Flag("format", "Output format").Required().Enum("diff", "patch"),
		NameOnly: cmd.Flag("name-only", "Only show changed filenames").Default("false").Bool(),
	}
}

func HandlePullsGetRawCmd(ctx context.Context, client *github.Client, event interface{}) (interface{}, error) {
	issue, err := RepoFromEvent(event)
	if err != nil {
		return nil, err
	}

	if issue == nil {
		issue = &pullsGetRawFlags.Issue
	}

	options := github.RawOptions{}

	switch *pullsGetRawFlags.Format {
	case "diff":
		options.Type = github.Diff
	case "patch":
		options.Type = github.Patch
	default:
		return nil, fmt.Errorf("invalid format value %s", *pullsGetRawFlags.Format)
	}
	raw, _, err := client.PullRequests.GetRaw(ctx, *issue.OwnerName, *issue.RepoName, *issue.ID, options)

	if *pullsGetRawFlags.NameOnly {
		files := map[string]interface{}{}
		for _, line := range strings.Split(raw, "\n") {
			if strings.HasPrefix(line, "+++ ") || strings.HasPrefix(line, "--- ") {
				line = line[4:]
				if line == "/dev/null" {
					continue
				}

				if strings.HasPrefix(line, "a/") || strings.HasPrefix(line, "b/") {
					files[line[2:]] = true
				}
			}
		}

		output := make([]string, len(files))
		index := 0
		for key, _ := range files {
			output[index] = key
			index += 1
		}

		sort.Strings(output)
		raw = strings.Join(output, "\n")
	}

	return raw, err
}
