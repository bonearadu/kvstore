# Key-Value Store

## Introduction

A multi-modal key-value store implementation in Go. Presents several implementation approaches, ranging from a
simple in-memory store to various persistent storage solutions.

This project is also used as a Go learning project, designed to explore the language's concurrency patterns, 
standard library, and idiomatic practices.

Key learning objectives included:
- Implementing concurrent data structures using Go's sync primitives
- Building HTTP servers with the standard library
- Exploring Go's generics for type-safe storage
- Developing clean, modular architecture
- Writing comprehensive tests

### Running the Server

The server supports several configuration flags:

```bash
go run main.go \
  -port 8080 \               # Port to listen on (default: 8080)
  -mode 0 \                  # Storage mode: 0=In-Memory (default), 1=Persistent, 2=Persistent with caching
  -store_path "" \           # Path for persistent storage (required for modes 1 and 2)
  -cache_capacity 100        # Cache size for persistent cached mode (default: 100)
```

Example for persistent storage with caching:
```bash
go run main.go -mode 2 -store_path ./data -cache_capacity 500
```

The server will start on the specified port and handle API requests according to the documented endpoints.

### Architecture

The project supports multiple storage implementations configurable via command-line flags:

0. In-Memory (default) - Fast volatile storage with no persistence between runs
1. Persistent - Disk-backed storage
2. Persistent with Caching - Disk-backed storage with in-memory LRU cache for optimal performance

The project is organized into the following components:

1. **Storage Layer** (`kv_store` package):
   - Defines a generic `KeyValueStore` interface
   - Provides the different key-value store implementations
   - Supports operations: Put, Get, Delete, and Entries

2. **Cache Layer** (`cache` package):
   - Defines a generic `Cache` interface
   - LRU (Least Recently Used) cache implementation
   - Configurable maximum capacity
   - Thread-safe concurrent access
   - Supports operations: Read, Write, Delete

3. **API Layer** (`api` package):
   - HTTP handlers for the RESTful API endpoints
   - Maps HTTP methods to storage operations
   - Handles request parsing and response formatting

4. **Server Layer** (`server` package):
   - Manages HTTP server lifecycle
   - Implements graceful shutdown
   - Handles signal processing

5. **Configuration Layer** (`config` package):
   - Parses command-line flags
   - Provides configuration options


### API

The key-value store exposes the following HTTP endpoints as a core API:

### Create or Update a Key-Value Pair

- **Endpoint**: `PUT /keys/{key}`
- **Description**: Store a value associated with the specified key. If the key already exists, its value will be updated.
- **Request Body**: Raw value content (any string/text)
- **Response**:
  - `201 Created` (for new keys)
  - `200 OK` (for updated keys)
  - `400 Bad Request` (if request is malformed)

### Retrieve a Value by Key

- **Endpoint**: `GET /keys/{key}`
- **Description**: Retrieve the value associated with the specified key.
- **Response**:
  - `200 OK` with value in response body
  - `404 Not Found` if key doesn't exist

### Delete a Key-Value Pair

- **Endpoint**: `DELETE /keys/{key}`
- **Description**: Remove the key-value pair for the specified key.
- **Response**:
  - `200 OK` on successful deletion
  - `404 Not Found` if key doesn't exist

### List All Entries

- **Endpoint**: `GET /keys`
- **Description**: Retrieve all key-value pairs in the store.
- **Response**:
  - `200 OK` with JSON object containing key-value pairs

### Implementation Details

- **Thread Safety**: The in-memory store uses `sync.RWMutex` to allow concurrent reads while ensuring exclusive access for writes.
- **Generic Types**: The store is implemented using Go's generics, allowing for type-safe storage of different key and value types.
- **RESTful Design**: The API follows RESTful principles with appropriate HTTP methods and status codes.
- **Extensibility**: The modular design makes it easy to add new features or replace components.

### Future improvements

Below are some ideas that can be considered to expand on our Key-Value store's capabilities:

- **More operations**: We can extend the InMemoryStore interface to support more operations which may streamline the
client's interactions. Examples: size, keySet, get multiple entries based on a regex.
- [DONE] **Persistent storage**: Persisting the storage is an important feature for a key-value store. This can be achieved
in a multitude of ways. One proposal would be to use a cache for frequently/recently accessed data, and a local database for keeping the state.
- **State snapshotting**: If we want to extend the in-memory solution, we can implement state snapshotting and operation logging.
On start, we can re-build the state starting from the latest snapshot and re-executing the operations that happened since, according to the logs.
This would take the solution in a completely different direction from the persistent storage one, making it more suitable
to smaller datasets requiring constant low-latency access.
- **A distributed design**: The Key-Value store problem is a popular one in system design. A larger extension would be
to implement our map store as a scalable distributed system, similar to Dynamo. This is a larger effort but a great
topic to explore at least from a theoretical standpoint.
