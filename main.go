package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"netcat/src" 
)

// Truncate the prevMessages file and prepare it for a new chat
func truncatePrevMessages() {
	file, err := os.OpenFile("data/prevMessages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	err = file.Truncate(0) // Truncate the file to clear previous messages
	if err != nil {
		fmt.Println("Error truncating file:", err)
		return
	}
}

// Save new chat start time in the logs
func saveChatStartTime() {
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	src.SaveToFile("data/logs.txt", "------------------------new chat started at ["+formattedTime+"]--------------------------------\n\n\n")
}

// Start the TCP server
func startTCPServer(port string) (net.Listener, error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return nil, err
	}
	fmt.Println("Listening on the Port:", port[1:])
	return listener, nil
}

// Accept connections from clients
func acceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go src.HandleClient(conn) // Handle each client in a new goroutine
	}
}

// Periodically list connected clients based on console input
func listClientsOnCommand() {
	for {
		var input string
		fmt.Scanln(&input)
		if input == "list" {
			src.ListClients() // List clients when the "list" command is entered
		}
	}
}

func main() {
	port := src.CheckPort()

	truncatePrevMessages()
	saveChatStartTime()

	// Start the TCP server
	listener, err := startTCPServer(port)
	if err != nil {
		return
	}
	defer listener.Close()

	go acceptConnections(listener)

	listClientsOnCommand()
}
