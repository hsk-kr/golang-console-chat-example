package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/hsk-kr/golang-console-chat-example/config"
	d "github.com/hsk-kr/golang-console-chat-example/data"
)

func Main() {
	fmt.Printf("Enter your nickname:")
	nickname := readLine()

	address := fmt.Sprintf("%s:%s", config.SERVER_HOST, config.SERVER_PORT)
	conn, err := net.Dial(config.SERVER_NETWORK, address)

	if err != nil {
		fmt.Printf("net.Dial Error: %s\n", err)
		os.Exit(1)
	}

	sendNicknameChangeRequest(&conn, nickname)
	go runReceiveData(&conn)

	for {
		message := readLine()
		sendMessage(&conn, message)
	}
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	str = strings.Replace(str, "\n", "", -1);
	return str
}

func sendNicknameChangeRequest(conn *net.Conn, nickname string) {
	d := d.GenerateNicknameData(nickname)
	(*conn).Write(d)
}

func sendMessage(conn *net.Conn, message string) {
	d := d.GenerateMessageData(message)
	(*conn).Write(d)
}

func runReceiveData(conn *net.Conn) {
	defer (*conn).Close()

	byteData := make([]byte, d.MAX_BUFFER)
	for {
		len, err := (*conn).Read(byteData)

		if err != nil {
			fmt.Printf("Client Read Error: %s(%d)\n", err, len)
			if err == io.EOF {
				os.Exit(1)
			}
		}

		data := d.ConvertBytesToData(byteData)
		fmt.Printf("%s\n", data.Message)
	}
}