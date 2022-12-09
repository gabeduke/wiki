package grifts

import (
	"github.com/gabeduke/wiki/models"

	"github.com/gobuffalo/grift/grift"
	"github.com/gobuffalo/pop/v6"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			if err := tx.TruncateAll(); err != nil {
				return err
			}
			w := &models.Workspace{
				ID:   1,
				Name: "default",
			}

			if err := tx.Create(w); err != nil {
				return err
			}

			return nil
		})
	})

})
