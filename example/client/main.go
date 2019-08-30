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
	ConnCount int64
}

func (cli *TcpClientHandler) OnConnect(c *conn.Connector) bool {
	atomic.AddInt64(&cli.ConnCount,1)
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
	atomic.AddInt64(&cli.ConnCount,-1)
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
	flag.IntVar(&interval,"i",1000,"interval to send data (Millisecond)")

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
	handle :=new(TcpClientHandler)
	for i:=0;i<termnum;i++{
		go func(){
			defer func(){
				if p:=recover();p!=nil{

				}
			}()
			cli := &conn.TcpClient{}
			cli.NewTcpClient(addr, port, handle)
			//cli.Start(&TcpClientHandler{})
			wg.Done()
			for {
				defer func(){
					if p:=recover();p!=nil{
						log.Println(p)
					}
				}()
				if cli.Conn==nil || !cli.Conn.IsConneted{
					cli.ReStart()
					time.Sleep(10 * time.Second)
					continue
				}
				//for{
				if cli.Conn!=nil && cli.Conn.IsConneted{
					atomic.AddInt64(&totalSend,int64(len([]byte(data))))
					atomic.AddInt64(&totalSendTime,int64(1))
					cli.Conn.SendChan<-conn.NewTcpData([]byte(data))
				}
				time.Sleep(time.Duration(int64(interval))*time.Millisecond)
				//}
			}
		}()
	}
	wg.Wait()
	fmt.Printf("start client num:%d  current:%d time:%v\n",termnum,len(clientChan),beginTime)

	go func(){
		var lastSendSize int64
		var lastInterval int64
		var firstTime bool =true
		onFirstTime :=func(){
			fmt.Print("\n\n")
			//fmt.Println("──────────┬────────────┬────────────┬────────────┬────────────")
			fmt.Println(fmt.Sprintf("%20v %20v%10v %20v %20v %20v","total_send_time","total_send_size","","cost_time","per_second_send","conn_count"))
			//fmt.Println("──────────┼────────────┼────────────┼────────────┬────────────")
			firstTime = false
		}
		interval:=time.Second*10
		ticker :=time.NewTicker(interval)
		perSecSendData :=func(curSize,lastSize int64,interval int64)float64{
			rsp :=float64(curSize-lastSize)/float64(interval-lastInterval)
			//fmt.Println(curSize,lastSize,interval)
			lastSendSize = curSize
			lastInterval = interval
			return rsp
		}
		printByteToMb := func(bSize float64)string{
			kbSize :=bSize/1024
			if kbSize<1024{
				return fmt.Sprintf("%.1f KB",kbSize)
			}
			return fmt.Sprintf("%.1f MB",kbSize/1024)
		}
		for{
			select{
			case <-ticker.C:
				if firstTime || handle.ConnCount!=int64(termnum){
					onFirstTime()
				}
				curTime :=time.Now().Unix()
				span :=curTime - beginTime
				total :=totalSend
				qps:=perSecSendData(total,lastSendSize,span)
				fmt.Println(fmt.Sprintf("%20d %20d%10v %20d %18v/s %20d",curTime,total,printByteToMb(float64(total)),span,printByteToMb(qps),handle.ConnCount))
			}
		}
	}()
	<-exit
}
