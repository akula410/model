package model

import (
	"database/sql"
	"github.com/akula410/models"
	"model/example/db"
)

type RBAC struct {
	models.Model
}

func (c RBAC)Init()RBAC{
	//Подключение к БД
	c.Conn = func() *sql.DB{
		return db.MySql.Connect()
	}

	//Отключение от БД (Можно не использовать если отключать в другом месте)
	c.ConnClose = func(){
		db.MySql.Close()
	}

	c.BeforeUpdates = append(c.BeforeUpdates, func(model *models.Model){
		fmt.Println(model.TableName)
	})


	return c
}
