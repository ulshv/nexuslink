package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
)

// Protocol identifier for our custom ping protocol
const pingProtocol = "/nexuslink.dev/ping/0.1.0"

func main() {
	// Parse command line flags
	isServer := flag.Bool("server", false, "Run as server (initiates connections)")
	listenAddr := flag.String("listen", "", "Address to listen on")
	targetAddr := flag.String("target", "", "Target peer to connect to (multiaddr)")
	debugLogs := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Set up logging
	logLevel := slog.LevelInfo
	if *debugLogs {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		slog.Info("Received interrupt signal, shutting down...")
		cancel()
	}()

	// Create libp2p host
	h, err := createHost(*listenAddr)
	if err != nil {
		slog.Error("Failed to create host", "error", err)
		os.Exit(1)
	}

	// Log host information
	slog.Info("Host created", "id", h.ID())

	// Display the peer ID address prominently for easy copying
	peerIDAddress := fmt.Sprintf("/p2p/%s", h.ID().String())
	slog.Info("═════════════════════════════════════════════")
	slog.Info("PEER ID ADDRESS (use this to connect to this node)", "address", peerIDAddress)
	slog.Info("═════════════════════════════════════════════")

	// Then continue with the existing code that logs all listening addresses
	for _, addr := range h.Addrs() {
		slog.Info("Listening on", "addr", fmt.Sprintf("%s/p2p/%s", addr, h.ID()))
	}

	// Set up ping protocol
	setupPingProtocol(h)

	// Set up DHT for peer discovery
	dht, err := setupDHT(ctx, h)
	if err != nil {
		slog.Error("Failed to set up DHT", "error", err)
		os.Exit(1)
	}

	// Start periodic announcements
	go announceToNetwork(ctx, h, dht)

	// If running as server and target is specified, connect to target
	if *isServer && *targetAddr != "" {
		go connectToTarget(ctx, h, *targetAddr, dht)
	}

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("Shutting down...")
}

// createHost creates a new libp2p host
func createHost(listenAddr string) (host.Host, error) {
	var opts []libp2p.Option

	// Create keys directory if it doesn't exist
	// This directory will store the node's private key for persistent identity
	if err := os.MkdirAll("./keys", 0700); err != nil {
		return nil, fmt.Errorf("failed to create keys directory: %w", err)
	}

	// Load or generate private key
	// This ensures the node has a consistent identity across restarts
	priv, err := loadOrGenerateKey("./keys/node.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load or generate key: %w", err)
	}
	opts = append(opts, libp2p.Identity(priv))

	// Handle different address formats
	if listenAddr == "" || strings.HasPrefix(listenAddr, "/p2p/") {
		// For empty address or peer ID only, listen on all interfaces with random ports
		slog.Info("Using random ports on all interfaces")
		opts = append(opts, libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
			"/ip6/::/tcp/0",
		))
	} else {
		// Use specified multiaddress
		addr, err := multiaddr.NewMultiaddr(listenAddr)
		if err != nil {
			return nil, fmt.Errorf("invalid listen address: %w", err)
		}
		opts = append(opts, libp2p.ListenAddrs(addr))
	}

	// Set up relay addresses
	relayAddrs := []peer.AddrInfo{}
	for _, addr := range dht.DefaultBootstrapPeers[:2] {
		pi, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			continue
		}
		relayAddrs = append(relayAddrs, *pi)
	}

	// Add NAT traversal options
	opts = append(opts,
		libp2p.EnableNATService(),
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
		libp2p.EnableAutoRelayWithStaticRelays(relayAddrs),
	)

	return libp2p.New(opts...)
}

// loadOrGenerateKey loads existing private key or creates a new one
func loadOrGenerateKey(keyPath string) (crypto.PrivKey, error) {
	// Check if key file exists
	if _, err := os.Stat(keyPath); err == nil {
		// Key file exists, load it
		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file: %w", err)
		}
		// Decode the base64 encoded key
		keyBytes, err := base64.StdEncoding.DecodeString(string(keyData))
		if err != nil {
			return nil, fmt.Errorf("failed to decode key: %w", err)
		}
		priv, err := crypto.UnmarshalPrivateKey(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal key: %w", err)
		}

		slog.Info("Loaded existing private key", "path", keyPath)
		return priv, nil
	}

	// Generate new Ed25519 key
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	keyBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}
	keyEncoded := base64.StdEncoding.EncodeToString(keyBytes)

	// Write the key to file with secure permissions (0600 = owner read/write only)
	// This prevents other users from accessing the private key
	if err := os.WriteFile(keyPath, []byte(keyEncoded), 0600); err != nil {
		return nil, fmt.Errorf("failed to write key file: %w", err)
	}

	slog.Info("Generated and stored new private key", "path", keyPath)
	return priv, nil
}

// setupPingProtocol registers handlers for ping protocol
func setupPingProtocol(h host.Host) {
	_ = ping.NewPingService(h)
	h.SetStreamHandler(protocol.ID(pingProtocol), handlePingStream)
}

// handlePingStream processes ping/pong messages
func handlePingStream(s network.Stream) {
	defer s.Close()

	buf := make([]byte, 4)
	remotePeer := s.Conn().RemotePeer()

	for {
		// Read from the stream
		n, err := s.Read(buf)
		if err != nil {
			// Don't log EOF as an error - it's a normal stream close
			if err.Error() != "EOF" {
				slog.Error("Error reading from stream", "error", err)
			}
			return
		}

		if n != 4 {
			continue
		}

		msg := string(buf)
		if msg == "ping" {
			slog.Info("Received ping", "from", remotePeer)
			if _, err := s.Write([]byte("pong")); err != nil {
				slog.Error("Error sending pong", "error", err)
				return
			}
		} else if msg == "pong" {
			slog.Info("Received pong", "from", remotePeer)
		}
	}
}

// setupDHT initializes the DHT for peer discovery
func setupDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {
	kdht, err := dht.New(ctx, h, dht.Mode(dht.ModeAuto))
	if err != nil {
		return nil, err
	}

	// Connect to bootstrap nodes in background
	go func() {
		for _, addr := range dht.DefaultBootstrapPeers {
			pi, err := peer.AddrInfoFromP2pAddr(addr)
			if err != nil || pi.ID == h.ID() {
				continue
			}

			connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			if err := h.Connect(connectCtx, *pi); err != nil {
				slog.Debug("Failed to connect to bootstrap node", "peer", pi.ID)
			}
			cancel()
		}

		bootstrapCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := kdht.Bootstrap(bootstrapCtx); err != nil {
			slog.Error("Failed to bootstrap DHT", "error", err)
		}
	}()

	return kdht, nil
}

// connectToTarget connects to a peer and starts sending pings
func connectToTarget(ctx context.Context, h host.Host, targetAddrStr string, dht *dht.IpfsDHT) {
	var peerID peer.ID
	var err error

	// Handle peer ID format
	if strings.HasPrefix(targetAddrStr, "/p2p/") {
		peerIDStr := strings.TrimPrefix(targetAddrStr, "/p2p/")
		peerID, err = peer.Decode(peerIDStr)
		if err != nil {
			slog.Error("Invalid peer ID", "error", err)
			return
		}

		// Try to find the peer with retries
		go func() {
			// Keep trying until context is canceled or connection succeeds
			retryTicker := time.NewTicker(5 * time.Second)
			defer retryTicker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-retryTicker.C:
					// Check if we're already connected
					if h.Network().Connectedness(peerID) == network.Connected {
						slog.Info("Already connected to peer", "peer", peerID)
						go sendPings(ctx, h, peerID)
						return
					}

					// Try direct connection first
					slog.Info("Attempting direct connection to peer", "peer", peerID)
					info := peer.AddrInfo{ID: peerID}
					if err := h.Connect(ctx, info); err == nil {
						slog.Info("Connected directly to peer", "peer", peerID)
						go sendPings(ctx, h, peerID)
						return
					}

					// If direct connection fails, try DHT lookup
					slog.Info("Direct connection failed, looking up peer in DHT", "peer", peerID)
					findCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
					peerInfo, err := dht.FindPeer(findCtx, peerID)
					cancel()

					if err != nil {
						slog.Warn("Failed to find peer in DHT, will retry", "peer", peerID, "error", err)
						continue
					}

					// Try to connect using the discovered addresses
					slog.Info("Found peer in DHT, connecting", "peer", peerID, "addrs", peerInfo.Addrs)
					if err := h.Connect(ctx, peerInfo); err != nil {
						slog.Warn("Failed to connect to peer, will retry", "peer", peerID, "error", err)
						continue
					}

					slog.Info("Connected to peer via DHT", "peer", peerID)
					go sendPings(ctx, h, peerID)
					return
				}
			}
		}()

		return
	}

	// Handle full multiaddress
	targetAddr, err := multiaddr.NewMultiaddr(targetAddrStr)
	if err != nil {
		slog.Error("Invalid target address", "error", err)
		return
	}

	info, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		slog.Error("Invalid peer address", "error", err)
		return
	}

	if err := h.Connect(ctx, *info); err != nil {
		slog.Error("Failed to connect to peer", "error", err)
		return
	}

	peerID = info.ID
	slog.Info("Connected to peer", "peer", peerID)
	go sendPings(ctx, h, peerID)
}

// sendPings sends periodic pings to a peer
func sendPings(ctx context.Context, h host.Host, peerID peer.ID) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Open stream with timeout
			streamCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			s, err := h.NewStream(streamCtx, peerID, protocol.ID(pingProtocol))
			cancel()

			if err != nil {
				slog.Error("Failed to open stream", "error", err)
				continue
			}

			// Set deadline
			if err := s.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				s.Close()
				continue
			}

			// Send ping
			slog.Info("Sending ping", "to", peerID)
			if _, err = s.Write([]byte("ping")); err != nil {
				slog.Error("Failed to send ping", "error", err)
				s.Close()
				continue
			}

			// Read response
			buf := make([]byte, 4)
			if _, err = s.Read(buf); err != nil {
				slog.Error("Failed to read pong", "error", err)
				s.Close()
				continue
			}

			if string(buf) == "pong" {
				slog.Info("Received pong", "from", peerID)
			}

			s.Close()
		}
	}
}

// Add this function to periodically announce our presence to the network
func announceToNetwork(ctx context.Context, h host.Host, dht *dht.IpfsDHT) {
	// Announce every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Do an initial announcement
	doAnnounce(ctx, h, dht)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			doAnnounce(ctx, h, dht)
		}
	}
}

// Helper function to perform the actual announcement
func doAnnounce(ctx context.Context, h host.Host, dht *dht.IpfsDHT) {
	// Create a context with timeout for the announcement
	announceCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	slog.Info("Announcing ourselves to the network")

	// Refresh the routing table to make ourselves more discoverable
	if err := dht.Bootstrap(announceCtx); err != nil {
		slog.Error("Failed to bootstrap DHT for announcement", "error", err)
		return
	}

	// Log our addresses for reference
	for _, addr := range h.Addrs() {
		fullAddr := addr.String() + "/p2p/" + h.ID().String()
		slog.Info("Available at address", "addr", fullAddr)
	}

	slog.Info("Successfully announced to the network")
}
