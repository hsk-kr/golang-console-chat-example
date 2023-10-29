package server

import (
	"fmt"
	"io"
	"net"

	"github.com/hsk-kr/golang-console-chat-example/data"
)

type client struct {
	conn *net.Conn
	cs *chatServer
	nickname string
}

func (c *client) run() {
	byteData := make([]byte, data.MAX_BUFFER)

	for {
		len, err := (*c.conn).Read(byteData)

		if err != nil {
			fmt.Printf("Client Read Error: %s(%d)\n", err, len)
			if err == io.EOF {
				(*c.cs).disconnect(c)
				break
			}
		}

		d := data.ConvertBytesToData(byteData)

		switch d.MessageType {
			case data.MESSAGE_TYPE_NICKNAME:
				c.setNickname(d.Message)
			case data.MESSAGE_TYPE_MESSAGE:
				c.cs.broadMessage(c, d.Message)
		}
	}
}

func (c *client) getIP() string {
	addr, err := net.ResolveTCPAddr("tcp", (*c.conn).RemoteAddr().String())
	if err != nil {
		fmt.Printf("Client ResolveTCPAddr Error: %s\n", err)
		return ""
	}

	return addr.IP.String()
}

func (c *client) setNickname(nickname string) {
	c.nickname = nickname
}

func (c *client) getNickname() string {
	if c.nickname == "" {
		return "unknown"
	}

	return c.nickname
}

func (c *client) sendMessage(data []byte) {
	(*c.conn).Write(data)
}