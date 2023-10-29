package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/hsk-kr/golang-console-chat-example/data"
)

type chatServer struct {
	server *net.Listener
	clients []*client
	muClients sync.Mutex
}

func (cs *chatServer) connect(conn *net.Conn) {
	newClient := client{}
	newClient.conn = conn
	newClient.cs = cs

	cs.muClients.Lock()
	defer cs.muClients.Unlock()
	cs.clients = append(cs.clients, &newClient)
	go newClient.run()
	fmt.Printf("client %s connected.\n", newClient.getIP())
}

func (cs *chatServer) disconnect(c *client) {
	idx := -1

	cs.muClients.Lock()
	defer cs.muClients.Unlock()
	for i, client := range cs.clients {
		if client == c {
			idx = i
			break
		}
	}

	if idx == -1 {
		fmt.Printf("client can't not be found. (%s)", c.getIP())
		return
	}

	targetClient := cs.clients[idx]
	fmt.Printf("client %s disconnected.\n", targetClient.getIP())
 	(*targetClient.conn).Close()
	cs.clients = append(cs.clients[:idx], cs.clients[idx + 1:]...)
}

func (cs *chatServer) broadMessage(c *client, message string) {
	if c != nil {
		message = fmt.Sprintf("%s:%s", c.getNickname(), message)
	}
	
	data := data.GenerateMessageData(message)
	cs.muClients.Lock()
	defer cs.muClients.Unlock()
	for _, client := range cs.clients {
		client.sendMessage(data)
	}
}