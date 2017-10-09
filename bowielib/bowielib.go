package bowielib

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/mattstratton/bowie/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	userName    string
	projectName string
	token       string
)

// RootCmd is the main command executed when bowie is run
var RootCmd = &cobra.Command{
	Use:   "bowie",
	Short: "Bowie is a pretty changelog generator",
	Long: `A pretty changelog generator,
built with love by mattstratton in Go.
	
Complete documentation is available at https://github.com/mattstratton/bowie`,

	Run: func(bowielib *cobra.Command, args []string) {
		ChangeLog(userName, projectName) // nolint: errcheck
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&userName, "user", "u", "", "Username of the owner of target GitHub repo")
	RootCmd.PersistentFlags().StringVarP(&projectName, "project", "p", "", "Name of project on GitHub")
	RootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "To make more than 50 requests per hour your GitHub token is required. You can generate it at: https://github.com/settings/tokens/new")
	if (os.Getenv("BOWIE_GITHUB_TOKEN") == "") && (token == "") {
		fmt.Fprintf(color.Output, "\n%s%s%s%s%s%s\n", color.RedString("ERROR: "), color.YellowString("You need to either set the "), color.GreenString("BOWIE_GITHUB_TOKEN"), color.YellowString(" environment variable or use the "), color.GreenString("--token "), color.YellowString("flag."))
		os.Exit(1)
	}
	var err error
	token, err = GetToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Execute is the main root command of bowie
func Execute() {
	RootCmd.Execute() // nolint: errcheck
}

// GetToken returns the GitHub token based on environment variable or flag (flag takes precedence)
func GetToken() (string, error) {
	if token != "" {
		return token, nil
	}
	if os.Getenv("BOWIE_GITHUB_TOKEN") != "" {
		return os.Getenv("BOWIE_GITHUB_TOKEN"), nil
	}
	return "", errors.Wrap(nil, "token set failed")
}

// ChangeLog generates the changelog with the given flags/configuration
func ChangeLog(username, project string) error {

	c, _ := client.NewGitHub(token)

	fmt.Fprintf(color.Output, "Your username is %s and your projectname is %s \n", color.GreenString(userName), color.GreenString(projectName))

	issues, _ := c.GetIssues(userName, projectName)

	for _, issue := range issues {
		// closedate := issue.GetClosedAt()
		myTag, _ := GetIssueTag(issue.GetClosedAt())
		// closedate = closedate.String()
		fmt.Fprintf(color.Output, "Issue title is: %s and it was closed by %s\n", color.GreenString(issue.GetTitle()), color.GreenString(myTag))

	}
	// fmt.Println("Issues:")
	// for k, v := range issues {
	// 	myTag, _ := GetIssueTag(v)
	// 	fmt.Println("Issue: " + GetIssueNameByID(k) + " resolved in tag " + myTag)
	// }
	return nil

}

// GetIssues gets the list of all closed issues for the username and project, and returns a map of them
// func GetIssues(c *githubClient) (map[int]time.Time, error) {

// 	client := client.NewGitHub

// 	// list all issues in the repo
// 	issueOpts := &github.IssueListByRepoOptions{
// 		State: "closed",
// 	}
// 	issues, _, err := client.Issues.ListByRepo(ctx, userName, projectName, issueOpts)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "GitHub issue list failed")
// 	}

// 	m := make(map[int]time.Time)
// 	for _, d := range issues {
// 		m[d.GetID()] = d.GetClosedAt()
// 	}

// 	return m, nil
// }

// GetTags gets the list of all tags for the username and project, and returns a map of them
func GetTags() (map[string]time.Time, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all tags in the repo
	tags, _, err := client.Repositories.ListTags(ctx, userName, projectName, nil)

	if err != nil {
		return nil, errors.Wrap(err, "GitHub tag list failed")
	}

	m := make(map[string]time.Time)
	for _, d := range tags {
		sha := d.Commit.GetSHA()
		tag, _, _ := client.Git.GetCommit(ctx, userName, projectName, sha)
		m[d.GetName()] = tag.Author.GetDate()
	}
	return m, nil
}

// GetIssueTag takes in a date (from an issue, probably?) and returns the associated tag where it was closed.
// This is not nearly clear enough; it should be something like GetTagFromDate but even that is stupid.
func GetIssueTag(date time.Time) (string, error) {
	var myTag string
	tags, _ := GetTags()
	for k, v := range tags {
		if date.After(v) {
			myTag = k
		}
	}
	return myTag, nil

}

// GetIssueNameByID takes in a GitHub issue ID and returns the Title of the issue.
func GetIssueNameByID(id int) (name string) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all issues in the repo
	issueOpts := &github.IssueListByRepoOptions{
		State: "closed",
	}
	issues, _, _ := client.Issues.ListByRepo(ctx, userName, projectName, issueOpts)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "GitHub issue list failed")
	// }

	for _, d := range issues {
		if d.GetID() == id {
			return d.GetTitle()
		}
	}

	return ""
}
