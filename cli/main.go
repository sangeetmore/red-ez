package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const serverAddress = "localhost:6969"

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", serverAddress, err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to redez server at %s\n", serverAddress)
	fmt.Println("Type 'QUIT' or 'EXIT' to close the client.")

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn) // Reader for server responses

	for {
		fmt.Printf("%s> ", serverAddress)

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue // Skip empty input
		}

		// Check for exit commands
		inputUpper := strings.ToUpper(input)
		if inputUpper == "QUIT" || inputUpper == "EXIT" {
			fmt.Println("Exiting client.")
			break
		}

		// Format command as RESP array of bulk strings
		parts := strings.Fields(input)
		respCmd := fmt.Sprintf("*%d\r\n", len(parts))
		for _, part := range parts {
			respCmd += fmt.Sprintf("$%d\r\n%s\r\n", len(part), part)
		}

		// Send command to server
		_, err = conn.Write([]byte(respCmd))
		if err != nil {
			fmt.Println("Error sending command:", err)
			break // Assume connection is broken
		}

		// Read response from server using RESP parser
		responseStr, err := readRespResponse(serverReader)
		if err != nil {
			fmt.Println("\nError reading response:", err)
			break // Assume connection is broken
		}

		// Print parsed response
		fmt.Println(responseStr)

		// TODO: Implement proper RESP response parsing for complex types (arrays, etc.)
		// This simple ReadString might not capture multi-line responses correctly.
	}
}

// readRespResponse reads and parses a single RESP response from the reader.
func readRespResponse(reader *bufio.Reader) (string, error) {
	typeByte, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	switch typeByte {
	case '+': // Simple String
		line, err := readLine(reader)
		return line, err
	case '-': // Error
		line, err := readLine(reader)
		return fmt.Sprintf("(error) %s", line), err
	case ':': // Integer
		line, err := readLine(reader)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(integer) %s", line), nil
	case '$': // Bulk String
		lengthStr, err := readLine(reader)
		if err != nil {
			return "", err
		}
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", fmt.Errorf("invalid bulk string length: %s", lengthStr)
		}

		if length == -1 {
			return "(nil)", nil // Null Bulk String
		}

		// Read the bulk string content + \r\n
		buf := make([]byte, length+2)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return "", err
		}
		// Return the string content without the trailing \r\n
		return fmt.Sprintf("\"%s\"", string(buf[:length])), nil
	case '*': // Array (basic handling - just shows count for now)
		countStr, err := readLine(reader)
		if err != nil {
			return "", err
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return "", fmt.Errorf("invalid array count: %s", countStr)
		}
		// TODO: Recursively parse array elements
		return fmt.Sprintf("(array, %d elements - parsing not fully implemented)", count), nil
	default:
		return "", fmt.Errorf("unknown RESP type prefix: %c", typeByte)
	}
}

// readLine reads a line ending in \r\n from the reader.
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// Trim the trailing \r\n
	return strings.TrimRight(line, "\r\n"), nil
}
