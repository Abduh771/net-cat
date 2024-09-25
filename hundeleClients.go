package main

import (
	"bufio"
	"net"
	"os"
	"strings"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	if len(clients) == 10 {
		conn.Write([]byte("Connection is Full in the server"))
		return
	}
	clientAddr := conn.RemoteAddr().String()

	// linuxLogo
	linuxLogo, _ := os.ReadFile("linuxLogo.txt")
	conn.Write(linuxLogo)

	for !isValidName(clients[clientAddr].Name) {
		// read name
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		reader := bufio.NewReader(conn)
		name, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		// name = name[:len(name)-1]
		name = strings.TrimSpace(name)
		clientsMutex.Lock()
		clients[clientAddr] = Client{Conn: conn, Name: name}
		clientsMutex.Unlock()
	}
	// fmt.Printf("Client connected: %s\n", clientAddr)
	reader := bufio.NewReader(conn)
	// enter message

	writeToClients(" has joined our chat...\n", clientAddr, false)

	clients[clientAddr].Conn.Write([]byte(prevMessage()))

	bl := true
	for {
		// check for empty message
		if bl {
			Status()
		}
		bl = true
		// Read data until newline or EOF (you can modify the delimiter if needed)

		message, err := reader.ReadString('\n')
		if err != nil {
			// leave message
			writeToClients(" has left our chat...\n", clientAddr, false)
			Status()
			delete(clients, clientAddr)
			// fmt.Printf("Client disconnected: %s\n", clientAddr)
			break // Exit loop if the client disconnects or an error occurs
		}

		// Process the received message (remove newline characters if necessary)

		// send message toclients

		if len(message) == 1 {
			// fmt.Printf("Received message from %s: empty message\n", clients[clientAddr].Name)
			conn.Write([]byte(geneateMessage(clients[clientAddr].Name)))
			bl = false
		} else {
			// fmt.Printf("Received message from %s: %s\n", clientAddr, message)
			clientsMutex.Lock()
			writeToClients(message, clientAddr, true)
			clientsMutex.Unlock()
		}
	}

	// Remove client from the map when they disconnect
	clientsMutex.Lock()
	delete(clients, clientAddr)
	clientsMutex.Unlock()
}
