package src

import (
	"bufio"
	"net"
	"os"
	"strings"
)

// HandleClient is the main handler for each connected client
func HandleClient(conn net.Conn) {
	defer conn.Close()

	if !isConnectionFull(conn) {
		return
	}

	clientAddr := conn.RemoteAddr().String()
	displayLinuxLogo(conn)

	name := requestClientName(conn, clientAddr)

	// Only write to clients if the name is valid
	if name != "" {
		WriteToClients(" has joined our chat...", clientAddr, false)
		clients[clientAddr].Conn.Write([]byte(prevMessage()))
	}

	handleMessages(conn, clientAddr, name)
}

func prevMessage() string {
	data, _ := os.ReadFile("data/prevMessages.txt")
	return string(data)
}

func isConnectionFull(conn net.Conn) bool {
	if len(clients) == 5 {
		conn.Write([]byte("Connection is Full in the server"))
		return false
	}
	return true
}

func displayLinuxLogo(conn net.Conn) {
	linuxLogo, _ := os.ReadFile("data/linuxLogo.txt")
	conn.Write(linuxLogo)
}

func requestClientName(conn net.Conn, clientAddr string) string {
	existingNames := make(map[string]bool) // Track used names
	var name string
	for {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		reader := bufio.NewReader(conn)
		inputName, err := reader.ReadString('\n')
		if err != nil {
			return ""
		}
		name = strings.TrimSpace(inputName)
		if !IsPrintable(name) {
			conn.Write([]byte("Invalid name: Please use only printable characters.\n"))
			continue
		}
		clientsMutex.Lock()
		if isValidName(name, existingNames) && !isNameTaken(name) {
			clients[clientAddr] = Client{Conn: conn, Name: name}
			existingNames[name] = true 
			clientsMutex.Unlock()
			break
		}
		clientsMutex.Unlock()
	}
	return name
}

func handleMessages(conn net.Conn, clientAddr, name string) {
    reader := bufio.NewReader(conn)
    showStatus := true
    for {
        if showStatus {
            Status()
        }
        showStatus = true

        message, err := reader.ReadString('\n')
        if err != nil {
            WriteToClients(" has left our chat...\n", clientAddr, false)
            Status()
            clientsMutex.Lock()
            delete(clients, clientAddr)
            clientsMutex.Unlock()
            break
        }

        message = strings.TrimSpace(message) // Trim whitespace from the message
        if !IsPrintable(message) {
            conn.Write([]byte("Invalid input: Please use only printable characters.\n"))
            continue
        }

        if len(message) == 0 {
            conn.Write([]byte(GenerateMessage(name))) // Show only the user's prompt
            showStatus = false
        } else {
            clientsMutex.Lock()
            WriteToClients(message, clientAddr, true) // Send message to other clients
            clientsMutex.Unlock()
        }
    }
}

// IsPrintable validates that the input contains only printable characters
func IsPrintable(input string) bool {
	for _, char := range input {
		if char < 32 || char > 126 {
			return false
		}
	}
	return true
}

func isValidName(name string, existingNames map[string]bool) bool {
	if len(name) == 0 || len(name) > 12 || !IsPrintable(name) || existingNames[name] {
		return false
	}
	return true
}

func isNameTaken(name string) bool {
	for _, client := range clients {
		if client.Name == name {
			return true
		}
	}
	return false
}