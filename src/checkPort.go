package src

import (
	"fmt"
	"os"
)

func atoi(s string) int {
	result := 0
	for _, c := range s {
		if c < '0' || c > '9' { 
			return 0 
		}
		result = result*10 + int(c-'0') // Convert character to int and build the number
	}
	return result
}

// CheckPort function to validate and set the port
func CheckPort() string {
	port := ":8989"       
	if len(os.Args) == 2 { 
		if atoi(os.Args[1]) != 0 { 
			port = ":" + os.Args[1] 
		} else {
			fmt.Println("[USAGE]: ./TCPChat $port") 
			os.Exit(1)
		}
	} else if len(os.Args) > 2 { 
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
	return port // Return the validated port
}