// Package context provides bowie context
//
// The context extends the standard library context and is inspired
// by github.com/goreleaser/goreleaser/context.
package context

import "github.com/goreleaser/goreleaser/config"

// Context is the overall context
type Context struct {
	ctx.Context
	Config config.Project
	Token  string
	//  @todo Add feature for extra release notes
	//  @body There should be a feature to add release commentary to the changelog.
	Validate    bool
	Parallelism int
}

// New context
func New(config config.Project) *Context {
	return &Context{
		Context:     ctx.Background(),
		Config:      config,
		Parallelism: 4,
	}
}
