package server

import (
	"fmt"
	"net"
	"os"

	"github.com/hsk-kr/golang-console-chat-example/config"
)

func Run() {
	endpoint := fmt.Sprintf("%s:%s", config.SERVER_HOST, config.SERVER_PORT) 
	server, err := net.Listen(config.SERVER_NETWORK, endpoint)

	if err != nil {
		fmt.Printf("net.Listen Error: %s\n", err)
		os.Exit(1)
	}
	defer server.Close()

	cs := new(chatServer)
	cs.server = &server

	fmt.Printf("Server has listened[%s:%s]\n", config.SERVER_HOST, config.SERVER_PORT)
	for {
		conn, err := server.Accept()
		
		if err != nil {
			fmt.Printf("server.Accept Error: %s\n", err)
			break
		}

		cs.connect(&conn)
	}
}