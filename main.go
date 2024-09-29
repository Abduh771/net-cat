package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"netcat/src"
)

func main() {
	port := src.CheckPort()

	// truncate the prev message file
	file, _ := os.OpenFile("data/prevMessages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.Truncate(0)
	// save start of new chat with time
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	src.SaveToFile("data/logs.txt", "------------------------new chat started at ["+formattedTime+"]--------------------------------\n\n\n")

	// Start TCP server
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on the Port: ", port[1:])

	// Accept connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}
			go src.HandleClient(conn)
		}
	}()

	// Periodically list connected clients (optional)
	for {
		var input string
		fmt.Scanln(&input)
		if input == "list" {
			src.ListClients() // List clients when the "list" command is entered
		}
	}
}