package main

import (
	"fmt"

	comm "github.com/tiptok/gotransfer/comm"
	conn "github.com/tiptok/gotransfer/conn"
)

var srv conn.TcpServer

func main() {
	var param comm.Param
	//F:/App/Go/WorkSpace/src/github.com/tiptok/gotransfer/
	err := comm.ReadConfig(&param, "param.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(param)
}
