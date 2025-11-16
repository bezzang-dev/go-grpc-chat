# gRPC Chat Application

A simple real-time chat application built with Go and gRPC, featuring bidirectional streaming for instant messaging.

## Prerequisites

- Go 1.24 or higher
- Protocol Buffers compiler (`protoc`)

## Installation

```bash
# Clone the repository
git clone https://github.com/bezzang-dev/go-grpc-chat.git
cd go-grpc-chat

# Install dependencies
go mod download
```

## Usage

### Start the Server

```bash
go run /server/main.go -port 50051
```

### Start Clients

Open multiple terminals and run:

```bash
# Client 1
go run /client/main.go -id Alice -addr localhost:50051

# Client 2
go run /client/main.go -id Bob -addr localhost:50051

# Client 3
go run /client/main.go -id Charlie -addr localhost:50051
```

### Chat Commands

- Type your message and press Enter to send
- Type `exit` to disconnect

## How It Works

### Server
- Listens for incoming gRPC connections
- Maintains a list of connected client streams
- Broadcasts messages to all connected clients
- Thread-safe with mutex for concurrent operations

### Client
- Establishes bidirectional streaming connection
- Sends messages from user input (main goroutine)
- Receives messages from server (background goroutine)
- Displays incoming messages in real-time

## Example

```bash
# Terminal 1 - Server
$ go run /server/main.go
Server listening on port 50051

# Terminal 2 - Alice
$ go run /client/main.go -id Alice
Hello everyone!
Sender:Bob Message:Hi Alice!
Sender:Charlie Message:Welcome!

# Terminal 3 - Bob  
$ go run /client/main.go -id Bob
Sender:Alice Message:Hello everyone!
Hi Alice!
Sender:Charlie Message:Welcome!
```