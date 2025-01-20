package main

import (
	"log"
	"net"
	"time"

	"github.com/rzolotuhin/go-zabbix-agent-protocol"
)

const agentAddress = "127.0.0.1:10050"
const requestedKey = "net.if.discovery"

func main() {
	conn, err := net.DialTimeout("tcp", agentAddress, time.Second*3)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	agent := zabbix.Agent{
		Transport: conn,
	}

	answer, err := agent.Get(requestedKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s [%s]: %s\n", agentAddress, requestedKey, answer)
}
