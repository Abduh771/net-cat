package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	port := ":8989"
	if len(os.Args) == 2 {
		_, e := strconv.Atoi(os.Args[1])
		if e == nil {
			port = ":" + os.Args[1]
		} else {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// turncate the prev message file
	file, _ := os.OpenFile("prevMessages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	file.Truncate(0)
	///save start of new chat with time
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	SaveToFile("logs.txt", "------------------------new chat started at ["+formattedTime+"]--------------------------------\n\n\n")

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
			go handleClient(conn)
		}
	}()

	// Periodically list connected clients (optional)
	for {
		var input string
		fmt.Scanln(&input)
		if input == "list" {
			listClients() // List clients when the "list" command is entered
		}
	}
}
