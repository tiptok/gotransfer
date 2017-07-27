package conn

import (
	"fmt"
	"log"
)

func MyRecover() {
	if err := recover(); err != nil {
		log.Println(err)
		fmt.Println(err)
	}
}
