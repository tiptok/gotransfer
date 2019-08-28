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
	flag.StringVar(&data,"d","hello tcp ","send data to remote server")
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
		totalSendTime int64
	)
	for i:=0;i<termnum;i++{
		go func(){
			defer func(){
				if p:=recover();p!=nil{

				}
			}()
			cli := &conn.TcpClient{}
			handle :=new(TcpClientHandler)
			cli.NewTcpClient(addr, port, handle)
			//cli.Start(&TcpClientHandler{})
			wg.Done()
			for {
				defer func(){
					if p:=recover();p!=nil{

					}
				}()
				if cli.Conn==nil || !cli.Conn.IsConneted{
					cli.ReStart() //重连
					time.Sleep(10 * time.Second)
					continue
				}
				//for{
				if cli.Conn!=nil && cli.Conn.IsConneted{
					atomic.AddInt64(&totalSend,int64(len([]byte(data))))
					atomic.AddInt64(&totalSendTime,int64(1))
					cli.Conn.SendChan<-conn.NewTcpData([]byte(data))
				}
				time.Sleep(time.Duration(int64(interval))*time.Second)
				//}
			}
		}()
	}
	wg.Wait()
	fmt.Printf("start client num:%d  current:%d time:%v\n",termnum,len(clientChan),beginTime)

	go func(){
		for{
			time.Sleep(time.Second*10)
			curTime :=time.Now().Unix()
			span :=curTime - beginTime
			log.Println(fmt.Sprintf("time:%v total-send-size:%v per-second-Send:%v 消耗时间:%v s total_send_time:%v",curTime,totalSend,totalSend/span,int64(span),totalSendTime))
		}
	}()
	<-exit
}
