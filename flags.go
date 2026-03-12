package main

import (
	"flag"
	"fmt"
	"strings"
)

type ConnectionInfo struct {
	Domain string
	Port   string
}

var hide ConnectionInfo
var serve ConnectionInfo
var isTls bool
var isLogging bool

func (connectionInfo *ConnectionInfo) String() string {
	return fmt.Sprintf("%s:%s", connectionInfo.Domain, connectionInfo.Port)
}

func (connectionInfo *ConnectionInfo) Set(flagValue string) error {
	split := strings.Split(flagValue, ":")

	connectionInfo.Domain = split[0]
	connectionInfo.Port = split[1]

	return nil
}

func FlagsInit() {
	flag.Var(&hide, "hide", "domain to substitute formatted as <address>:<port>")
	flag.Var(&serve, "serve", "replace to serve formatted as <address>:<port>")
	flag.BoolVar(&isTls, "tls", false, "pass to use https")
	flag.BoolVar(&isLogging, "isLogging", false, "pass to log onto a file")
}
