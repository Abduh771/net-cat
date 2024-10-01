# net-cat
### Introduction
A recreation of NetCat using Go, featuring a server-client architecture. The server listens for incoming connections, while clients can connect and transmit messages.

### Features
TCP connection for multiple clients (up to 10)
Clients must provide a name
Message broadcasting with timestamps and usernames
New clients receive previous messages
Notifications when clients join/leave
Default port: 8989
Installation
Ensure Go (v1.22.3+) is installed.
Clone the repository:
git clone https://learn.zone01oujda.ma/git/asoudri/net-cat.git

### Project Structure
```
├── go.mod                   # Go module file
├── main.go                  # Main server file
├── data/                    # Data files
│   ├── linuxLogo.txt        # ASCII art logo
│   ├── logs.txt             # Server chat log
│   └── prevMessages.txt     # Previous messages for new clients
└── src/                     # Source files
    ├── checkport.go         # Port checks
    ├── getmessages.go       # Fetches/saves messages
    └── handleclients.go     # Manages client interactions
```
