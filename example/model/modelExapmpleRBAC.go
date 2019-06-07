package model

var ExampleRBAC RBAC

func init(){
	ExampleRBAC.TableName = "test_data"
	ExampleRBAC.PrimaryKey = "data_id"
}
