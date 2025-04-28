# Redez - A Simple Redis Clone in Go

## Description

Redez is a basic implementation of a Redis-like in-memory key-value store built with Go. This project serves as an educational tool to understand the core concepts behind Redis, including the Redis Serialization Protocol (RESP), command handling, concurrency, and persistence mechanisms like AOF.

## Features

* **TCP Server:** Listens for client connections on port 6969. *(Note: The current `main.go` accepts only a single connection after startup; a loop is needed for true concurrent handling)*.
* **RESP Implementation:** Includes a parser for handling basic RESP data types: Arrays, Bulk Strings, Simple Strings, Errors, and Integers. Also includes functions for serializing responses back into RESP format.
* **Command Handling:** Supports a subset of Redis commands: `PING`, `ECHO`, `SET`, `GET`.
* **In-Memory Storage:** Utilizes a thread-safe `map[string]string` protected by a `sync.RWMutex` for storing key-value data.
* **AOF Persistence:** Implements Append-Only File (AOF) persistence. Successful `SET` commands are appended to the `redez.aof` file.
* **Basic CLI Client:** A command-line client is included (`./cli`) for interacting with the Redez server.

## Getting Started

### Prerequisites

* Go (latest stable version recommended)

### Building

1.  **Build & Run the Server:**
    ```bash
    cd /path/to/redez # Navigate to the project root
    make run
    ```
2.  **Build & Run the Client:**
    Open another terminal window:
    ```bash
    cd /path/to/redez # Navigate to the project root
    make run-cli
    ```

#### You can now type Redis commands (like `PING`, `SET key value`, `GET key`) into the client terminal. Type `QUIT` or `EXIT` to close the client.

## Implemented Commands

* `PING [message]` - Responds with `PONG` if no message is provided, or echoes the `message` if one is given.
* `ECHO message` - Responds with the provided `message`.
* `SET key value` - Stores the `key` and `value`. This command is logged to the AOF file. Responds with `OK` on success.
* `GET key` - Retrieves the value associated with the `key`. Responds with the value as a bulk string, or `(nil)` if the key does not exist.

## Future Work / TODOs

* Implement a proper connection loop in `main.go` to handle multiple concurrent clients.
* Implement the logic to load existing data from the AOF file upon server startup (`aof.go`).
* Expand command support (e.g., `DEL`, `EXISTS`, List commands, Hash commands).
* Add support for key expiration (`EX` option for `SET`, `TTL` command).
* Enhance error handling throughout the application.
* Improve the client's RESP response parser to handle complex types like nested arrays correctly (`cli/main.go`).
* Consider adding configuration options (e.g., port, AOF filename).
* Implement periodic `fsync` for the AOF file to improve data durability (`aof.go`).
