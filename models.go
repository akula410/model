package models

import (
	"database/sql"
	"fmt"
	"strings"
)


type Model struct {
	TableName string
	PrimaryKey string

	BeforeUpdate func()
	AfterUpdate func()
	BeforeInsert func()
	AfterInsert func()

	BeforeUpdates []func(model *Model)
	AfterUpdates []func(model *Model)
	BeforeInserts []func(model *Model)
	AfterInserts []func(model *Model)

	Conn func()*sql.DB
	ConnClose func()

	fields map[string]interface{}
	ok bool
	id interface{}
}

func (c *Model) Find(id interface{}) bool{

	selectWhere, params := c.fieldToSelect(id)
	textSql := fmt.Sprintf("SELECT * FROM %s WHERE %s", c.TableName, selectWhere)
	c.sql(textSql, params)
	return c.ok
}

func (c *Model) FindByCondition(condition map[string]interface{}) bool{
	selectWhere, params := c.fieldsToSelect(condition)
	textSql := fmt.Sprintf("SELECT * FROM %s WHERE %s", c.TableName, selectWhere)
	c.sql(textSql, params)
	return c.ok
}

func (c *Model) Update() int64{
	var result int64
	if c.GetId() == nil {
		result = c.Insert()
	} else {
		c.beforeUpdates()

		fieldSlice, updSet := c.fieldsToUpdate()
		textSql := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", c.TableName, updSet, c.PrimaryKey)

		ins, err := c.Conn().Prepare(textSql)
		if err != nil {
			panic(err.Error())
		}

		res, err := ins.Exec(fieldSlice...)

		if err != nil {
			panic(err.Error())
		}

		aff, err := res.RowsAffected()
		if err != nil {
			panic(err.Error())
		}

		result = aff
		c.afterUpdates()

		err = ins.Close()
		if err != nil {
			panic(err.Error())
		}
		c.callSqlClose()
	}


	c.callSqlClose()

	return result
}

func (c *Model) InsertFind() bool {
	c.Insert()
	return c.Find(c.id)
}


func (c *Model) Insert() int64{
	c.beforeInserts()

	c.id = nil
	c.ok = false
	fieldSlice, strInsert, strValue := c.fieldsToInsert()
	textSql := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", c.TableName, strInsert, strValue)

	ins, err := c.Conn().Prepare(textSql)
	if err != nil {
		panic(err.Error())
	}
	res, err := ins.Exec(fieldSlice...)

	if err != nil {
		panic(err.Error())
	}

	aff, err := res.RowsAffected()

	if err != nil {
		panic(err.Error())
	}


	c.id, err = res.LastInsertId()

	if err != nil {
		panic(err.Error())
	}

	c.ok = true

	c.afterInserts()
	err = ins.Close()
	if err != nil {
		panic(err.Error())
	}
	c.callSqlClose()

	return aff
}

func (c *Model) Field(field string)interface{}{
	return c.fields[field]
}

func (c *Model) GetFields() map[string]interface{} {
	return c.fields
}

func (c *Model) SetField(field string, value interface{}) *Model{

	if len(c.fields)==0 {
		c.fields = make(map[string]interface{})
	}
	c.fields[field] = value

	return c
}

func (c *Model) GetId() interface{}{
	if c.ok == false {
		return nil
	} else{
		return c.id
	}
}

func (c *Model) FlushData()*Model{
	c.fields = nil
	c.ok = false
	c.id = nil
	return c
}

func (c *Model) fieldToSelect(id interface{}) (string, []interface{}){
	selectWhere := fmt.Sprintf("%s=?", c.PrimaryKey)
	params := make([]interface{}, 0, 1)
	params = append(params, id)
	return selectWhere, params

}

func (c *Model) fieldsToSelect(condition map[string]interface{}) (string, []interface{}){
	fieldsWhere := make([]string, 0, len(condition))
	params := make([]interface{}, 0, len(condition))

	for field, value := range condition{
		switch value {
		case nil:
			fieldsWhere = append(fieldsWhere, fmt.Sprintf("%s IS ?", field))
		default:
			fieldsWhere = append(fieldsWhere, fmt.Sprintf("%s=?", field))
		}

		params = append(params, value)
	}

	return  strings.Join(fieldsWhere, " AND "), params
}

func (c *Model) fieldsToUpdate() ([]interface{}, string) {
	transformFields := make([]interface{}, 0, len(c.fields))
	resultArray := make([]string, 0, len(c.fields)-1)
	for field,value := range c.fields {
		if field != c.PrimaryKey {
			transformFields = append(transformFields, value)
			resultArray = append(resultArray, fmt.Sprintf("%s=?", field))

		}
	}
	transformFields = append(transformFields, c.id)
	return transformFields, strings.Join(resultArray, ", ")
}

func (c *Model) fieldsToInsert() ([]interface{}, string, string) {
	countFld := len(c.fields)
	if _, ok := c.fields[c.PrimaryKey]; ok {
		countFld--
	}

	transformFields := make([]interface{}, 0, countFld)
	resultInsertArray := make([]string, 0, countFld)
	resultValueArray := make([]string, 0, countFld)
	for field,value := range c.fields {
		if field != c.PrimaryKey {
			transformFields = append(transformFields, value)
			resultInsertArray = append(resultInsertArray, field)
			resultValueArray = append(resultValueArray, "?")
		}
	}

	return transformFields, strings.Join(resultInsertArray, ", "), strings.Join(resultValueArray, ", ")
}


func (c *Model) sql(textSql string, params []interface{}){
	c.ok = false
	rows, err := c.Conn().Query(textSql, params...)

	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()

	if err != nil {
		panic(err.Error())
	}

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	fields := make(map[string]interface{})


	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			panic(err.Error())
		}

		c.ok = true

		for i, col := range columns {

			var v interface{}

			val := values[i]

			b, ok := val.([]byte)

			if ok {
				v = string(b)
			} else {
				v = val
			}
			if col == c.PrimaryKey {
				c.id = v
			}
			fields[col] = v
		}
		break
	}

	c.fields = fields
	err = rows.Close()
	if err != nil {
		panic(err.Error())
	}
	c.callSqlClose()
}





func (c *Model) beforeUpdates(){
	c.callSliceFunc(c.BeforeUpdates)
	c.callFunc(c.BeforeUpdate)

}
func (c *Model) beforeInserts(){
	c.callSliceFunc(c.BeforeInserts)
	c.callFunc(c.BeforeInsert)
}

func (c *Model) afterUpdates(){
	c.callSliceFunc(c.AfterUpdates)
	c.callFunc(c.AfterUpdate)
}

func (c *Model) afterInserts(){
	c.callSliceFunc(c.AfterInserts)
	c.callFunc(c.AfterInsert)
}

func (c *Model) callSliceFunc(fn []func(model *Model)){
	if len(fn) > 0 {
		for _, f := range fn{
			f(c)
		}
	}
}

func (c *Model) callFunc(fn func()){
	if fn != nil {
		fn()
	}
}

func (c *Model) callSqlClose(){
	if c.ConnClose != nil {
		c.ConnClose()
	}
}
