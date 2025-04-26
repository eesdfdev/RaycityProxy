package Network

import (
	"RCProxy/Logger"
	"net"
)

var Clients []*Client

func Run() {
	listener, err := net.Listen("tcp", ":2080")
	if err != nil {
		Logger.Errorf("서버 소켓 바인딩 에러 : %v", err)
	}
	defer listener.Close()
	var socketId = 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			Logger.Errorf("서버 소켓 Accpet 에러 : %v", err)
			return
		}
		socketId += 1
		client := &Client{
			ConnId: socketId,
			Conn:   conn,
		}
		Clients = append(Clients, client)
		go HandleClient(client)
	}
}
