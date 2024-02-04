# Word of Wisdom TCP Server

## Disclaimer

This project is a homework assignment and, as such, should be considered a learning exercise rather than a
production-ready solution.
The implementation, especially the Proof of Work (PoW) component, is conceptual and may not align with real-world best
practices in cryptography or network security.
Key points to note:

1. **PoW Resource Consumption**: The PoW mechanism implemented here may consume more resources than the request handler
   itself.
   This approach was chosen to fulfill the assignment's requirements and demonstrate the concept.
   In a real-world scenario, the resource consumption and efficiency of such a system would need careful evaluation and
   likely a different approach.

2. **Level of Abstraction and Complexity**: I've tried to keep a balance between the complexity needed for a
   production-ready application and a minimal solution that gets the job done.
   The result is what I might have developed as a proof-of-concept for a new service.
   It includes the basic code layout structure, an example of unit testing, fuzz testing, things like graceful shutdown, and basic
   build tooling.

## Table of Contents

1. [Overview](#overview)
2. [Project Structure](#project-structure)
3. [Running Instructions](#running-instructions)
4. [Protocol Description](#protocol-description)
5. [Database File Format](#database-file-format)

## Overview

This project implements a TCP server that provides quotes from a specified collection.
To protect against DDoS attacks, the server uses a simple Proof of Work (PoW) challenge-response mechanism.
Clients must solve a PoW challenge before receiving a quote.

## Project Structure

- `cmd`: This directory contains the entry points for the server and client applications.
- `pkg`: Houses the core packages used by applications.

### Dockerfiles

Dockerfiles for building container images are located in their respective subdirectories under `cmd`.

### Other Files

- `db.txt`: A text file serving as a simple database for storing quotes. This file is used by the server to retrieve
  quotes.

## Running Instructions

Provided examples should be run from the project root.

### Server Usage

```shell
go run ./cmd/wow -h
Usage of wow:
  -addr string
    	application network address (default ":1111")
  -db string
    	database file (default "db.txt")
  -difficulty int
    	challenge difficulty (default 4)
  -timeout duration
    	connection idle timeout (default 3s)
```

### Client Usage

Client provides an interface similar to `ab`, you can specify the total number of requests and a concurrency level
using `-n` and `-c` flags.
To print the quotes to standard output, use `-v` flag.

```shell
go run ./cmd/client -h
Usage of client:
  -addr string
    	server address (default "localhost:1111")
  -c int
    	request concurrency (default 1)
  -n int
    	request number (default 1)
  -v	print quotes
```

### Local Environment

#### Prerequisites

- Go (version 1.21 or later)

#### Run server

```shell
go run ./cmd/wow
```

#### Run client

```shell
go run ./cmd/client -v -n 5 -c 3
```

### Docker Environment

Images should be built from the root directory of the project.

#### Prerequisites

- Docker

#### Build server

```shell
docker build -f cmd/wow/Dockerfile -t wow .
```

#### Build client

```shell
docker build -f cmd/client/Dockerfile -t client .
```

#### Run server

```shell
docker run -p 1111:1111 wow
```

#### Run client

```shell
docker run --net=host client -v -n 3 -c 3
```

## Protocol Description

### Connection Establishment

The client initiates the interaction by establishing a TCP connection to the server. The established connection may be
reused to issue subsequent requests.

### Request Command

Upon establishing a connection, the client sends a `get` command to the server to request a quote.

### PoW Challenge-Response

1. The server responds to the `get` command by sending a PoW challenge. The challenge comprises a difficulty level and a
   random string, separated by a colon `:`.
2. The client is required to compute a solution to the PoW challenge by finding a nonce. This nonce, when concatenated
   with the challenge string and hashed using SHA-256, should produce a hash with the specified number of leading zeros.
3. The client then sends the computed nonce back to the server as the solution.

### Quote Response

1. Upon receiving the nonce, the server verifies the client's PoW solution.
2. If the solution is valid, the server selects a quote from its database using a round-robin method and sends it to the
   client.
3. The server indicates the end of the quote by sending an additional newline character, resulting in an empty line.

### Connection Termination

The server may terminate the connection if the client's PoW solution is incorrect or after a specified idle timeout
period. The client can also close the connection.

### Example Interaction

1. The client connects to the server.
2. The client issues a `get` command.
3. The server sends a PoW challenge: `"4:abc123"`.
4. The client calculates and returns the solution: `"56789"`.
5. The server verifies the solution and sends a quote followed by an empty
   line: `"hello, world\n--Brian Kernighan, Programming in C: A Tutorial\n\n"`.
6. The client closes the connection.

## Database File Format

The database file is a simple text file.

### Structure

- Each quote in the database file is separated by a blank line.
- A single quote can span multiple lines, allowing for longer quotes.
- The end of a quote is marked by a newline character, followed by a blank line.
- The file ends with a blank line following the last quote.

### Example

Here is an example illustrating the format of the db file:

```
The only true wisdom is in knowing you know nothing.
-- Socrates

The unexamined life is not worth living.
-- Socrates

To know, is to know that you know nothing. That is the meaning of true knowledge.
-- Socrates

```
