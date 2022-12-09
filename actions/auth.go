package actions

import (
	"fmt"
	"github.com/gabeduke/wiki/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/markbates/going/defaults"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
	"time"
)

type AuthClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func (c *AuthClaims) Valid() error {
	if c.ID == uuid.Nil {
		return errors.New("invalid id")
	}
	return nil
}

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/google/callback")),
	)
}

func AuthToken(c buffalo.Context) error {
	u := c.Value("current_user").(*models.User)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		u.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 6, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "wiki",
			Subject:   "api_auth",
			ID:        uuid.Must(uuid.NewV4()).String(),
			Audience:  []string{u.Email.String},
		},
	})

	s, err := envy.MustGet("JWT_SECRET")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	tokenString, err := token.SignedString([]byte(s))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	return c.Render(http.StatusOK, r.JSON(map[string]string{"token": tokenString}))
}

func AuthCallback(c buffalo.Context) error {
	gu, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("provider = ? and provider_id = ?", gu.Provider, gu.UserID)
	exists, err := q.Exists("users")
	if err != nil {
		return errors.WithStack(err)
	}
	u := &models.User{}
	if exists {
		if err = q.First(u); err != nil {
			return errors.WithStack(err)
		}
	}
	u.Name = defaults.String(gu.Name, gu.NickName)
	u.Provider = gu.Provider
	u.ProviderID = gu.UserID
	u.Email = nulls.NewString(gu.Email)
	if err = tx.Save(u); err != nil {
		return errors.WithStack(err)
	}

	c.Session().Set("current_user_id", u.ID)
	if err = c.Session().Save(); err != nil {
		return errors.WithStack(err)
	}
	c.Flash().Add("success", "You have been logged in")
	return c.Redirect(http.StatusPermanentRedirect, "/")
}

func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out")
	return c.Redirect(302, "/")
}

func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			if err := tx.Find(u, uid); err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

func parseToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		s, err := envy.MustGet("JWT_SECRET")
		if err != nil {
			return nil, err
		}
		return []byte(s), nil
	})

	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return claims, err
}
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		const BearerSchema = "Bearer "
		authHeader := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(authHeader, BearerSchema) {
			token := authHeader[len(BearerSchema):]

			claim, err := parseToken(token)
			if err != nil {
				return c.Error(http.StatusUnauthorized, err)
			}

			if userId := claim.ID; userId != uuid.Nil {
				c.Session().Set("current_user_id", claim.ID)
				if err = c.Session().Save(); err != nil {
					return errors.WithStack(err)
				}
				return next(c)
			}
		}

		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}
