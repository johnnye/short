package main 

import "testing"
import "reflect"

func TestCreateShortURL(t *testing.T){
	var d Data
	v := createShortURL("http://john.com")
	if reflect.TypeOf(v) != reflect.TypeOf(d){
		t.Fail()
	}
}