package app

import (
	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gobeli/pocketbase-htmx/lib"
	"github.com/gobeli/pocketbase-htmx/model"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func TodosGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		var record *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)
		todos, err := (&model.Todo{}).FindAll(e.App.Dao(), record)

		if err != nil {
			return err
		}

		return lib.Render(c, 200, TodosList(todos))
	}
}

func TodoDelete(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		var authRecord *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)

		id := c.PathParam("id")
		todo := model.Todo{}
		err := (&todo).FindById(e.App.Dao(), authRecord, id)

		if err == nil {
			err = todo.Delete(e.App.Dao())
		}

		if err != nil {
			c.NoContent(400)
			return nil
		}

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
		var authRecord *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)
		todo := model.Todo{Name: c.FormValue("name"), User: authRecord.Id}
		err := todo.Save(e.App.Dao())

		if err == nil {
			return lib.HtmxRedirect(c, "/app/todos")
		}

		component := lib.HtmxRender(
			c,
			func() templ.Component { return TodoAddForm(&todo, err) },
			func() templ.Component { return TodoAdd(&todo, err) },
		)

		return lib.Render(c, 200, component)
	}
}
