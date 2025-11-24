package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"p3/chat"
	"strconv"
	"strings"
)

var clients []net.Conn

func StartServer(messages *[]chat.Message) {
	userId := 1
	const (
		SERVER_HOST = "localhost"
		SERVER_PORT = "9988"
		SERVER_TYPE = "tcp"
	)

	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		clients = append(clients, connection)
		for _, v := range *messages {
			connection.Write([]byte(strconv.Itoa(v.Userid) + ": " + v.Content + "\n"))
		}
		go processClient(connection, messages, &userId)

	}
}

func broadcast(message *chat.Message, sender net.Conn) {
	for i, client := range clients {
		if client != sender {
			client.Write([]byte("\r\033[K" + strconv.Itoa(message.Userid) + ": " + message.Content + "\n" + strconv.Itoa(i+1) + "> "))
		}
	}
}

func processClient(connection net.Conn, messages *[]chat.Message, userId *int) {
	// connection.RemoteAddr().String()

	// TODO use IPv6 adress to fingerprint client and make a ID for the client based upon that IPv6 adress
	currentUserId := *userId
	*userId = *userId + 1
	reader := bufio.NewReader(connection)
	broadcast(chat.NewMessage(0, fmt.Sprintf("*** User %v connected ***", currentUserId)), connection)
	connection.Write([]byte(strconv.Itoa(currentUserId) + "> "))

	for {
		recieved, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			removeClient(connection)
			return
		}

		recieved = strings.TrimSpace(recieved)
		newMessage := *chat.NewMessage(currentUserId, recieved)
		*messages = append(*messages, newMessage)
		fmt.Println("Received:", recieved)
		broadcast(&newMessage, connection)
		connection.Write([]byte(strconv.Itoa(currentUserId) + "> "))
	}
}

func removeClient(connection net.Conn) {
	for i, client := range clients {
		if client == connection {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	connection.Close()
}
