package main

import (
	"fmt"
	"net"
	"os"
)

type Communication struct {
	socketPath string
	listen     net.Listener
	cli        chan string
}

func NewCommunication(socketPath string, cli chan string) (*Communication, error) {
	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}
	return &Communication{socketPath: socketPath, listen: listener, cli: cli}, nil
}

func (c *Communication) Close() {
	defer c.listen.Close()
	os.Remove(c.socketPath)
}

func (c *Communication) Listen() {
	for {
		conn, err := c.listen.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			os.Exit(1)
		}
		go c.handleConnection(conn)
	}
}

func (c *Communication) handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	message := string(buffer[:n])

	c.cli <- message

	fmt.Println("Received message:", message)

	conn.Close()
}
