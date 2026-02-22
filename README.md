# TITAN C2 FRAMEWORK

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.21-blue)
![License](https://img.shields.io/badge/license-MIT-green)

**TITAN** is an advanced, production-grade Command & Control (C2) framework designed for Red Team operations. It features a high-performance gRPC-based communication layer, a modern React-based dashboard, and a modular architecture for rapid extension.

## 🚀 Features

- **High Performance Transport**: Uses **gRPC/Protobuf** for efficient, low-latency bidirectional communication between agents and server.
- **Secure Communication**: Built-in support for TLS and AES-256 encryption (configurable).
- **Cross-Platform Agents**: Written in Go, agents can be compiled for Windows, Linux, and macOS.
- **Modern Web Dashboard**: Real-time React UI for managing agents, tracking status, and issuing commands.
- **Shell Execution**: Native shell integration (PowerShell on Windows, /bin/sh on *nix).
- **Persistent Job Queue**: Commands are queued and retrieved by agents during heartbeat cycles.

## 🏗 Architecture

The framework is divided into three main components:

1.  **C2 Server (`cmd/server`)**: The core controller. Manages state in memory (or DB), handles gRPC connections from agents, and serves the HTTP REST API + Web UI.
2.  **Titan Agent (`cmd/agent`)**: The implant. Runs on target machines, registers with the C2, and executes instructions.
3.  **Web Dashboard (`web/`)**: A React Single Page Application (SPA) for operators.

### Directory Structure

```
titan-c2-framework/
├── cmd/
│   ├── agent/          # Agent source code
│   └── server/         # Server source code
├── internal/
│   ├── db/             # In-memory store (state management)
│   └── server/         # gRPC server implementation
├── pkg/
│   ├── crypto/         # Encryption helpers (AES-256)
│   ├── pb/             # Protocol Buffer definitions
│   └── transport/      # gRPC connection wrappers
├── web/                # React Frontend
├── Dockerfile          # Production container build
└── Makefile            # Build automation
```

## 🛠 Usage

### Prerequisites
- Go 1.21+
- Node.js 18+ (for UI)
- Protoc (Protocol Buffers Compiler)

### Building from Source

1.  **Generate Protobufs**:
    ```bash
    make proto
    ```

2.  **Build Server & Agent**:
    ```bash
    make all
    ```
    Binaries will be placed in `bin/`.

3.  **Build Web UI**:
    ```bash
    cd web && npm install && npm run build
    ```

### Running via Docker

The easiest way to deploy the server is via Docker:

```bash
docker build -t titan-c2 .
docker run -p 8080:8080 -p 9090:9090 titan-c2
```

Access the dashboard at `http://localhost:8080`.

## 🛡 Security Note

This tool is created for **authorized Red Team engagements and educational purposes only**. The authors are not responsible for misuse.

## 🤝 Contributing

Contributions are welcome. Please ensure all pull requests pass the linter and include relevant tests.

## 🗺 Roadmap

- [x] **Phase 1**: Core C2 architecture with gRPC.
- [x] **Phase 2**: Secure transport layer (mTLS + AES-GCM) and rate-limiting.
- [ ] **Phase 3**: Steganographic listener endpoints.
- [ ] **Phase 4**: Decentralized peer-to-peer relay system.
- [ ] **Phase 5**: Multi-tenant UI with RBAC.

## 📄 License

MIT

## ✍️ Author

Olivier Robert-Duboille
