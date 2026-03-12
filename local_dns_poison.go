package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const POISON_C = "## POISON ##"
const POISON_C_END = "## END POISON ##"
const LOCAL_DNS_FILE = "/etc/hosts"

type PoisonInfo struct {
	BlockEnd        int
	BlockStart      int
	ExistingContent []string
}

func shouldPoison(domain string) bool {
	return net.ParseIP(domain) == nil
}

func hasPoisonBlock(text string) bool {
	return strings.Contains(text, POISON_C)
}

func hasPoisonEnd(text string) bool {
	return strings.Contains(text, POISON_C_END)
}

func cleanLocalDns(localDnsInfo PoisonInfo) {
	hostFile, err := os.OpenFile(LOCAL_DNS_FILE, os.O_RDWR|os.O_TRUNC, 0644)

	if err != nil {
		log.Panicln("failed to clean up poison")
	}

	defer hostFile.Close()

	upToBlock := localDnsInfo.ExistingContent[:localDnsInfo.BlockStart]
	afterBlock := localDnsInfo.ExistingContent[localDnsInfo.BlockEnd+1:]

	noPoisonFile := strings.Join(append(upToBlock, afterBlock...), "\n")
	hostFile.Write([]byte(noPoisonFile))
}

func writeLocalDns(localDnsInfo PoisonInfo, domain string) {
	hostFile, err := os.OpenFile(LOCAL_DNS_FILE, os.O_RDWR|os.O_TRUNC, 0644)

	if err != nil {
		log.Panicln("failed local poison")
	}

	defer hostFile.Close()

	if localDnsInfo.BlockEnd == 0 {
		localDnsInfo.ExistingContent = append(localDnsInfo.ExistingContent, POISON_C)
		localDnsInfo.ExistingContent = append(localDnsInfo.ExistingContent, fmt.Sprintf("127.0.0.1 %s", domain))
		localDnsInfo.ExistingContent = append(localDnsInfo.ExistingContent, POISON_C_END)
	} else {
		localDnsInfo.ExistingContent = appendAtIndex(localDnsInfo.ExistingContent, localDnsInfo.BlockEnd, fmt.Sprintf("127.0.0.1 %s", domain))
	}

	jointLines := strings.Join(localDnsInfo.ExistingContent, "\n")
	hostFile.Write([]byte(jointLines))
}

func createLocalDnsInfo(domain string) PoisonInfo {
	hostFile, err := os.OpenFile(LOCAL_DNS_FILE, os.O_RDONLY, 0644)

	if err != nil {
		log.Panicln("FAILED TO READ LOCAL DNS")
	}

	defer hostFile.Close()
	existingContent := make([]string, 0)

	counter := 0
	blockEnd := 0
	blockStart := 0
	scanner := bufio.NewScanner(hostFile)

	for scanner.Scan() {
		line := scanner.Text()
		existingContent = append(existingContent, line)

		if hasPoisonBlock(line) {
			blockStart = counter
		}

		if hasPoisonEnd(line) {
			blockEnd = counter
		}

		counter++
	}

	return PoisonInfo{
		BlockEnd:        blockEnd,
		ExistingContent: existingContent,
		BlockStart:      blockStart,
	}
}

func PoisonLocalDns(domain string) {
	dnsInfo := createLocalDnsInfo(domain)
	writeLocalDns(dnsInfo, domain)
}

func UnPoisonLocalDns(domain string) {
	localDnsInfo := createLocalDnsInfo(domain)
	cleanLocalDns(localDnsInfo)
}
