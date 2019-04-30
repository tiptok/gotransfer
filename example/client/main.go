package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tiptok/gotransfer/conn"
)

type TcpClientHandler struct {
}

func (cli *TcpClientHandler) OnConnect(c *conn.Connector) bool {
	log.Println("client On connect", c.RemoteAddress)
	return true
}
func (cli *TcpClientHandler) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	//cli.BizNode.SendToSrc(d)
	//log.Println("Rece Data:%v", d)
	return true
}
func (cli *TcpClientHandler) OnClose(c *conn.Connector) {
	//从列表移除
	log.Println("Client OnClose:", c.RemoteAddress)
}

func main() {
	var data string
	var termnum int
	var interval int
	var port   int
	var addr   string

	flag.StringVar(&addr,"addr","127.0.0.1","remote server IP")
	flag.IntVar(&port,"p",9927,"remote server port")
	flag.StringVar(&data,"d","hello tcp","send data to remote server")
	flag.IntVar(&termnum,"n",1,"client size link to remote server")
	flag.IntVar(&interval,"i",1,"interval to send data")

	flag.Parse()

	clientChan :=make(chan *conn.TcpClient,termnum)
	exit := make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(termnum)
	for i:=0;i<termnum;i++{
		go func(){
			defer func(){
				if p:=recover();p!=nil{

				}
			}()
			cli := &conn.TcpClient{}
			cli.NewTcpClient(addr, port, 500, 500)
			cli.Start(&TcpClientHandler{})
			wg.Done()
			for {
				if !cli.Conn.IsConneted {
					cli.ReStart() //重连
				}
				go func(){
					for{
						if cli.Conn.IsConneted{
							cli.Conn.SendChan<-conn.NewTcpData([]byte(data))
						}
						time.Sleep(time.Duration(int64(interval))*time.Second)
					}
				}()
				//log.Println("cli stop isconnnet:", cli.Conn.IsConneted, &cli)
				time.Sleep(10 * time.Second)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("start client num:%d  current:%d\n",termnum,len(clientChan))
	<-exit
}
