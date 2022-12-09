package grifts

import (
	"github.com/gabeduke/wiki/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
