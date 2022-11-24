package grifts

import (
	"github.com/gobuffalo/grift/grift"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"wiki/models"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			if err := tx.TruncateAll(); err != nil {
				return err
			}
			u := &models.User{
				//Email:                "mark@example.com",
				//Password:             "password",
				//PasswordConfirmation: "password",
			}

			item := &models.Item{
				Title:  "Buy Milk",
				UserID: u.ID,
			}

			if err := tx.Create(item); err != nil {
				return err
			}

			item = &models.Item{
				Title:  "Learn Go",
				Body:   nulls.NewString("What a great language"),
				UserID: u.ID,
			}

			if err := tx.Create(item); err != nil {
				return err
			}

			item = &models.Item{
				Title:     "Learn Buffalo",
				UserID:    u.ID,
				Completed: true,
			}

			if err := tx.Create(item); err != nil {
				return err
			}
			return nil
		})
	})

})
