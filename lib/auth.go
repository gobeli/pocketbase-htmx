package lib

import (
	"fmt"
	"net/http"

	"github.com/gobeli/pocketbase-htmx/middleware"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

type Users struct {
	models.Record
}

func (*Users) TableName() string {
	return "users"
}

func Login(e *core.ServeEvent, c echo.Context, username string, password string) error {
	user, err := e.App.Dao().FindAuthRecordByUsername("users", username)
	if err != nil {
		return fmt.Errorf("Login failed")
	}

	valid := user.ValidatePassword(password)
	if !valid {
		return fmt.Errorf("Login failed")
	}

	return setAuthToken(e.App, c, user)
}

func Register(e *core.ServeEvent, c echo.Context, username string, password string, passwordRepeat string) error {
	user, _ := e.App.Dao().FindAuthRecordByUsername("users", username)
	if user != nil {
		return fmt.Errorf("username already taken")
	}

	if password != passwordRepeat {
		return fmt.Errorf("passwords don't match")
	}

	collection, err := e.App.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	newUser := models.NewRecord(collection)
	newUser.SetPassword(password)
	newUser.SetUsername(username)

	if err = e.App.Dao().SaveRecord(newUser); err != nil {
		return err
	}

	return setAuthToken(e.App, c, newUser)
}

func setAuthToken(app core.App, c echo.Context, user *models.Record) error {
	s, tokenErr := tokens.NewRecordAuthToken(app, user)
	if tokenErr != nil {
		return fmt.Errorf("Login failed")
	}

	c.SetCookie(&http.Cookie{
		Name:     middleware.AuthCookieName,
		Value:    s,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})

	return nil
}
