package main

import (
	"crypto/tls"
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

	var listen net.Listener
	var listenErr error

	if isTls {
		cert, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")

		if err != nil {
			log.Fatalln("failed to load cert")
		}

		configs := tls.Config{Certificates: []tls.Certificate{cert}}
		tlsListen, tlsErr := tls.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", hide.Port), &configs)
		listen = tlsListen
		listenErr = tlsErr
	} else {
		nonTlslisten, NonTlslistenErr := net.Listen("tcp", hide.String())
		listen = nonTlslisten
		listenErr = NonTlslistenErr
	}

	if listenErr != nil {
		fmt.Println(listenErr)
		log.Fatalln("proxy failed to start at origin")
	}

	for {
		conn, _ := listen.Accept()

		go proxy(conn)
	}
}
