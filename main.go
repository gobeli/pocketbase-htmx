package main

import (
	"log"

	"github.com/gobeli/pocketbase-htmx/app"
	"github.com/gobeli/pocketbase-htmx/auth"
	"github.com/gobeli/pocketbase-htmx/middleware"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const AuthCookieName = "Auth"

func main() {
	pb := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Static("/public", "public")

		authGroup := e.Router.Group("/auth", middleware.LoadAuthContextFromCookie(pb))
		auth.RegisterLoginRoutes(e, *authGroup)

		app.InitAppRoutes(e, pb)
		return nil
	})

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
