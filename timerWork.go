package main

import (
	"log"
	"time"

	"github.com/tiptok/gotransfer/conn"
)

func TimerReconn() {
	for {
		log.Println("TimerReconn On Work 连接重连")
		select {
		case <-timerReconn.C:
			for key1, _ := range tcpTransferList {
				for key2, v := range tcpTransferList[key1] {
					if !((*v.Conn).IsConneted) {
						if v.Start(&TransferClientHandler{}) {
							//启动成功
							log.Println(key2, "->", v.LocalAdr, " Reconn...")
						} else {
							continue
						}
					}
				}
			}
		}

		if !Alive {
			//timerReconn.Stop()
			break
		}
	}

}

func TimerClear() {
	for {
		select {
		case <-timerClear.C:
			log.Println("TimerClear On Work")
			for _, udp := range udpTo {
				timeDiff := time.Now().Sub(udp.HeartTime).Seconds()
				if timeDiff > 60 {
					//心跳时间相差60s 移除掉
					sLocal := udp.RemoteAddress
					if _, ok := udpFrom[sLocal]; ok {
						for _, trans := range udpFrom[sLocal] {
							defer func() {
								conn.MyRecover()
							}()
							if (*trans.Conn).IsConneted {
								(*trans.Conn).Close() //先关闭
							}
						}
						delete(udpFrom, sLocal) //移除
						log.Println("移除 离线upd连接：", sLocal)
					}
					if udp.IsConneted {
						(udp).Close()
					}
					//delete(udpTo, key) //移除
				}
			}
			timerClear.Reset(time.Second * 60) //重置定时器
		}
	}
}
