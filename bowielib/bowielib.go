package bowielib

import (
	"context"
	"fmt"
	"os"
	"sort"
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

type myTag struct {
	Name string
	Date time.Time
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

	tags, _ := GetTags()
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Date.After(tags[j].Date)
	})
	for i, d := range tags {
		fmt.Fprintf(color.Output, "%s %s \n", color.GreenString("Release:"), color.BlueString(d.Name))
		fmt.Fprintf(color.Output, "%s \n", color.CyanString("Fixed issues:"))
		// fmt.Println("Current tag: " + d.Name)
		if i != (len(tags) - 1) { //@TODO@ Figure out how to get the last issues to show up
			// fmt.Println("Next tag: " + tags[i+1].Name)
			for _, issue := range issues {
				if (issue.GetClosedAt().Before(d.Date)) && (issue.GetClosedAt().After(tags[i+1].Date)) {
					fmt.Fprintf(color.Output, "* %s \n", color.GreenString(issue.GetTitle()))
					// fmt.Println(issue.GetTitle())
				}
			}

		}
	}

	return nil

}

// GetTags gets the list of all tags for the username and project, and returns a map of them
func GetTags() ([]*myTag, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc) //@TODO@ Change GetTags to be a methon on client

	// list all tags in the repo
	tags, _, err := client.Repositories.ListTags(ctx, userName, projectName, nil)

	if err != nil {
		return nil, errors.Wrap(err, "GitHub tag list failed")
	}
	taggers := []*myTag{}
	for _, d := range tags {
		sha := d.Commit.GetSHA()
		tag, _, _ := client.Git.GetCommit(ctx, userName, projectName, sha)
		someTag := new(myTag)
		someTag.Name = d.GetName()
		someTag.Date = tag.Author.GetDate()
		taggers = append(taggers, someTag)
	}
	return taggers, nil
}
