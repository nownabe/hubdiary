package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

type githubRepo struct {
	repositories *github.RepositoriesService
	*config
}

func newGithubRepo(ctx context.Context, cfg *config) *githubRepo {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.PAT},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := github.NewClient(tc)

	return &githubRepo{
		repositories: c.Repositories,
		config:       cfg,
	}
}

func (r *githubRepo) ReadContent(ctx context.Context, path string) (string, string, error) {
	opts := &github.RepositoryContentGetOptions{Ref: r.Branch}

	repoContent, _, _, err := r.repositories.GetContents(ctx, r.Owner, r.Repo, path, opts)
	if err != nil {
		return "", "", fmt.Errorf("failed to get content: %w", err)
	}

	content, err := repoContent.GetContent()
	if err != nil {
		return "", "", fmt.Errorf("failed to decode content: %w", err)
	}

	return content, repoContent.GetSHA(), nil
}

func (r *githubRepo) WriteContent(ctx context.Context, path, content, sha string) error {
	committer := &github.CommitAuthor{
		Name:  github.String(r.CommitterName),
		Email: github.String(r.CommitterEmail),
	}

	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(""),
		Content:   []byte(content),
		Branch:    github.String(r.Branch),
		Committer: committer,
	}

	now := time.Now().Format(time.RFC3339)

	if sha == "" {
		opts.Message = github.String("Created by hubdiary at " + now)
		_, _, err := r.repositories.CreateFile(ctx, r.Owner, r.Repo, path, opts)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	} else {
		opts.Message = github.String("Modified by hubdiary at " + now)
		opts.SHA = github.String(sha)
		_, _, err := r.repositories.UpdateFile(ctx, r.Owner, r.Repo, path, opts)
		if err != nil {
			return fmt.Errorf("failed to update file: %w", err)
		}
	}

	return nil
}
