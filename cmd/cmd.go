package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"golang.org/x/oauth2"

	"text/template"

	rice "github.com/GeertJohan/go.rice"
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

	Run: func(cmd *cobra.Command, args []string) {
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

// MyTag is the struct for a tag
type MyTag struct {
	Name   string
	Date   time.Time
	RefURL string
}

type ChangeTag struct {
	Name         string
	Date         time.Time
	RefURL       string
	Enhancements []*github.Issue
	Bugs         []*github.Issue
	ClosedIssues []*github.Issue
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

	// find a rice.Box
	// to compile,run `rice embed-go`
	templateBox, err := rice.FindBox("../templates")
	if err != nil {
		log.Fatal(err)
	}

	templateName, _ := templateBox.String("CHANGELOG.md.tmpl")

	t, err := template.New("CHANGELOG.md").Parse(templateName)
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Create("CHANGELOG.md")
	// if err != nil {
	// 	return "Cannot create", err
	// }
	defer f.Close()

	var ChangeTags []ChangeTag
	for i, d := range tags {
		var thisTag ChangeTag
		thisTag.Name = d.Name
		thisTag.Date = d.Date
		thisTag.RefURL = d.RefURL
		if i != (len(tags) - 1) {
			enhancements := []*github.Issue{}
			bugs := []*github.Issue{}
			closed := []*github.Issue{}
			for _, issue := range issues {

				if (issue.GetClosedAt().Before(d.Date)) && (issue.GetClosedAt().After(tags[i+1].Date)) {

					switch GetIssueType(issue) {
					case "enhancement":
						enhancements = append(enhancements, issue)
					case "bug":
						bugs = append(bugs, issue)
					case "closed":
						closed = append(closed, issue)
					}
					if len(enhancements) > 0 {
						thisTag.Enhancements = enhancements
					}
					if len(bugs) > 0 {
						thisTag.Bugs = bugs
					}
					if len(closed) > 0 {
						thisTag.ClosedIssues = closed
					}

				}
			}

		} else {
			enhancements := []*github.Issue{}
			bugs := []*github.Issue{}
			closed := []*github.Issue{}
			for _, issue := range issues {

				if issue.GetClosedAt().Before(d.Date) {

					switch GetIssueType(issue) {
					case "enhancement":
						enhancements = append(enhancements, issue)
					case "bug":
						bugs = append(bugs, issue)
					case "closed":
						closed = append(closed, issue)
					}
					if len(enhancements) > 0 {
						thisTag.Enhancements = enhancements
					}
					if len(bugs) > 0 {
						thisTag.Bugs = bugs
					}
					if len(closed) > 0 {
						thisTag.ClosedIssues = closed
					}

				}
			}

		}
		ChangeTags = append(ChangeTags, thisTag)
	}
	t.Execute(f, ChangeTags)
	if err != nil {
		fmt.Println(err, "template execute error")
	} else {
		fmt.Println("Created changelog")
	}
	return nil

}

func GetIssueType(issue *github.Issue) string {
	for _, l := range issue.Labels {
		switch l.GetName() {
		case "enhancement":
			return "enhancement"
		case "bug":
			return "bug"
		default:
			return "closed"
		}
	}
	return "closed"
}

// GetTags gets the list of all tags for the username and project, and returns a map of them
func GetTags() ([]*MyTag, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	var allTags []*github.RepositoryTag

	client := github.NewClient(tc) //@TODO Change GetTags to be a methon on client
	tagOpts := &github.ListOptions{
		PerPage: 50,
	}
	for {
		tags, resp, err := client.Repositories.ListTags(
			ctx,
			userName,
			projectName,
			tagOpts,
		)

		if err != nil {
			return nil, errors.Wrap(err, "GitHub tag list failed")
		}
		allTags = append(allTags, tags...)
		if resp.NextPage == 0 {
			break
		}
	}

	taggers := []*MyTag{}
	for _, d := range allTags {
		sha := d.Commit.GetSHA()
		tag, _, _ := client.Git.GetCommit(ctx, userName, projectName, sha)
		ref, _, _ := client.Git.GetRef(ctx, userName, projectName, d.GetName())
		someTag := new(MyTag)
		someTag.Name = d.GetName()
		someTag.Date = tag.Author.GetDate()
		someTag.RefURL = ref.GetURL()
		taggers = append(taggers, someTag)
	}
	return taggers, nil
}
