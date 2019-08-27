package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tiptok/gotransfer/conn"
)

type TcpClientHandler struct {
}

func (cli *TcpClientHandler) OnConnect(c *conn.Connector) bool {
	log.Println("client On connect", c.RemoteAddress)
	return true
}
func (cli *TcpClientHandler) OnReceive(c *conn.Connector, d *conn.TcpData) bool {
	//cli.BizNode.SendToSrc(d)
	//log.Println("Rece Data:%v", d)
	return true
}
func (cli *TcpClientHandler) OnClose(c *conn.Connector) {
	//从列表移除
	log.Println("Client OnClose:", c.RemoteAddress,c.Status())
}

func main() {
	var data string
	var termnum int
	var interval int
	var port   int
	var addr   string

	flag.StringVar(&addr,"addr","127.0.0.1","remote server IP")
	flag.IntVar(&port,"p",9928,"remote server port")
	flag.StringVar(&data,"d","hello tcp","send data to remote server")
	flag.IntVar(&termnum,"n",1,"client size link to remote server")
	flag.IntVar(&interval,"i",1,"interval to send data")

	flag.Parse()

	clientChan :=make(chan *conn.TcpClient,termnum)
	exit := make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(termnum)

	//统计信息
	var(
		beginTime int64 = time.Now().Unix()
		totalSend int64
	)
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
					defer func(){
						if p:=recover();p!=nil{

						}
					}()
					for{
						if cli.Conn.IsConneted{
							atomic.AddInt64(&totalSend,int64(len(data)))
							cli.Conn.SendChan<-conn.NewTcpData([]byte(data))
						}
						time.Sleep(time.Duration(int64(interval))*time.Second)
					}
				}()
				time.Sleep(10 * time.Second)
				//log.Println("cli stop isconnnet:", cli.Conn.IsConneted, &cli)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("start client num:%d  current:%d\n",termnum,len(clientChan))

	go func(){
		for{
			time.Sleep(time.Second*10)
			curTime :=time.Now().Unix()
			span :=curTime - beginTime
			fmt.Println(fmt.Printf("TotalSendSize:%v PersecondSend:%v Run:%v S",totalSend,totalSend/span,span))
		}
	}()
	<-exit
}
