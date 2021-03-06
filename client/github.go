package client

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type githubClient struct {
	client *github.Client
}

// NewGitHub returns a github client implementation
func NewGitHub(token string) (Client, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &githubClient{client}, nil
}

// GetIssues gets the list of all closed issues for the username and project, and  returns a map of them
func (c *githubClient) GetIssues(userName, projectName string) ([]*github.Issue, error) {
	ctx := context.Background()

	// list all issues in the repo
	issueOpts := &github.IssueListByRepoOptions{
		State:       "closed",
		ListOptions: github.ListOptions{PerPage: 300},
	}

	var allIssues []*github.Issue
	for {
		issues, resp, err := c.client.Issues.ListByRepo(
			ctx,
			userName,
			projectName,
			issueOpts,
		)
		fmt.Print("* ")
		if err != nil {
			return nil, errors.Wrap(err, "GitHub issue list failed")
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		issueOpts.Page = resp.NextPage
	}
	return allIssues, nil
}
