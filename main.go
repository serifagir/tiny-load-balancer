package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	listenAddress = "localhost:8080"

	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}

	tracker = 0
)

func main() {
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %s ", err)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept backend server connection: %s", err)
		}

		backendServer := selectBackendServer()
		go func() {
			err := proxy(backendServer, connection)
			if err != nil {
				log.Printf("FATAL: proxying process failed: %v", err)
			}
		}()
		proxy(backendServer, connection)
	}

}

func proxy(backendServer string, c net.Conn) error {
	backendConnection, err := net.Dial("tcp", backendServer)
	if err != nil {
		return fmt.Errorf("Failed to connect backend server %s", err)
	}

	go io.Copy(backendConnection, c)

	go io.Copy(c, backendConnection)

	return nil
}

func selectBackendServer() string {
	serverIdx := tracker % len(server)
	tracker++
	return server[serverIdx]
}
