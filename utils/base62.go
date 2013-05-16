// 
// Copied from https://github.com/andrew/base62.go
// With modifications by John Nye
//
// heavily influcened by http://golang.org/src/pkg/encoding/base64/base64.go

// Package base62 implements base62 encoding 
package base62

import (
	"math"
)

/*
 * Encodings
 */


var encodeStd2 = []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","0","1","2","3","4","5","6","7","8","9"}

func EncodeInt(value int64) string{
	var output string
	var remainder int64
	var xyz int64

	base := int64(62)
	remainder = value % base
	output = encodeStd2[remainder]
	xyz = int64(math.Floor(float64(xyz / base)))
	for i := 0; i < 10; i++ {
		if xyz == 0 {
			break
		}
		i = 0
		remainder = xyz % base

		xyz = int64(math.Floor(float64(xyz / base)))
		output = encodeStd2[remainder]+output
		if(xyz == 0){
			i=10
		}
	}
	return output
}