package model

var ExampleRegular RBAC

func init(){
	ExampleRegular.TableName = "test_data"
	ExampleRegular.PrimaryKey = "data_id"
}
