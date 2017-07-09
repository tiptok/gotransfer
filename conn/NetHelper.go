package conn

import (
	"log"
)

func MyRecover() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}

