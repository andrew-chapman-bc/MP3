package main

import (
	"./unicast"
	"./utils"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
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

	messageChannel := make(chan utils.Message)
	s := getCmdLine()
	
	
	connections, err := utils.GetConnections()
	if err != nil {
		fmt.Println("Error reading json", err)
	}
	
	port := s
	serv, err := unicast.NewTCPServer(port, connections)
	if err != nil {
		fmt.Println(err)
	}
	
	
	err1 := serv.RunServ(messageChannel)
	if err1 != nil {
		fmt.Println(err)
	}

	portArr := utils.GetConnectionsPorts(connections)
	var cliArr [3]*unicast.Client
	for index := range cliArr {
		cli, err := unicast.NewTCPClient(portArr[index], connections)
		if err != nil {
			fmt.Println(err)
		}
		cliArr[index] = cli
	}

	go func() {
		newMessage := <- messageChannel
		for _, client := range cliArr {
			err := client.SendMessageToServer(newMessage)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	// 


}

