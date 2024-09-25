package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	Name string
	Conn net.Conn
}

var (
	clients      = make(map[string]Client) // Map to store connected clients
	clientsMutex sync.Mutex                // Mutex to synchronize access to clients map
)

func geneateMessage(name string) string {
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	return "[" + formattedTime + "]" + "[" + name + "]" + ":"
}

func writeToClients(message string, clientAddr string, bl bool) {
	if bl {
		message = "\n" + geneateMessage(clients[clientAddr].Name) + message
		SaveToFile("prevMessages.txt", message[1:])
	} else {
		message = "\n" + clients[clientAddr].Name + message
		SaveToFile("logs.txt", geneateMessage("Client Name: "+clients[clientAddr].Name+" || Client Adress "+clientAddr)+message[1:])
	}

	
	loopAll(message, clientAddr)
}

func Status() {
	for i, j := range clients {
		if j.Name != "" {
			j.Conn.Write([]byte(geneateMessage(clients[i].Name)))
		}
	}
}

// Function to handle each client connection

func loopAll(message, clientAddr string) {
	for i, j := range clients {
		if i != clientAddr {
			if j.Name != "" {
				j.Conn.Write([]byte(message))
			}
		}
	}
}

// Function to display all connected clients
func listClients() {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	fmt.Println("Connected clients:")
	for addr := range clients {
		fmt.Println(addr)
	}
}
