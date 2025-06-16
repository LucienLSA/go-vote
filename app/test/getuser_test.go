package test

import (
	"fmt"
	"govote/app/model"
	"testing"
)

func TestGetUserV1(t *testing.T) {
	model.NewMysql()
	defer model.Close()
	ret, _ := model.GetUserV1("lucien")
	fmt.Printf("ret:%+v", ret)
}
