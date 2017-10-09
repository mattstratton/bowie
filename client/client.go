// Package client contains the client implementations for SCM provider.
package client

import "github.com/google/go-github/github"

// Info of the repository
type Info struct {
	Description string
	Homepage    string
	URL         string
}

// Client interface
type Client interface {
	GetIssues(userName, projectName string) ([]*github.Issue, error)
}
