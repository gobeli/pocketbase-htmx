package app

import (
	"github.com/gobeli/pocketbase-htmx/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func InitAppRoutes(e *core.ServeEvent, pb *pocketbase.PocketBase) {
	appGroup := e.Router.Group("/app", middleware.LoadAuthContextFromCookie(pb), middleware.AuthGuard)

	appGroup.GET("/profile", ProfileGet)
	appGroup.GET("/todos", TodosGet(e))
	appGroup.GET("/todos/add", TodoAddGet)
	appGroup.POST("/todos/add", TodoAddPost(e))
	appGroup.POST("/todos/:id/delete", TodoDelete(e))
}
