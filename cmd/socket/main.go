package main

import (
	"fmt"
	"net"
	"os"
	"pass/socket"
	"pass/storage"
)

func main() {
	socketPath := "/tmp/pass.sock"
	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Socket listening on: ", socketPath)
	defer listener.Close()
	masterKey, err := storage.LoadMasterKey("Master.key")
	if err != nil {
		panic(err)
	}

	vault, err := storage.LoadVault(masterKey, "Vault.vault")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Request comming in")
		socket.Handle(conn, masterKey, &vault, "Vault.vault")
	}
}
