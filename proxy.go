package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func proxy(poisonedConn net.Conn) {
	serveConn, destError := net.Dial("tcp", serve.String())

	if destError != nil {
		fmt.Println("problem reaching the substitute")
	}

	go func() {
		io.Copy(poisonedConn, serveConn)
	}()

	io.Copy(serveConn, poisonedConn)
}

func main() {
	FlagsInit()
	flag.Parse()

	PoisonLocalDns(hide.Domain)

	listen, listenErr := net.Listen("tcp", hide.String())

	if listenErr != nil {
		log.Fatalln("proxy failed to start at origin")
	}

	for {
		conn, _ := listen.Accept()

		go proxy(conn)
	}
}
