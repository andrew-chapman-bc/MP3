package main

import (
	"./unicast"
	"./utils"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"sync"
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

	var wg sync.WaitGroup
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
	

	wg.Add(1)
	go func(){
		defer wg.Done()
		err1 := serv.RunServ(messageChannel)
		if err1 != nil {
			fmt.Println(err)
		}
	}()

	portArr := utils.GetConnectionsPorts(connections)
	portArrLen := len(portArr)

	cliArr := make([]*unicast.Client, portArrLen)
	for index := range cliArr {
		fmt.Println(portArr[index])
		cli, err := unicast.NewTCPClient(portArr[index], connections)
		if err != nil {
			fmt.Println(err)
		}
		cliArr[index] = cli
	}

	wg.Add(1)
	go func() {
		for {
			newMessage := <- messageChannel
			fmt.Println(newMessage)
			for _, client := range cliArr {
				err := client.SendMessageToServer(newMessage)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
	// 
	wg.Wait()

}

