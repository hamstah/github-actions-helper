package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/v29/github"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	quiet = kingpin.Flag("quiet", "Do not show any output").Default("false").Bool()
)

func ParseEvent() (interface{}, error) {
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	if eventPath == "" {
		return nil, nil
	}

	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if eventName == "" {
		return nil, nil
	}

	bytes, err := ioutil.ReadFile(eventPath)
	if err != nil {
		return nil, err
	}

	var event interface{}
	switch eventName {
	case "push":
		event = &github.PushEvent{}
	case "pull_request":
		event = &github.PullRequestEvent{}
	case "issue_comment":
		event = &github.IssueCommentEvent{}
	default:
		if !*quiet {
			fmt.Println(eventName)
			fmt.Println(string(bytes))
		}
	}

	if event == nil {
		return nil, fmt.Errorf("unsupported event name %s", eventName)
	}

	err = json.Unmarshal(bytes, event)
	return event, err
}

func FatalOnError(err error) {
	if err != nil {
		if !*quiet {
			log.Fatalln(err)
		} else {
			os.Exit(1)
		}
	}
}

func main() {
	var result interface{}
	var err error

	client := NewClient(os.Getenv("GITHUB_TOKEN"))
	event, err := ParseEvent()
	FatalOnError(err)

	ctx := context.Background()

	switch kingpin.Parse() {
	case issuesCommentsCreate.FullCommand():
		result, err = HandleIssuesCommentsCreateCmd(ctx, client, event)
	case issuesCommentsFilter.FullCommand():
		result, err = HandleIssuesCommentsFilterCmd(ctx, client, event)
	case pullsGet.FullCommand():
		result, err = HandlePullsGetCmd(ctx, client, event)
	case pullsMerge.FullCommand():
		result, err = HandlePullsMergeCmd(ctx, client, event)
	}
	FatalOnError(err)

	if result != nil && !*quiet {
		bytes, err := json.MarshalIndent(result, "", "  ")
		FatalOnError(err)
		fmt.Println(string(bytes))
	}
}
