package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/motemen/go-gitconfig"
)

const ()

func main() {
	cfg := parseConfig()

	ctx := context.Background()
	r := newGithubRepo(ctx, cfg)

	var date time.Time

	if len(os.Args) == 2 {
		var err error
		date, err = time.Parse("2006-01-02", os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
		date = time.Now()
	}

	path := date.Format("2006/01/02.md")

	content, sha, err := r.ReadContent(ctx, path)
	if err != nil {
		msg := err.Error()
		if msg[len(msg)-16:] != "404 Not Found []" {
			panic(err)
		}

		content = fmt.Sprintf("# %s\n\n* \n", date.Format("2006-01-02"))
	}

	e := &editor{editor: cfg.Editor}
	content, err = e.Edit(content)
	if err != nil {
		panic(err)
	}

	if err := r.WriteContent(ctx, path, content, sha); err != nil {
		panic(err)
	}
}

func parseConfig() *config {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`hubdiary

Usage:
  %s [2006-01-02]

`,
			os.Args[0],
		)

		flag.PrintDefaults()
	}

	configPath := flag.String("config", "",
		`Path to config file. Default path is $XDG_CONFIG_HOME/hubdiary/config.json.

This is an example:

{
	"owner": "owner",
	"repo": "diaryrepo",
	"branch": "main",
	"committer_name": "owner",
	"committer_email": "owner@users.noreply.github.com",
	"pat": "xxx",
	"editor": "vim"
}

"committer_name" and "committer_email" is used for making git commits.
If you specify these values as command line arguments, values from the
config file will be overwritten. To see details of each option,
run "hubdiary -h".`)

	owner := flag.String("owner", "", "Owner of repository to store diary in. Default comes from user.name of git config.")
	repo := flag.String("repo", "", "Repository to store diary in. Default is 'diary'.")
	branch := flag.String("branch", "", "Branch name to store diary in. Default is 'main'.")
	committerName := flag.String("committerName", "", "Commit author name. Default comes from user.name of git config.")
	committerEmail := flag.String("committerEmail", "", "Commit author email. Default comes from user.email of git config.")
	pat := flag.String("pat", "", fmt.Sprintf("GitHub Personal Access Token. Default comes from $%s.", patEnvName))
	editor := flag.String("editor", "", fmt.Sprintf("Editor path. Default comes from $%s", editorEnvName))

	flag.Parse()

	cfg, err := newConfig(gitconfig.Default)
	if err != nil {
		panic(err)
	}

	if err := cfg.loadFile(*configPath); err != nil {
		panic(err)
	}

	if *owner != "" {
		cfg.Owner = *owner
	}
	if *repo != "" {
		cfg.Repo = *repo
	}
	if *branch != "" {
		cfg.Branch = *branch
	}
	if *committerName != "" {
		cfg.CommitterName = *committerName
	}
	if *committerEmail != "" {
		cfg.CommitterEmail = *committerEmail
	}
	if *pat != "" {
		cfg.PAT = *pat
	}
	if *editor != "" {
		cfg.Editor = *editor
	}

	return cfg
}
