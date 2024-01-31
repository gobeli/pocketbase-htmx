package lib

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

func IsHtmxRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

type TemplRenderFunc func() templ.Component

func HtmxRender(c echo.Context, htmxTemplate TemplRenderFunc, nonHtmxTemplate TemplRenderFunc) templ.Component {
	if IsHtmxRequest(c) {
		return htmxTemplate()
	} else {
		return nonHtmxTemplate()
	}
}

func HtmxRedirect(c echo.Context, path string) error {
	if IsHtmxRequest(c) {
		c.Response().Header().Set("HX-Location", path)
		return c.NoContent(204)
	}

	return c.Redirect(302, path)
}
