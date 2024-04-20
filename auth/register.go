package auth

import (
	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gobeli/pocketbase-htmx/lib"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type RegisterFormValue struct {
	username       string
	password       string
	passwordRepeat string
}

func (lfv RegisterFormValue) Validate() error {
	return validation.ValidateStruct(&lfv,
		validation.Field(&lfv.username, validation.Required, validation.Length(3, 50)),
		validation.Field(&lfv.password, validation.Required),
	)
}

func getRegisterFormValue(c echo.Context) RegisterFormValue {
	return RegisterFormValue{
		username:       c.FormValue("username"),
		password:       c.FormValue("password"),
		passwordRepeat: c.FormValue("passwordRepeat"),
	}
}

func RegisterRegisterRoutes(e *core.ServeEvent, group echo.Group) {
	group.GET("/register", func(c echo.Context) error {
		if c.Get(apis.ContextAuthRecordKey) != nil {
			return c.Redirect(302, "/app/profile")
		}

		return lib.Render(c, 200, Register(RegisterFormValue{}, nil))
	})

	group.POST("/register", func(c echo.Context) error {
		form := getRegisterFormValue(c)
		err := form.Validate()

		if err == nil {
			err = lib.Register(e, c, form.username, form.password, form.passwordRepeat)
		}

		if err != nil {
			component := lib.HtmxRender(
				c,
				func() templ.Component { return RegisterForm(form, err) },
				func() templ.Component { return Register(form, err) },
			)
			return lib.Render(c, 200, component)
		}

		return lib.HtmxRedirect(c, "/app/profile")
	})
}
