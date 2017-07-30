package main

import (
	"log"
	"time"
)

func Test() {
	t1 := time.Now()
	time.Sleep(time.Second * 3)
	t2 := time.Now()
	diff := t2.Sub(t1).Seconds()
	log.Panicln(diff)
}
