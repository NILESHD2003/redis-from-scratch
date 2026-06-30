# Redis From Scratch

A Redis-compatible server built from scratch in Go to deeply understand how Redis works internally—from TCP networking and the RESP protocol to command execution, metadata, and the architecture behind an in-memory database.

This project is focused on learning by implementation rather than simply replicating features. Every subsystem is built from scratch to understand the design decisions behind Redis.

> **Current Status:** The project includes a concurrent TCP server, RESP2 encoder/decoder, command registry, execution engine, Redis-style command metadata, `COMMAND DOCS` support, and a fully functional `PING` command compatible with `redis-cli`.

---

# Table of Contents

- Overview
- Features
- Architecture
- Project Structure
- Server Workflow
- RESP2 Implementation
- Command Execution Engine
- Command Registry & Metadata
- Supported Commands
- Getting Started
- Testing
- Design Principles
- Roadmap

---

# Overview

Redis is much more than a key-value store.

Internally it consists of multiple independent systems working together:

- TCP networking
- RESP protocol parsing
- Command lookup
- Argument validation
- Command execution
- Response serialization

This project rebuilds those systems one at a time to understand how Redis processes requests from the moment a client connects until a response is written back.

Current implementation includes:

- Concurrent TCP server
- RESP2 encoder
- RESP2 decoder
- Command execution engine
- Redis-style command registry
- Command documentation metadata
- `COMMAND DOCS` support
- `PING` command
- Extensive unit tests

---

# Features

## Networking

- TCP server built using Go's standard library
- Concurrent client handling using goroutines
- Multiple client support
- Atomic connection tracking

## RESP2

Implemented from scratch:

### Decoder

- Simple Strings
- Errors
- Integers
- Bulk Strings
- Arrays
- Nested Arrays
- Null Bulk Strings
- Null Arrays

### Encoder

- Simple Strings
- Error Strings
- Integers
- Bulk Strings
- Arrays
- Null Bulk Strings
- Null Arrays

---

## Command System

- Redis-style command registry
- Execution engine
- Arity validation
- Command metadata
- `COMMAND DOCS`
- RESP-aware handler return types

---

## Testing

- Table-driven tests
- Decoder tests
- Encoder tests
- Invalid protocol tests
- Nested RESP parsing tests

---

# Architecture

```text
                +----------------+
                | Redis Client   |
                +----------------+
                        |
                        |
                 TCP Connection
                        |
                        ▼
                +----------------+
                | TCP Server     |
                +----------------+
                        |
                        ▼
                +----------------+
                | RESP Decoder   |
                +----------------+
                        |
                        ▼
                +----------------+
                | Execution      |
                | Engine         |
                +----------------+
                        |
                +-------+--------+
                |                |
                ▼                ▼
       Command Registry     Argument Validation
                |
                ▼
        Command Handler
                |
                ▼
        RESPValue Response
                |
                ▼
        RESP Encoder
                |
                ▼
            TCP Socket
                |
                ▼
          Redis Client
```

Each subsystem has a single responsibility which keeps the implementation modular and easy to extend.

---

# Server Workflow

For every connected client, the server continuously performs the following steps:

1. Read incoming bytes
2. Decode RESP request
3. Extract command and arguments
4. Lookup command in registry
5. Validate arguments
6. Execute handler
7. Encode returned RESP value
8. Write response back to client

Unlike simple echo servers, the networking layer never knows anything about Redis commands. Its only responsibility is transporting RESP messages.

---

# RESP2 Implementation

Redis communicates using RESP (Redis Serialization Protocol).

This project implements RESP2 completely from scratch.

Supported types:

| Type | Supported |
|-------|-----------|
| Simple String | ✅ |
| Error | ✅ |
| Integer | ✅ |
| Bulk String | ✅ |
| Array | ✅ |
| Nested Arrays | ✅ |
| Null Bulk String | ✅ |
| Null Array | ✅ |

Example request:

```text
*2
$4
PING
$5
hello
```

Decoded into:

```go
[]string{
    "PING",
    "hello",
}
```

Instead of returning raw strings, handlers return protocol-aware values which are encoded later by the RESP encoder.

---

# Command Execution Engine

The execution engine is responsible for everything after RESP decoding.

Responsibilities:

- Lookup commands
- Validate arity
- Execute handlers
- Return RESP values
- Convert errors into RESP error responses

Execution flow:

```text
RESP Request
      │
      ▼
Execution Engine
      │
      ▼
Lookup Command
      │
      ▼
Validate Arguments
      │
      ▼
Execute Handler
      │
      ▼
RESPValue
      │
      ▼
RESP Encoder
```

Keeping this layer separate allows networking and protocol parsing to remain completely independent of Redis commands.

---

# RESP-aware Handler Design

Instead of returning strings, handlers return RESP types.

Example:

```go
func HandlePING(args []string) (core.RESPValue, error) {
    if len(args) == 0 {
        return core.SimpleString("PONG"), nil
    }

    return core.BulkString(args[0]), nil
}
```

This design allows handlers to naturally return:

- Simple Strings
- Bulk Strings
- Integers
- Arrays
- Errors
- Null values

without the execution engine needing to infer the response type.

---

# Command Registry

Commands are stored in a centralized registry.

Each command contains:

```go
type CommandDefinition struct {
    Name    string
    Arity   int
    Handler Handler
    Docs    CommandDocs
}
```

The execution engine performs a lookup in this registry before executing any command.

This design makes adding new commands straightforward and keeps networking code completely generic.

---

# Command Metadata

Each command stores documentation similar to Redis itself.

Metadata includes:

- Summary
- Complexity
- Since Version
- Command Group
- Arguments
- Key Specifications
- Tips

This metadata powers Redis-compatible introspection commands such as:

```text
COMMAND DOCS
```

rather than maintaining documentation separately.

---

# Supported Commands

Currently implemented:

| Command | Status |
|----------|--------|
| PING | ✅ |
| COMMAND DOCS | ✅ |

More Redis commands will be implemented incrementally.

---

# Getting Started

Clone the repository:

```bash
git clone https://github.com/NILESHD2003/redis-from-scratch.git
cd redis-from-scratch
```

Run the server:

```bash
go run .
```

Using redis-cli:

```bash
redis-cli -p 6380 PING
```

Output:

```text
PONG
```

You can also inspect command metadata:

```bash
redis-cli -p 6380 COMMAND DOCS PING
```

---

# Testing

Run all tests:

```bash
go test ./...
```

Run RESP tests:

```bash
go test ./core -v
```

Current test coverage includes:

- RESP decoding
- RESP encoding
- Invalid inputs
- Missing CRLF
- Nested arrays
- Bulk strings
- Null values
- Error handling

---

# Design Principles

The goal of this project is not just feature parity with Redis, but understanding the engineering behind it.

Key principles include:

- Build every subsystem from scratch
- Keep networking independent of command execution
- Separate protocol parsing from business logic
- Prefer small, testable components
- Follow Redis architecture wherever practical
- Learn through implementation rather than abstraction

---

# Roadmap

## Foundation

- ✅ Concurrent TCP Server
- ✅ RESP2 Decoder
- ✅ RESP2 Encoder
- ✅ Command Registry
- ✅ Execution Engine
- ✅ RESP-aware handlers
- ✅ Command Metadata
- ✅ COMMAND DOCS
- ✅ PING

## Storage Engine

- ⬜ In-memory key-value store
- ⬜ GET
- ⬜ SET
- ⬜ DEL
- ⬜ EXISTS
- ⬜ KEYS

## Expiration

- ⬜ TTL
- ⬜ EXPIRE
- ⬜ PTTL

## Persistence

- ⬜ RDB Snapshots
- ⬜ Append Only File (AOF)

## Advanced Features

- ⬜ Transactions
- ⬜ Pub/Sub
- ⬜ Replication
- ⬜ Streams
- ⬜ Lua Scripting

## Engineering

- ⬜ Streaming RESP parser
- ⬜ Graceful shutdown
- ⬜ Benchmarks
- ⬜ Integration tests
- ⬜ Observability
- ⬜ Performance optimizations

---

# Educational Goals

This project is intended to provide a practical understanding of:

- TCP networking in Go
- Concurrent server design
- Wire protocol implementation
- Redis internals
- Command dispatch systems
- Protocol serialization
- Software architecture
- Test-driven development

Rather than treating Redis as a black box, the objective is to understand how every component works internally by implementing it from scratch.
