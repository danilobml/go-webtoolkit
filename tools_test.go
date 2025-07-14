package gowebtoolkit

import "testing"

func TestTools_GenerateRandomString(t *testing.T) {
	var testTools Tools
	length := 10

	s := testTools.GenerateRandomString(length)
	
	if len(s) != length {
		t.Error("wrong length of the generated string")
	}
}
