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