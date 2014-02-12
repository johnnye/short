package main 

import "testing"
import "log"

func TestCreateShortURL(t *testing.T){
	v := createShortURL("http://john.com")
	log.Print(v)
	t.Fail()
}