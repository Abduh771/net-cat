package src

import (
	"fmt"
	"net"
	"os"
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

func GenerateMessage(name string) string {
	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	return "[" + formattedTime + "]" + "[" + name + "]" + ":"
}

func WriteToClients(message string, clientAddr string, bl bool) {
    if bl {
        message = "\n" + GenerateMessage(clients[clientAddr].Name) + message + "\n" // Add newline at the end
        SaveToFile("data/prevMessages.txt", message[1:]) // Remove leading newline for saving
    } else {
        message = "\n" + clients[clientAddr].Name + message + "\n" // Add newline at the end
        SaveToFile("data/logs.txt", GenerateMessage("Client Name: "+clients[clientAddr].Name+" || Client Address "+clientAddr)+message[1:]) // Remove leading newline for saving
    }

    LoopAll(message, clientAddr)
}

func Status() {
	for i, j := range clients {
		if j.Name != "" {
			j.Conn.Write([]byte(GenerateMessage(clients[i].Name)))
		}
	}
}

func LoopAll(message, clientAddr string) {
	for i, j := range clients {
		if i != clientAddr {
			if j.Name != "" {
				j.Conn.Write([]byte(message))
			}
		}
	}
}

// ListClients displays all connected clients
func ListClients() {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	fmt.Println("Connected clients:")
	for addr := range clients {
		fmt.Println(addr)
	}
}

func SaveToFile(name, message string) {
	file, _ := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	file.WriteString(message)
}