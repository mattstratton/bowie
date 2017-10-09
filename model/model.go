package model

// GitHubURLs holds the URLs to be used when using github enterprise
type GitHubURLs struct {
	API      string `yaml:"api,omitempty"`
	Upload   string `yaml:"upload,omitempty"`
	Download string `yaml:"download,omitempty"`
}

// Repo represents any kind of repo (GitHub, GitLab, etc)
type Repo struct {
	Owner string `yaml:",omitempty"`
	Name  string `yaml:",omitempty"`

	// Capture all undefined fields and should be empty after loading
	XXX map[string]interface{} `yaml:",inline"`
}

// String is the concatenated repo name, i.e.,  owner/name
func (r Repo) String() string {
	return r.Owner + "/" + r.Name
}

// CommitAuthor is the author of a Git commit
type CommitAuthor struct {
	Name  string `yaml:",omitempty"`
	Email string `yaml:",omitempty"`
}

// Project includes all project configuration
type Project struct {
	ProjectName string `yaml:"project_name,omitempty"`

	// should be set if using github enterprise
	GitHubURLs GitHubURLs `yaml:"github_urls,omitempty"`

	// Capture all undefined fields and should be empty after loading
	XXX map[string]interface{} `yaml:",inline"`
}
