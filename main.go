cat > main.go << 'EOF'
package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	listenAddr = getEnv("LAGLESS_LISTEN", "0.0.0.0:7001")
	gameAddr   = getEnv("LAGLESS_GAME",   "127.0.0.1:7002")
	headerSize = 12
)

func main() {
	fmt.Println("Lagless Relay Server v0.1")
	fmt.Printf("Listening on  : %s\n", listenAddr)
	fmt.Printf("Game server   : %s\n", gameAddr)

	gameUDPAddr, err := net.ResolveUDPAddr("udp", gameAddr)
	if err != nil { log.Fatalf("Bad game addr: %v", err) }

	listenUDPAddr, _ := net.ResolveUDPAddr("udp", listenAddr)
	conn, err := net.ListenUDP("udp", listenUDPAddr)
	if err != nil { log.Fatalf("Failed to listen: %v", err) }
	defer conn.Close()

	gameConn, err := net.DialUDP("udp", nil, gameUDPAddr)
	if err != nil { log.Fatalf("Failed to connect to game server: %v", err) }
	defer gameConn.Close()

	buf := make([]byte, 4096)
	var forwarded, acked uint64
	lastLog := time.Now()
	fmt.Println("Ready — waiting for packets...")

	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil { continue }
		if n < headerSize { continue }

		seq := binary.BigEndian.Uint64(buf[0:8])
		payload := buf[headerSize:n]
		gameConn.Write(payload)
		forwarded++

		clientAckAddr := &net.UDPAddr{IP: clientAddr.IP, Port: 7003}
		ackPayload := make([]byte, headerSize)
		binary.BigEndian.PutUint64(ackPayload[0:8], seq)
		conn.WriteToUDP(ackPayload, clientAckAddr)
		acked++

		if time.Since(lastLog) > 5*time.Second {
			fmt.Printf("[relay] forwarded: %d | acked: %d\n", forwarded, acked)
			lastLog = time.Now()
		}
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok { return val }
	return fallback
}
EOF