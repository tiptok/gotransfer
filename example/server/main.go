package main

import (
	"flag"
	"fmt"
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/tiptok/gotransfer/conn"
)

type SimpleServerHandler struct {
}

//OnConnect handler
func (trans *SimpleServerHandler) OnConnect(c *conn.Connector) bool {
	log.Println(c.RemoteAddress, "On Connect.")
	return true
}

//OnClose handler
func (trans *SimpleServerHandler) OnClose(c *conn.Connector) {
	log.Println(c.RemoteAddress, "On Close.")
}

var exit chan int

//OnReceive handler
func (trans *SimpleServerHandler) OnReceive(c *conn.Connector, d *conn.TcpData) bool {
	//log.Printf("%v On Receive Data : %v\n", c.RemoteAddress, hex.EncodeToString(d.Bytes()))
	c.SendChan <- d
	return true
}

func main() {
	defer func(){
		if p:=recover();p!=nil{
			log.Println("panic,recover...")
			log.Println(p)
		}
	}()

	var (
		port int
	)
	flag.IntVar(&port,"p",9928,"server listen port")

	flag.Parse()
	//启动tcp服务
	go func() {
		config :=&conn.Conifg{
			ListenPort:port,
			SendSize:100,
			ReceiveSize:100,
			PackageSize:128,
			PoolMinSize:64,
			PoolMaxSize:1024*1024*10,
		}
		var srv = conn.NewTcpServer(config,&SimpleServerHandler{},nil)
		srv.Start()
	}()

	go func() {
		err:= http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 9929), nil)
		if err!=nil{
			log.Println(err)
		}
	}()
	//等待退出
	<-exit
}
