package model

import (
	"database/sql"
	"github.com/akula410/models"
	"model/example/db"
)

type Regular struct {
	models.Model
}

func (c Regular)Init()Regular{
	//Подключение к БД
	c.Conn = func() *sql.DB{
		return db.MySql.Connect()
	}

	//Отключение от БД (Можно не использовать если отключать в другом месте)
	c.ConnClose = func(){
		db.MySql.Close()
	}

	return c
}
