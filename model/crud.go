package model

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Crud interface {
	models.Model
	GetUser() string
	Validate() error
}

func FindAll[C Crud](model Crud, dao *daos.Dao, authRecord *models.Record) ([]C, error) {
	items := []C{}
	error := dao.ModelQuery(model).Where(dbx.NewExp("user = {:id}", dbx.Params{"id": authRecord.Id})).All(&items)

	return items, error
}

func FindById(model Crud, dao *daos.Dao, authRecord *models.Record, id string) error {
	err := dao.FindById(model, id)

	if err != nil || model.GetUser() != authRecord.Id {
		return fmt.Errorf("record not found")
	}

	return nil
}

func Save(model Crud, dao *daos.Dao) error {
	err := model.Validate()
	if err == nil {
		return dao.Save(model)
	}

	return err
}
