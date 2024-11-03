package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	mode := flag.String("mode", "tcp-server", "Mode: tcp-server, tcp-client, udp-server, udp-client")
	address := flag.String("address", "localhost:8080", "Server address (e.g., localhost:8080)")
	flag.Parse()

	switch *mode {
	case "tcp-server":
		runTCPServer(*address)
	case "tcp-client":
		runTCPClient(*address)
	case "udp-server":
		runUDPServer(*address)
	case "udp-client":
		runUDPClient(*address)
	default:
		fmt.Println("Unknown mode:", *mode)
		os.Exit(1)
	}
}

func runTCPServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("TCP server started at", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 20480)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}
		message := string(buf[:n])
		fmt.Println("Received message:", message)

		_, err = conn.Write([]byte("Received: " + message))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}

func runTCPClient(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to TCP server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to TCP server", address)
	sendMessages(conn)
}

func runUDPServer(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP server address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server started at", address)

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			continue
		}
		message := string(buf[:n])
		fmt.Println("Received from", clientAddr, ":", message)

		_, err = conn.WriteToUDP([]byte("Received: "+message), clientAddr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
}

func runUDPClient(address string) {
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to UDP server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to UDP server", address)
	sendMessages(conn)
}

func sendMessages(conn net.Conn) {
	var input string
	for {
		fmt.Print("Enter message (or exit to quit): ")
		fmt.Scanln(&input)
		if strings.ToLower(input) == "exit" {
			fmt.Println("Closing connection...")
			return
		}
		if len(input) > 20480 {
			fmt.Println("Message exceeds 20 KB, it will be split.")
		}
		_, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		buf := make([]byte, 20480)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error receiving response:", err)
			return
		}
		fmt.Println("Server response:", string(buf[:n]))
	}
}
