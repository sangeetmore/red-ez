// Init Server and Client requests to the server of the redez clone

package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Global AOF instance
var globalAof *Aof

// Global in-memory data store and mutex
var (
	dataStore = make(map[string]string)
	storeMutex = sync.RWMutex{}
)

func main() {
	// Initialize AOF
	var err error
	globalAof, err = NewAof("redez.aof")
	if err != nil {
		fmt.Printf("Failed to initialize AOF: %v\n", err)
		// Server should exit if persistence fails
		return
	}
	// Ensure AOF file is closed when main function exits
	defer globalAof.Close()

	//Create a TCP server on port 6969
	l, err := net.Listen("tcp", ":6969")
	if err != nil { //print server error
		fmt.Println(err)
		return
	} else {
		fmt.Println("Listening on PORT 6969")
	}

	//Listen for connections to the server
	conn, err := l.Accept()
	if err != nil { //print conn error
		fmt.Println(err)
		return
	} else {
		fmt.Println("Server now accepting connections")
	}

	defer conn.Close() //close connection once it is finished

	// Create a RESP reader for the connection using the function from resp.go
	respReader := NewResp(conn)

	for {
		// Parse the command from the client using the method from resp.go
		value, err := respReader.Parse()
		if err != nil {
			// If client closes connection, EOF will be returned
			if err == io.EOF {
				fmt.Println("Client closed connection")
				break
			}
			fmt.Println("Error parsing command:", err)
			// Decide if we should continue or close connection on error
			// For now, let's just break the loop
			break 
		}

		// Print the parsed command for debugging
		fmt.Printf("Received: %+v\n", value)

		// Handle the command using the handler function (defined in handler.go)
		response := handleCommand(value)

		// Write the response back to the client
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error writing response:", err)
			// Close connection on write error
			break
		}
	}

}
