package test

import (
	"fmt"
	"testing"
)

type T interface {
	tes()
}

type Test1 struct {
	age int
}

func (t *Test1) tes() {

}

func TestNew(t *testing.T) {
	var tt T
	tt = new(Test1)
	tt1 := &Test1{}
	fmt.Println(tt)
	fmt.Println(tt1)
}
