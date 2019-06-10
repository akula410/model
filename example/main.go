package main

import (
	"fmt"
	"model/example/model"
)

func main(){
	//Example 1
	modelRBAC := model.ExampleRBAC.Init()
	modelRBAC.Find(218)
	modelRBAC.Delete()


	modelRegular := model.ExampleRegular()
	if modelRegular.Data.Find(210) == true {
		fmt.Println(modelRegular.TrfField())
	}
}
