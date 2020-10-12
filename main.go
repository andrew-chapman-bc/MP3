package main

import (
	"./unicast"
	"bufio"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"strings"
	"sync"
	"errors"
)





func getCmdLine() string {
	// Use argparse library to get accurate command line data
	parser := argparse.NewParser("", "Concurrent TCP Channels")
	s := parser.String("s", "string", &argparse.Options{Required: true, Help: "String to print"})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}
	return *s
}

func main() {

	messageChannel := make(chan message)
	s := getCmdLine()
	
	
	connections, err := utils.GetConnections()
	if err != nil {
		fmt.Println("Error reading json", err)
	}
	
	
	var serv *unicast.Server
	port := cmdLineArr[1]
	serv, err := unicast.NewTCPServer(port, connections)
	if err != nil {
		fmt.Println(err)
	}
	
	
	go func() {
		err := unicast.serv.RunServ(messageChannel)
	}

	portArr := utils.GetConnectionsPorts(connections)
	var cliArr [3]*unicast.Client
	for index, _ := range cliArr {
		cli, err := unicast.NewTCPClient(portArr[index], connections)
		cliArr[index] = cli
	}

	go func() {
		newMessage := <- messageChannel
		for index, client := range cliArr {
			unicast.client.sendMessageToServer(newMessage)
		}
	}
	// 


}

