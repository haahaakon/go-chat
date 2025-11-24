package main

import (
	"fmt"
	"p3/chat"
	"p3/server"
)

func main() {
	messages := []chat.Message{}
	messages = append(messages, *chat.NewMessage(0, "Welcome to this chat!"))
	messages = append(messages, *chat.NewMessage(0, "---------------------"))

	fmt.Println(messages)

	server.StartServer(&messages)

}
