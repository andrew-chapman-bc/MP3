package main

import (
	"./unicast"
	"./utils"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"sync"
	"time"
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
	go func() {
		defer wg.Done()
		err1 := serv.RunServ(messageChannel)
		if err1 != nil {
			fmt.Println(err1)
		}
	}()
	
	portArr := utils.GetConnectionsPorts(connections)
	portArrLen := len(portArr)

	cliArr := make([]*unicast.Client, portArrLen)
	for index := range cliArr {
		// 1234 will always be the first cli here
		// Need to make it so we send our port's data over first instead of always 1234
		cli, err := unicast.NewTCPClient(portArr[index], connections)
		if err != nil {
			fmt.Println(err)
		}
		cliArr[index] = cli
		
		err = cli.RunCli()
		if err != nil {
			fmt.Println(err)
		}
		newMessage, err := cli.FetchInitialState()
		err2 := cli.SendMessageToServer(newMessage)
		if err2 != nil {
			fmt.Println(err2)
			break
		}
		time.Sleep(15 * time.Second)
	}

	wg.Add(1)
	go func() {
		fmt.Println("hi")
		// if n-f
		// 		calculateAvg
		//		sendMessage
		// for {
		// 	newMessage := <- messageChannel
		// 	fmt.Println("this is the message we get from the channel", newMessage)
		// 	for _, client := range cliArr {
		// 		err := client.SendMessageToServer(newMessage)
		// 		if err != nil {
		// 			fmt.Println(err)
		// 		}
		// 	}
		// }
	}()
	// 
	wg.Wait()

}

