package model

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
	result := c.Data.TrfToInt("field1")
	result *= 10
	return result
}
