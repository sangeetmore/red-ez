// Create a parser (serializer and deserializer) for the values and keys in the following syntax:
// Eg:=> $11\r\nDATATOSTORE\r\n
// a. '$' => first byte that represents the data type (here, it's a 'string')
// b. 11 => represents the no. of bytes needed to store the data
// c. \r\n => additional 2 bytes after the first 2 bytes, key and value; it's called CLRF & it represents the end of a line

//We will use bufio for the parser

package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Defining a const for first byte data type
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Struct used for serialization and deserialization process
type Value struct {
	typ   string  //type of the Value
	str   string  //holds string value of simple_strings
	num   int     //holds the integer value
	bulk  string  //holds the string value of bulk strings
	array []Value //array that can contain items of struct Value
}

// Type of reader used to read the data from the buffer
type Resp struct {
	reader *bufio.Reader
}

// New Resp reader to read the data from the buffer
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readLine reads a line from the reader until \r\n
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n++
		line = append(line, b)
		// Check if the line ends with \r\n
		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}
	// Return the line without the trailing \r\n
	return line[:len(line)-2], n, nil
}

// readInteger reads an integer value after the type prefix
func (r *Resp) readInteger() (int, error) {
	line, _, err := r.readLine()
	if err != nil {
		return 0, err
	}
	// Simple integer parsing for now
	i := 0
	_, err = fmt.Sscan(string(line), &i)
	if err != nil {
		return 0, fmt.Errorf("error parsing integer: %w", err)
	}
	return i, nil
}

// readBulk reads a bulk string value after the type prefix
func (r *Resp) readBulk() (string, error) {
	line, _, err := r.readLine()
	if err != nil {
		return "", err
	}
	// Read the bulk string length
	length, err := strconv.Atoi(string(line))
	if err != nil {
		return "", fmt.Errorf("error parsing bulk string length: %w", err)
	}
	// Read the bulk string value
	bulk := make([]byte, length)
	_, err = io.ReadFull(r.reader, bulk)
	if err != nil {
		return "", err
	}
	// Read the trailing \r\n
	_, _, err = r.readLine()
	if err != nil {
		return "", err
	}
	return string(bulk), nil
}

// readArray reads an array value after the type prefix
func (r *Resp) readArray() ([]Value, error) {
	line, _, err := r.readLine()
	if err != nil {
		return nil, err
	}
	// Read the array length
	length, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, fmt.Errorf("error parsing array length: %w", err)
	}
	array := make([]Value, length)
	for i := range array {
		value, err := r.Parse()
		if err != nil {
			return nil, err
		}
		array[i] = value
	}
	return array, nil
}

// Parse the data from the buffer
func (r *Resp) Parse() (Value, error) {
	// Read the first byte to determine the data type
	typeByte, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typeByte {
	case ARRAY:
		array, err := r.readArray()
		if err != nil {
			return Value{}, err
		}
		return Value{typ: "array", array: array}, nil
	case BULK:
		bulk, err := r.readBulk()
		if err != nil {
			return Value{}, err
		}
		return Value{typ: "bulk", bulk: bulk}, nil
	case STRING:
		line, _, err := r.readLine()
		if err != nil {
			return Value{}, err
		}
		return Value{typ: "string", str: string(line)}, nil
	case ERROR:
		line, _, err := r.readLine()
		if err != nil {
			return Value{}, err
		}
		return Value{typ: "error", str: string(line)}, nil
	case INTEGER:
		num, err := r.readInteger()
		if err != nil {
			return Value{}, err
		}
		return Value{typ: "integer", num: num}, nil
	default:
		return Value{}, fmt.Errorf("invalid RESP type byte: %c", typeByte)
	}
}
