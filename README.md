# GOLang CUI based Multichat Example

[Dev.to Post](https://dev.to/lico/golang-cli-based-multichat-example-37go)

---

# Post

![program output](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/601egax7ai2btq3rk7y0.png)

At the start of the program, you can choose whether it operates as a server or client. When a client provides their nickname, they can have conversations with other clients.

---

# What is Socket Programming?

Socket programming is a way of connecting two nodes on a network to communicate with each other. One socket(node) listens on a particular port at an IP, while the other socket reaches out to the other to form a connection. - [geeksforgeeks.org](https://www.geeksforgeeks.org/socket-programming-cc/)

---

# Source Code

## Structure

![Structure Diagram](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/2rf46676sd3re85pd25j.png)
[-excalidraw](https://excalidraw.com/)

I have included this diagram to help you get a better sense of how the project is organized.
Let's now explore the code for each module step by step.

---

## Config

**config.go**
```go
package config

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = "3000"
	SERVER_NETWORK = "tcp"
)
```

The constants are utilized for both the server and the client.

---

## Data

**data.go**
```go
package data

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

type Data struct {
	MessageType byte
	Message string
}

const (
	MAX_BUFFER = 1024
)

const (
	MESSAGE_TYPE_NICKNAME = 1
	MESSAGE_TYPE_MESSAGE = 2
)

func GenerateNicknameData(nickname string) []byte {
	return dataToBytes(MESSAGE_TYPE_NICKNAME, nickname)
}

func GenerateMessageData(message string) []byte {
	return dataToBytes(MESSAGE_TYPE_MESSAGE, message)
}

func ConvertBytesToData(b []byte) Data {
	data := Data{}
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	err := decoder.Decode(&data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return data
}

func dataToBytes(messageType byte, message string) []byte {
	d := Data{}
	d.MessageType = messageType
	d.Message = message

	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(d)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	return buffer.Bytes()
}
```

---

```go
type Data struct {
	MessageType byte
	Message string
}

const (
	MAX_BUFFER = 1024
)

const (
	MESSAGE_TYPE_NICKNAME = 1
	MESSAGE_TYPE_MESSAGE = 2
)
```

The `Data` structure consists of two fields `MessageType` and `Message`.
The `MessageType` can take one of two values, `MESSAGE_TYPE_NICKNAME` or `MESSAGE_TYPE_MESSAGE`.

When `MessageType` is set to `MESSAGE_TYPE_NICKNAME`, it indicates the user's intention to change their nickname. On the other hand, when set to `MESSAGE_TYPE_MESSAGE`, the message is sent to all connected clients.

The `MAX_BUFFER` specifies the maximum data size in bytes. As there is no size limit imposed in this app, issues may arise if the data exceeds 1024 bytes. it is going to be a problem. To prevent such problems, you have the option to create a custom protocol.

For instance, you can design the protocol as below.

```
Header(1 byte) | MessageSize(2 byte) | Message(X bytes)
```

The `Header` serves the purpose of marking the start of the data.

Subsequently, two bytes are used to represent the size of the data, allowing you to determine the data's length.

However, in the context of this example, we won't delve further into this particular scenario.

```go
func dataToBytes(messageType byte, message string) []byte {
	d := Data{}
	d.MessageType = messageType
	d.Message = message

	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(d)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	return buffer.Bytes()
}
```

the `dataToBytes` function converts the `Data` struct into a `[]byte`.

```go
func ConvertBytesToData(b []byte) Data {
	data := Data{}
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	err := decoder.Decode(&data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return data
}
```

the `ConvertBytesToData` function converts `[]byte` into the `Data` struct.

---

## Server

**main.go**
```go
package server

func Main() {
	Run()
}
```

This function is an entry point and is called from the main function.

**server.go**
```go
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
```

---

```go
endpoint := fmt.Sprintf("%s:%s", config.SERVER_HOST, config.SERVER_PORT) 
server, err := net.Listen(config.SERVER_NETWORK, endpoint)
```

The `net.Listen` function opens the server and accepts incoming connections.

```go
for {
	conn, err := server.Accept()
		
	if err != nil {
		fmt.Printf("server.Accept Error: %s\n", err)
		break
	}

	cs.connect(&conn)
}
```

In the main loop, it accepts the incoming connections and passes them to the `chatServer.connect` function. The core server logic is implemented within the `chatServer`.

**chat_server.go**
```go
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
```

---

```go
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
```

The `connect` function creates a client and adds it to the slice.

Since the `clients` slice can be accessed by other goroutines, I use a Mutex to ensure synchronization.

```go
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
```

The `disconnect` function takes the client as input. It proceeds to close the connection and remove the client from the slice.

```go
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
```

The `broadMessage` function is responsible for transmitting messages to all clients. Within this function, the client's nickname is prefixed to the message when the client is not `nil`.

**client.go**
```go
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
```

---

```go
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
```

The `run` function accepts data and checks the message type.

If the message type is `MESSAGE_TYPE_NICKNAME`, it updates the nickname.

If the message type is `MESSAGE_TYPE_MESSAGE`, it forwards the message to the `broadMessage` function within the `chatServer`.

---

## Client

**main.go**
```go
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
```

---

```go
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
```

At the program's start, it gets an input of the user's nickname and sends it to the server after establishing a connection.

Then, whenever the user inputs a message, it sends it to the server.

```go
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
```

The `runReceiveData` function receives the data from the server and displays the message on the screen.

To allow the user to send and receive messages simultaneously, this function runs as a goroutine.

```go
func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	str = strings.Replace(str, "\n", "", -1);
	return str
}
```

The `readLine` function inputs a string from the command line and returns it.

As the input data includes the new line character, it gets rid of the new line character using `strings.Replace`.

---

## Main

**main.go**
```go
package main

import (
	"fmt"

	"github.com/hsk-kr/golang-console-chat-example/client"
	"github.com/hsk-kr/golang-console-chat-example/server"
)

func main() {
	var t string

	for t != "c" && t != "s"{
		fmt.Print("Which one do you want to execute as server(s) or client(c):")
		fmt.Scanf("%s", &t);
	}
	
	if t == "s" {
		server.Main()
	} else {
		client.Main()
	}
}
```

The program requests the user to choose whether they want to run it as a server or a client. It then calls the respective `Main` function accordingly.

---

# Conclusion

This is my very first Go project, and it reminds me of the time I studied socket programming a little bit with C/C++ about 10 years ago. Actually, I made the code in a similar way, I think I haven't fully used Go's strengths.

These days, starting a coding project is kind of daunting. The pressure to create something big or impressive makes it hard to begin, especially when I see amazing projects from others, I feel kind of overwhelmed.

But now, I have made a decision to keep coding without overthinking. I will work on small pieces of code, one at a time, and write about them here. I believe that these small steps will add up and be worth it!

I hope you guys found it useful. 

Happy Coding!
