// 
// Copied from https://github.com/andrew/base62.go
// With modifications by John Nye
//
// heavily influcened by http://golang.org/src/pkg/encoding/base64/base64.go

// Package base62 implements base62 encoding 
package base62

import (
	"math"
	"strings"
)

/*
 * Encodings
 */


var encodeStd2 = [62]string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9"}

func EncodeInt(value int64) string{
	var output string
	var remainder int64
	var numberOfCycles int64

	base := int64(62)
	remainder = value % base
	output = encodeStd2[remainder]
	numberOfCycles = int64(math.Floor(float64(value / base)))-1

	if numberOfCycles >=0 && numberOfCycles < 62{
		output = strings.Join([]string{encodeStd2[numberOfCycles], output}, "")
	}

	for i := 0; i < 10; i++ {
		if numberOfCycles < 62 {
			break
		}
		i = 0
		remainder = numberOfCycles % base

		numberOfCycles = int64(math.Floor(float64(numberOfCycles / base)))-1
		output = strings.Join([]string{encodeStd2[remainder],output}, "")
		if(numberOfCycles <= 0){
			i=10
		}
	}
	return output
}