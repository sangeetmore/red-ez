package main

import (
	"fmt"
	"strings"
)

// handleCommand processes the parsed RESP command and returns the response bytes.
func handleCommand(value Value) []byte {
	if value.typ != "array" {
		return []byte("-ERR Command must be an array\r\n")
	}

	if len(value.array) == 0 {
		return []byte("-ERR Empty command array\r\n")
	}

	command := value.array[0]
	if command.typ != "bulk" {
		return []byte("-ERR Command must be a Bulk String\r\n")
	}

	cmdUpper := strings.ToUpper(command.bulk)
	args := value.array[1:] // Arguments are the rest of the array

	switch cmdUpper {
	case "PING":
		// PING command doesn't typically take arguments, but Redis accepts them.
		if len(args) == 0 {
			return []byte("PONG\r\n")
		} else if len(args) == 1 && args[0].typ == "bulk" {
			// Respond with the argument if PING has one argument
			return marshalBulkString(args[0].bulk)
		} else {
			return []byte("-ERR wrong number of arguments for 'ping' command\r\n")
		}
	case "ECHO":
		if len(args) != 1 || args[0].typ != "bulk" {
			return marshalError("ERR wrong number of arguments for 'echo' command")
		}
		return marshalBulkString(args[0].bulk)
	case "SET":
		if len(args) != 2 || args[0].typ != "bulk" || args[1].typ != "bulk" {
			return marshalError("ERR wrong number of arguments for 'set' command")
		}
		key := args[0].bulk
		val := args[1].bulk
		storeMutex.Lock()
		dataStore[key] = val
		storeMutex.Unlock()

		// Write command to AOF
		if err := globalAof.Write(value); err != nil {
			// Log the error, but maybe don't fail the operation for the client?
			// Or return a specific server error?
			fmt.Println("Failed to write to AOF:", err)
			// Potentially return an error to the client here
		}

		return marshalSimpleString("OK")
	case "GET":
		if len(args) != 1 || args[0].typ != "bulk" {
			return marshalError("ERR wrong number of arguments for 'get' command")
		}
		key := args[0].bulk
		storeMutex.RLock()
		value, ok := dataStore[key]
		storeMutex.RUnlock()
		if !ok {
			return nullBulkString // Key not found
		}
		return marshalBulkString(value)
	default:
		return []byte(fmt.Sprintf("-ERR unknown command `%s`\r\n", command.bulk))
	}
}

// marshalBulkString converts a string into RESP Bulk String format.
func marshalBulkString(str string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(str), str))
}

// marshalSimpleString converts a string into RESP Simple String format.
func marshalSimpleString(str string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", str))
}

// marshalError converts a string into RESP Error format.
func marshalError(str string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", str))
}

// Null Bulk String representation
var nullBulkString = []byte("$-1\r\n")
