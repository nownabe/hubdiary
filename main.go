package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/motemen/go-gitconfig"
)

func main() {
	var (
		configPath string
		repo       string
		user       string
		email      string
		pat        string
	)

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

	flag.StringVar(&configPath, "config", "",
		`Path to config file. Default path is $XDG_CONFIG_HOME/hubdiary/config.json.

This is an example:

{
	"repo": "owner/repo",
	"user": "owner",
	"email": "owner@users.noreply.github.com",
	"pat": "xxx"
}

"user" and "email" is used for making git commits. If you specify these
values as environment variables or command line arguments, values from
the config file will be overwritten. To see details of each option,
run "hubdiary -h".`)

	flag.StringVar(&repo, "repo", "", "Repository to store diary in. Default is ${user}/diary.")
	flag.StringVar(&user, "user", "", "Commit author name. Default is user.name of .gitconfig.")
	flag.StringVar(&email, "email", "", "Commit author email. Default is user.email of .gitconfig.")
	flag.StringVar(&pat, "pat", "",
		fmt.Sprintf("GitHub Personal Access Token. You can specify it with %s environment variable.", patEnvName))

	flag.Parse()

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

	ctx := context.Background()

	cl := &configLoader{
		gitConfig: gitconfig.Default,
		path:      configPath,
		envPAT:    os.Getenv(patEnvName),
	}
	cfg, err := cl.Load(repo, user, email, pat)

	r := newGithubRepo(ctx, cfg)

	path := date.Format("2006/01/02.md")

	content, sha, err := r.ReadContent(ctx, path)
	if err != nil {
		msg := err.Error()
		if msg[len(msg)-16:] != "404 Not Found []" {
			panic(err)
		}

		content = fmt.Sprintf("# %s\n\n* \n", date.Format("2006-01-02"))
	}

	e := &editor{editor: "nvim"} // TODO: Move to config
	content, err = e.Edit(content)
	if err != nil {
		panic(err)
	}

	if err := r.WriteContent(ctx, path, content, sha); err != nil {
		panic(err)
	}
}
