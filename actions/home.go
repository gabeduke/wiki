package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	c.Session().Set("current_workspace", 1)
	if err := c.Session().Save(); err != nil {
		return errors.WithStack(err)
	}

	if c.Value("current_user") != nil && c.Value("current_workspace") != nil {
		return c.Redirect(302, "/items")
	}
	return c.Render(200, r.HTML("index.plush.html"))
}
