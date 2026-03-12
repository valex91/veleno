package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

func proxy(poisonedConn net.Conn, logger io.Writer) {
	serveConn, destError := net.Dial("tcp", serve.String())

	if destError != nil {
		fmt.Println("problem reaching the substitute")
	}

	go func() {
		io.Copy(poisonedConn, serveConn)
	}()

	var poisonedConnReader io.Reader

	if logger != nil {
		poisonedConnReader = io.TeeReader(poisonedConn, logger)
	} else {
		poisonedConnReader = poisonedConn
	}

	io.Copy(serveConn, poisonedConnReader)
}

func main() {
	FlagsInit()
	flag.Parse()
	PoisonLocalDns(hide.Domain)

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)
	connectionChannel := make(chan net.Conn)

	var listen net.Listener
	var listenErr error
	var logger *os.File

	if isLogging {
		logger = LogToFile()
		defer logger.Close()
	}

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

	go func() {
		for {
			conn, _ := listen.Accept()
			connectionChannel <- conn
		}
	}()

	for {
		select {
		case <-interruptChannel:
			UnPoisonLocalDns(hide.Domain)
			os.Exit(0)

		case conn := <-connectionChannel:
			go proxy(conn, logger)
		}
	}
}
