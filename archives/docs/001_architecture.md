# P2P Network Architecture

## Overview

This document describes the architecture of a peer-to-peer (P2P) network implementation in Go. The network supports both server nodes (full nodes) and clients, with a distributed public key registry system for peer discovery and routing.

## Network Components

### 1. Node Types

#### Server Nodes
- Run on dedicated servers (Ubuntu, macOS)
- Maintain connections to multiple peers
- Store and forward the distributed public key registry
- Handle routing between peers
- Participate in network governance
- Run continuously

#### Client Nodes
- Run on various platforms (Windows, macOS, Linux, Android, iOS)
- Connect to one or more server nodes
- Can initiate and receive connections
- May have intermittent connectivity
- Limited participation in network operations

### 2. Identity System

#### Public Key Registry
- Each node/client has a unique public/private key pair
- Public keys are encoded in base58 format for human readability
- Format: `key_[base58-encoded-string]`
- Registry is distributed across server nodes
- Contains:
  - Public key
  - Last known connection endpoints
  - Node type (server/client)
  - Last seen timestamp
  - Optional metadata

### 3. Connection Flow

1. **Initial Connection**
   - New node generates key pair
   - Connects to known server node(s)
   - Shares public key and node type
   - Receives partial registry of other nodes

2. **Peer Discovery**
   - Node queries connected servers for specific public keys
   - Servers respond with last known connection paths
   - Multiple paths may be available for redundancy

3. **Connection Establishment**
   - Node A wants to connect to Node B
   - Node A queries network for Node B's public key
   - Network returns possible routes to Node B
   - Connection is established through intermediate nodes
   - End-to-end encryption is used for security

### 4. Routing System

#### Path Finding
- Uses distributed hash table (DHT) for efficient lookups
- Multiple paths may exist between nodes
- Paths are ranked by:
  - Number of hops
  - Connection latency
  - Node reliability

#### Traffic Routing
- Messages are encrypted end-to-end
- Intermediate nodes only see next hop
- Support for different traffic types:
  - Direct messages
  - Broadcast messages
  - Service discovery
  - Registry updates

## Security Considerations

### Authentication
- All connections require public key authentication
- Challenge-response protocol using key pairs
- No central authority required

### Encryption
- All traffic is encrypted end-to-end
- TLS for node-to-node connections
- Custom protocol for client-to-client encryption

### Privacy
- Nodes can operate anonymously
- Only public keys are shared
- Intermediate nodes cannot read message content
- Optional metadata can be encrypted

## Implementation Details

### Technology Stack
- Language: Go (golang)
- Key cryptographic libraries:
  - `crypto/ed25519` for signatures
  - `x/crypto/nacl` for encryption
- Network stack:
  - TCP for reliable connections
  - UDP for peer discovery
  - Optional WebSocket support for web clients

### Data Structures

```go
type Node struct {
  PublicKey string
  NodeType NodeType
  LastSeen time.Time
  Endpoints []string
  Metadata map[string]string
}
type Route struct {
  Source string
  Destination string
  Hops []string
  Metrics RouteMetrics
}
```
