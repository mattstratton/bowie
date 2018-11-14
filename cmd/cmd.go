package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"golang.org/x/oauth2"

	"text/template"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dimiro1/banner"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	colorable "github.com/mattn/go-colorable"
	"github.com/mattstratton/bowie/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	userName          string
	projectName       string
	token             string
	dateFormat        string
	outputFile        string
	bugsTitle         string
	enhancementsTitle string
	closedTitle       string
	logTitle          string
	bugLabel          string
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
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString(bowieBanner))

	RootCmd.PersistentFlags().StringVarP(&userName, "user", "u", "", "Username of the owner of target GitHub repo")
	RootCmd.PersistentFlags().StringVarP(&projectName, "project", "p", "", "Name of project on GitHub")
	RootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "To make more than 50 requests per hour your GitHub token is required. You can generate it at: https://github.com/settings/tokens/new")
	if (os.Getenv("BOWIE_GITHUB_TOKEN") == "") && (token == "") {
		fmt.Fprintf(color.Output, "\n%s%s%s%s%s%s\n", color.RedString("ERROR: "), color.YellowString("You need to either set the "), color.GreenString("BOWIE_GITHUB_TOKEN"), color.YellowString(" environment variable or use the "), color.GreenString("--token "), color.YellowString("flag."))
		os.Exit(1)
	}
	RootCmd.PersistentFlags().StringVarP(&dateFormat, "date-format", "f", "", "Date format. Default is %Y-%m-%d")
	RootCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "", "Output file. Default is CHANGELOG.md")
	RootCmd.PersistentFlags().StringVarP(&bugsTitle, "bugs-title", "", "", `Custom title for bug-fixes section. Default is "**Fixed bugs:**"`)
	RootCmd.PersistentFlags().StringVarP(&enhancementsTitle, "enhancements-title", "", "", "Custom title for enhancements section. Default is \"**Implemented enhancements:**\"")
	RootCmd.PersistentFlags().StringVarP(&closedTitle, "issues-title", "", "", `Custom title for closed-issues section. Default is "Closed issues:"`)
	RootCmd.PersistentFlags().StringVarP(&logTitle, "header-title", "", "", `Custom header title. Default is "# Change Log"`)
	RootCmd.PersistentFlags().StringVarP(&bugLabel, "bug-label", "", "", "Issues with the specified labels will be always added to \"Fixed bugs\" section. Default is 'bug'")

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

// ChangeTag is a git tag type, which includes all of the associated issues, etc
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
func ChangeLog(username, project string) error { // nolint: gocyclo

	c, _ := client.NewGitHub(token)

	fmt.Fprintf(color.Output, "Your username is %s and your projectname is %s \n", color.GreenString(userName), color.GreenString(projectName))

	fmt.Println("Getting list of all issues.....")

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

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

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
					case enhancementDisplayName:
						enhancements = append(enhancements, issue)
					case bugDisplayName:
						bugs = append(bugs, issue)
					case closedDisplayName:
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
					case enhancementDisplayName:
						enhancements = append(enhancements, issue)
					case bugDisplayName:
						bugs = append(bugs, issue)
					case closedDisplayName:
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

	defer func() {
		if terr := t.Execute(f, ChangeTags); terr != nil && err == nil {
			err = terr
		}
	}()

	// t.Execute(f, ChangeTags)
	// if err != nil {
	// 	return errors.Wrap(err, "template creation failure")
	// } else {
	// 	fmt.Println("Created changelog")
	// }
	return nil

}

// GetIssueType is a utility function to obtain the type of issue (enhancement, bug, or other)
func GetIssueType(issue *github.Issue) string {
	for _, l := range issue.Labels {
		switch l.GetName() {
		case enhancementTag:
			return enhancementDisplayName
		case bugTag:
			return bugDisplayName
		default:
			return closedDisplayName
		}
	}
	return closedDisplayName
}

// GetTags gets the list of all tags for the username and project, and returns a map of them
func GetTags() ([]*ChangeTag, error) {
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

	tags := []*ChangeTag{}
	fmt.Println("Getting list of all tags...")
	for _, d := range allTags {
		// sha := d.Commit.GetSHA()
		tag, _, _ := client.Git.GetCommit(ctx, userName, projectName, d.Commit.GetSHA())
		ref, _, _ := client.Git.GetRef(ctx, userName, projectName, d.GetName())

		thisTag := new(ChangeTag)
		thisTag.Name = d.GetName()
		thisTag.Date = tag.Author.GetDate()
		thisTag.RefURL = ref.GetURL()
		tags = append(tags, thisTag)
	}
	return tags, nil
}
