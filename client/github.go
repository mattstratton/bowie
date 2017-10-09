package client

import (
	"context"
	"net/url"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubClient struct {
	client *github.Client
}

// NewGitHub returns a github client implementation
func NewGitHub(ctx *context.Context) (Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ctx.Token},
	)
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	if ctx.Config.GitHubURLs.API != "" {
		api, err := url.Parse(ctx.Config.GitHubURLs.API)
		if err != nil {
			return &githubClient{}, err
		}
		upload, err := url.Parse(ctx.Config.GitHubURLs.Upload)
		if err != nil {
			return &githubClient{}, err
		}
		client.BaseURL = api
		client.UploadURL = upload
	}

	return &githubClient{client}, nil
}
