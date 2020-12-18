package main

import ( "testing")
func TestNotEmptyString(t *testing.T){
	args := []string{""  , "path",  "/my/path"}
	for _, arg := range args{
		NotEmptyString(arg)
	}
}
