package model

import (
	"fmt"
	"strconv"
)

type tmp struct {
	Data Regular
}


func ExampleRegular() *tmp{
	tmp := &tmp{}
	tmp.Data.TableName = "test_data"
	tmp.Data.PrimaryKey = "data_id"
	tmp.Data = tmp.Data.Init()
	return tmp
}

func (c *tmp) TrfField() int64{
	var result int64

	fields := c.Data.GetFields()
	field, ok := fields["field1"]
	if ok == false {
		result = 0
	}else {
		i, err := strconv.ParseInt(fmt.Sprintf("%v", field), 10, 64)
		if err != nil {
			panic(err)
		}
		result = i*10
	}

	return result
}
