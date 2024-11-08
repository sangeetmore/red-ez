// Init Server and Client requests to the server of the redez clone

package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	//Create a TCP server on port 6379
	l, err := net.Listen("tcp", ":6379")
	if err != nil { //print server error
		fmt.Println(err)
		return
	} else {
		fmt.Println("Listening on PORT 6379")
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

	for {
		buf := make([]byte, 1024)

		//read msgs from client

		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("err reading from client", err.Error())
			os.Exit(1)
		}
		//ignore req and send back PONG
		conn.Write([]byte("+HEHE\r\n"))
	}

}
