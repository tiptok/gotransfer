package main

import "log"

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
