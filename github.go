package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

type githubRepo struct {
	repositories   *github.RepositoriesService
	owner          string
	repo           string
	branch         string
	committerName  string
	committerEmail string
}

func newGithubRepo(ctx context.Context, cfg *config) *githubRepo {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.PAT},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := github.NewClient(tc)

	ownerAndRepo := strings.Split(cfg.Repo, "/")

	return &githubRepo{
		repositories:   c.Repositories,
		owner:          ownerAndRepo[0],
		repo:           ownerAndRepo[1],
		branch:         "main", // TODO: Move to config
		committerName:  cfg.User,
		committerEmail: cfg.Email,
	}
}

func (r *githubRepo) ReadContent(ctx context.Context, path string) (string, string, error) {
	opts := &github.RepositoryContentGetOptions{Ref: r.branch}

	repoContent, _, _, err := r.repositories.GetContents(ctx, r.owner, r.repo, path, opts)
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
		Name:  github.String(r.committerName),
		Email: github.String(r.committerEmail),
	}

	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(""),
		Content:   []byte(content),
		Branch:    github.String(r.branch),
		Committer: committer,
	}

	if sha == "" {
		opts.Message = github.String("Created by hubdiary")
		_, _, err := r.repositories.CreateFile(ctx, r.owner, r.repo, path, opts)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	} else {
		opts.Message = github.String("Modified by hubdiary")
		opts.SHA = github.String(sha)
		_, _, err := r.repositories.UpdateFile(ctx, r.owner, r.repo, path, opts)
		if err != nil {
			return fmt.Errorf("failed to update file: %w", err)
		}
	}

	return nil
}
