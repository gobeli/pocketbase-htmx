package app

import (
	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gobeli/pocketbase-htmx/lib"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func TodosGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		var record *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)
		todos := []*models.Record{}
		e.App.Dao().RecordQuery("todos").Where(dbx.NewExp("user = {:id}", dbx.Params{"id": record.Id})).All(&todos)
		return lib.Render(c, 200, TodosList(todos))
	}
}

func TodoDelete(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		var authRecord *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)
		id := c.PathParam("id")
		record, err := e.App.Dao().FindRecordById("todos", id)

		if err != nil || record.Get("user") != authRecord.Id {
			c.NoContent(400)
			return nil
		}

		e.App.Dao().DeleteRecord(record)

		return lib.HtmxRedirect(c, "/app/todos")
	}
}

type AddTodoFormValue struct {
	name string
}

func (atfv AddTodoFormValue) Validate() error {
	return validation.ValidateStruct(&atfv,
		validation.Field(&atfv.name, validation.Required, validation.Length(1, 50)),
	)
}

func TodoAddGet(c echo.Context) error {
	return lib.Render(c, 200, TodoAdd(nil, nil))
}

func TodoAddPost(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		todo := AddTodoFormValue{name: c.FormValue("name")}
		err := todo.Validate()

		if err == nil {
			var authRecord *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)
			todos, _ := e.App.Dao().FindCollectionByNameOrId("todos")
			record := models.NewRecord(todos)

			record.Load(map[string]any{
				"user": authRecord.Id,
				"name": todo.name,
			})

			err = e.App.Dao().SaveRecord(record)

			if err == nil {
				return lib.HtmxRedirect(c, "/app/todos")
			}
		}

		component := lib.HtmxRender(
			c,
			func() templ.Component { return TodoAddForm(&todo, err) },
			func() templ.Component { return TodoAdd(&todo, err) },
		)

		return lib.Render(c, 200, component)
	}
}
