package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients = make(map[net.Conn]string)
	mutex   = &sync.Mutex{} // Mutex to synchronize access to messageChan
)

func main() {
	port := "8989"
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) != 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
	Addr := ":" + port

	// Start listening on the specified port
	listener, err := net.Listen("tcp", Addr)
	if err != nil {
		fmt.Println("Error starting TCP listener:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Server is listening on port %s...\n", port)
	// Loop to accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue // continue the loop even if there's an error
		}
		// Handle each connection in a separate goroutine
		go handleConnection(conn)
	}
}

// Function to read from the connection and process client messages
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Get the client's name
	conn.Write([]byte("Enter Your Name: "))
	nameBuffer := make([]byte, 1024)
	n, err := conn.Read(nameBuffer)
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}

	name := string(nameBuffer[:n-1]) // Remove newline character
	if !checkName(conn, name) {
		return // Exit if the name is invalid or taken
	}

	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()
	fmt.Printf("%s connected.\n", name)
	for {
		// Read from the connection
		n, err := conn.Read([]byte(nameBuffer))
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%s disconnected", name)
			} else {
				fmt.Println("Error reading from connection:", err)
			}
			return // Exit the loop if there's an error or client disconnects
		}
		// Process the received message
		BroadcastLoop(string(nameBuffer[:n]), conn)
	}
}

func BroadcastLoop(msg string, sender net.Conn) {
	for conn := range clients {
		if conn == sender {
			continue
		}
		_, err := conn.Write([]byte(clients[sender] + " : " + msg))
		if err != nil {
			fmt.Println("Error sending message:", err)
			conn.Close()
			delete(clients, conn)
		}

	}
}

func checkName(sender net.Conn, name string) bool {
	// Trim whitespace and check for empty name
	name = strings.TrimSpace(name)
	if name == "" {
		sender.Write([]byte("Name cannot be empty!\n"))
		return false
	}

	// Check for invalid characters
	for _, n := range name {
		if n < 32 || n > 126 {
			sender.Write([]byte("Name contains invalid characters!\n"))
			return false
		}
	}

	// Check if the name is already taken
	mutex.Lock()
	defer mutex.Unlock()
	for _, existingName := range clients {
		if existingName == name {
			sender.Write([]byte("Username already taken!\n"))
			return false
		}
	}

	return true
}
