package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Todo struct {
	models.BaseModel
	User string `db:"user" json:"user"`
	Name string `db:"name" json:"name"`
}

func (*Todo) TableName() string {
	return "todos"
}

func (todo *Todo) GetUser() string {
	return todo.User
}

func (todo *Todo) Validate() error {
	return validation.ValidateStruct(todo,
		validation.Field(&todo.Name, validation.Required, validation.Length(1, 50)),
	)
}

func (todo *Todo) FindById(dao *daos.Dao, authRecord *models.Record, id string) error {
	return FindById(todo, dao, authRecord, id)
}

func (todo *Todo) FindAll(dao *daos.Dao, authRecord *models.Record) ([]*Todo, error) {
	return FindAll[*Todo](todo, dao, authRecord)
}

func (todo *Todo) Save(dao *daos.Dao) error {
	return Save(todo, dao)
}

func (todo *Todo) Delete(dao *daos.Dao) error {
	return dao.Delete(todo)
}
