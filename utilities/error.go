package utilities

import "log"

func CheckError(err error) {
	if err != nil {
		log.Println("ERROR: ", err)
	}
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}