// Create an Append Only File (AOF) for persistence of the Redis Store
package main

import (
	"fmt"
	"os"
	"sync"
)

// Aof manages the Append-Only File persistence.
type Aof struct {
	file *os.File // File handle for AOF
	mu   sync.Mutex // Mutex for file writes
}

// NewAof initializes the AOF persistence layer.
// It opens or creates the specified AOF file.
func NewAof(path string) (*Aof, error) {
	// Open the file with flags: Read/Write, Create if not exists, Append on writes
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening AOF file %s: %w", path, err)
	}

	a := &Aof{
		file: f,
	}

	// TODO: Add logic here later to load existing data from the AOF file

	fmt.Println("AOF initialized, using file:", path)
	return a, nil
}

// Close closes the AOF file.
func (a *Aof) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.file.Close()
}

// Write appends the given command Value (expected to be an Array) to the AOF file.
func (a *Aof) Write(value Value) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if value.typ != "array" {
		// Should not happen if called correctly, but good to check
		return fmt.Errorf("AOF Write expects an array Value, got %s", value.typ)
	}

	// Marshal the Value back to RESP format
	respBytes := marshalValue(value)

	// Write the RESP bytes to the file
	_, err := a.file.Write(respBytes)
	if err != nil {
		return fmt.Errorf("error writing to AOF file: %w", err)
	}

	// TODO: Consider adding fsync periodically for durability

	return nil
}

// marshalValue converts a Value back into its RESP byte representation.
// Note: This is a simplified version, might need more robust handling
// especially for nested arrays if those become supported.
func marshalValue(v Value) []byte {
	switch v.typ {
	case "simple_string":
		return []byte(fmt.Sprintf("+%s\r\n", v.str))
	case "error":
		return []byte(fmt.Sprintf("-%s\r\n", v.str))
	case "integer":
		// Assuming integer stored as string in Value struct based on current resp.go -- Correction: v.num is int
		return []byte(fmt.Sprintf(":%d\r\n", v.num)) // Use Sprintf with %d for integer
	case "bulk":
		if v.bulk == "" { // Handle potential empty bulk string if needed
			return []byte("$0\r\n\r\n")
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.bulk), v.bulk))
	case "null":
		return []byte("$-1\r\n")
	case "array":
		resp := fmt.Sprintf("*%d\r\n", len(v.array))
		for _, elem := range v.array {
			resp += string(marshalValue(elem)) // Recursively marshal elements
		}
		return []byte(resp)
	default:
		// Should not happen
		return []byte{}
	}
}

// TODO: Implement loading existing data from the AOF file.