package main

import (
	"model/example/model"
)

func main(){
	modelRBAC := model.ExampleRBAC.Init()
	modelRBAC.Find(218)
	modelRBAC.Delete()

	modelRegular := model.ExampleRegular.Init()
	modelRegular.Find(218)
	modelRegular.Delete()
}
