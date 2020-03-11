package main

import "gopkg.in/alecthomas/kingpin.v2"

func RepoFlags(cmd *kingpin.CmdClause) Repo {
	return Repo{
		OwnerName: cmd.Flag("owner", "Owner name").String(),
		RepoName:  cmd.Flag("repo", "Repo name").String(),
	}
}

func IssueFlags(cmd *kingpin.CmdClause) Issue {
	return Issue{
		Repo: RepoFlags(cmd),
		ID:   cmd.Flag("id", "PR id").Int(),
	}
}
