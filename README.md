# GoRAT: A Simple Remote Administration Tool

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.18%2B-blue?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Python-3.7%2B-blue?style=for-the-badge&logo=python" alt="Python Version">
  <img src="https://img.shields.io/badge/Status-Educational-orange?style=for-the-badge" alt="Project Status">
</p>

<p align="center">
  <i>A simple, cross-platform RAT developed in Go and Python for a short-term project.</i>
</p>

---


---

> **:warning: Disclaimer:** This project was developed for educational purposes only. It is a proof-of-concept to demonstrate how a Remote Administration Tool (RAT) can be built. Using such tools for unauthorized access to computer systems is illegal. The author is not responsible for any misuse of this code.

## Overview

GoRAT is a simple, cross-platform RAT developed as part of a short project. It consists of a client (implant) written in Go and a command-and-control (C2) server written in Python. The goal was to create a basic framework for remote command execution and data exfiltration in a limited time frame.

## Features

-   **Remote Shell:** Execute shell commands on the client machine.
-   **File Transfer:** Upload files to the client and download files from it.
-   **Screenshot:** Capture the client's screen.
-   **Persistence:** A basic persistence mechanism for Linux/macOS using `crontab`.
-   **Destructive Payload:** Includes a `fork_bomb` command as an example payload.

## Architecture

The project is split into two main components:

### Server (C2)

-   **Language:** Python 3
-   **File:** `server.py`
-   **Functionality:**
    -   Listens for incoming TCP connections from the client.
    -   Provides an interactive command-line interface for the operator.
    -   Sends commands to the client and receives/displays the output.
    -   Handles file transfers and saves received files (like screenshots and downloads).

### Client (Implant)
-   **Language:** Go
-   **File:** `main.go`
-   **Functionality:**
    -   Connects back to the hardcoded C2 server address.
    -   Continuously retries connection if the server is unavailable.
    -   Receives and parses commands from the server.
    -   Executes commands (e.g., shell commands, file operations, screenshot).
    -   Sends results back to the C2 server.

## Getting Started

### Prerequisites

-   Go compiler (version 1.18+ recommended)
-   Python 3 (version 3.7+ recommended)
-   Go dependencies:
    ```sh
    go get github.com/kbinani/screenshot
    ```

### 1. Server Setup

Run the Python server. It will start listening for connections on `127.0.0.1:1234`.

```sh
python3 server.py
[*] Server started. Listening for incoming connections...
```

### 2. Build and Run the Client

The C2 server address is hardcoded in `main.go`. Modify this if your server is running on a different IP or port.

```go
const c2 string = "127.0.0.1:1234"
```

Build the Go client:

```sh
go build -o client main.go
```

Running directly
```sh
go run main.go
```

Run the client on the target machine. Once it connects, you will see a message on the server terminal.

## Available Commands

You can issue commands from the server's `$` prompt.

| Command               | Description                                                              |
| --------------------- | ------------------------------------------------------------------------ |
| `<shell_command>`     | Executes any shell command (e.g., `whoami`, `ls -l`).                      |
| `cd <directory>`      | Changes the current working directory on the client.                     |
| `download <filepath>` | Downloads a file from the client to the server's directory.              |
| `upload <filepath>`   | Uploads a file from the server's directory to the client's CWD.          |
| `screenshot`          | Takes a screenshot of the client's primary display and saves it as `screenshot.png`. |
| `persist`             | Establishes persistence on the client via a cron job (`@reboot`).        |
| `fork_bomb`           | **(DANGEROUS)** Triggers a fork bomb on the client, likely crashing it.   |
| `q` / `quit`          | Closes the connection to the client and shuts down the server.           |
