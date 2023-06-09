package utils

import "log"

func CheckErr(errMsg string, err error) {
	if err != nil {
		log.Println(errMsg, err)
	}
}
