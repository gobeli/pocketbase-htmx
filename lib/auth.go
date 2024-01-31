package lib

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tokens"
)

func Login(e *core.ServeEvent, username string, password string) (*string, error) {
	user, err := e.App.Dao().FindAuthRecordByUsername("users", username)
	if err != nil {
		return nil, fmt.Errorf("Login failed")
	}

	valid := user.ValidatePassword(password)
	if !valid {
		return nil, fmt.Errorf("Login failed")
	}

	s, tokenErr := tokens.NewRecordAuthToken(e.App, user)
	if tokenErr != nil {
		return nil, fmt.Errorf("Login failed")
	}

	return &s, nil
}
