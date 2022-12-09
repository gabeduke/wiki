package actions

import (
	"github.com/gabeduke/wiki/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"net/http"
)

// UsersResource is the resource for the Item model
type UsersResource struct {
	buffalo.Resource
}

func (v UsersResource) scope(c buffalo.Context) *pop.Query {
	tx := c.Value("tx").(*pop.Connection)
	cu, ok := c.Value("current_user").(*models.User)
	if !ok {
		return tx.Q()
	}
	return tx.BelongsTo(cu)
}

func UsersNew(c buffalo.Context) error {
	u := models.User{}
	c.Set("user", u)
	c.Bind(u)
	return c.Render(http.StatusOK, r.Auto(c, u))
}

// List gets all Users. This function is mapped to the path
// GET /users
func (v UsersResource) List(c buffalo.Context) error {
	users := &models.Users{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := v.scope(c).PaginateFromParams(c.Params())
	q = q.Order("created_at desc")

	// Retrieve all Items from the DB
	if err := q.All(users); err != nil {
		return err
	}

	// Make Items available inside the html template
	c.Set("users", users)

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, users))
}

//// UsersCreate registers a new user with the application.
//func UsersCreate(c buffalo.Context) error {
//	u := &models.User{}
//	if err := c.Bind(u); err != nil {
//		return err
//	}
//
//	tx := c.Value("tx").(*pop.Connection)
//	verrs, err := u.Create(tx)
//	if err != nil {
//		return err
//	}
//
//	if verrs.HasAny() {
//		c.Set("user", u)
//		c.Set("errors", verrs)
//		return c.Render(200, r.HTML("users/new.plush.html"))
//	}
//
//	c.Session().Set("current_user_id", u.ID)
//	c.Flash().Add("success", "Welcome to Toodo!")
//
//	return c.Redirect(302, "/items")
//}
