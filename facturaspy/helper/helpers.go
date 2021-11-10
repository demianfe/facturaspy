package helper

import (
	"fmt"
)

// general helpers

//Contains returns true if a string is insede an array
func Contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

//CheckError function
func CheckError(err error) {
	//TODO: improve error handling
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
