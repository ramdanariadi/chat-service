package utils

import "log"

func LogIfError(err error) {
	if nil != err {
		log.Print(err)
	}
}
