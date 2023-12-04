package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	socketPath := "/tmp/ootp.sock"
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		fmt.Println("No command specified")
		return
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}

	message := os.Args[1]
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}

	conn.Close()
}
