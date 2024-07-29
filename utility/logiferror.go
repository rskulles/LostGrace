package utility

import "log"

func LogIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}