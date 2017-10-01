package main

import (
	"fmt"
	"log"
	"net"
)

type client struct {
	Addr string
	Name string
	Conn net.Conn
}

func main() {
	ln, err := net.Listen("tcp", ":8089")

	if err != nil {
		log.Fatal(err)
	}

	newConn := make(chan net.Conn)

	go func() {
		for {
			conn, err := ln.Accept()

			if err != nil {
				panic(err)
			}

			fmt.Println("new connection")
			newConn <- conn
		}
	}()

	handleConnection(newConn)

}

func handleConnection(newConn chan net.Conn) {
	conns := []net.Conn{}
	msgs := make(chan []byte)

	for {
		select {
		case conn := <-newConn:
			conns = append(conns, conn)

			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)

					if err != nil {
						fmt.Println("connection closed")
						break
					}
					data := buf[:n]
					msgs <- data
				}
			}()

		case msg := <-msgs:

			fmt.Println("message receieved", string(msg))
			for _, conn := range conns {
				conn.Write(msg)
			}
		}
	}
}
