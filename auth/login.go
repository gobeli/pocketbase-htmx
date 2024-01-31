package auth

import (
	"net/http"

	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gobeli/pocketbase-htmx/lib"
	"github.com/gobeli/pocketbase-htmx/middleware"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type LoginFormValue struct {
	username string
	password string
}

func (lfv LoginFormValue) Validate() error {
	return validation.ValidateStruct(&lfv,
		validation.Field(&lfv.username, validation.Required, validation.Length(3, 50)),
		validation.Field(&lfv.password, validation.Required),
	)
}

func getLoginFormValue(c echo.Context) LoginFormValue {
	return LoginFormValue{
		username: c.FormValue("username"),
		password: c.FormValue("password"),
	}
}

func RegisterLoginRoutes(e *core.ServeEvent, group echo.Group) {
	group.GET("/login", func(c echo.Context) error {
		if c.Get(apis.ContextAuthRecordKey) != nil {
			return c.Redirect(302, "/app/profile")
		}

		return lib.Render(c, 200, Login(LoginFormValue{}, nil))
	})

	group.POST("/login", func(c echo.Context) error {
		form := getLoginFormValue(c)
		err := form.Validate()

		var token *string
		if err == nil {
			token, err = lib.Login(e, form.username, form.password)
		}

		if err != nil {
			component := lib.HtmxRender(
				c,
				func() templ.Component { return LoginForm(form, err) },
				func() templ.Component { return Login(form, err) },
			)
			return lib.Render(c, 200, component)
		}

		c.SetCookie(&http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    *token,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		})

		return lib.HtmxRedirect(c, "/app/profile")
	})

	group.POST("/logout", func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    "",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			MaxAge:   -1,
		})

		return lib.HtmxRedirect(c, "/auth/login")
	})
}
