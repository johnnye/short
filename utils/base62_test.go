package base62 

import (
	"testing"
	"fmt"
)

func TestEncodeIntSmall(t *testing.T){
	var v string
	v = EncodeInt(0)
	if v != "A"{
		t.Fail()
	}
}

func TestMANYNumbers(t * testing.T){
	var out string 
	out = EncodeInt(int64(3843))
	if out != "89"{
		t.Fail()
	}
}

func TestEncodeLargeNumber(t *testing.T){
	v := EncodeInt(int64(99))
	if v != "Al"{
		fmt.Println("This is the:",v)
		t.Fail()
	}
}

