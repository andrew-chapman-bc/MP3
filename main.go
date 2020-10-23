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
	messagesChannel := make(chan utils.Message)
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
		err1 := serv.RunServ(messagesChannel)
		if err1 != nil {
			fmt.Println(err1)
		}
		defer wg.Done()
	}()
	
	portArr := utils.GetConnectionsPorts(connections)
	portArrLen := len(portArr)

	cliArr := make([]*unicast.Client, portArrLen)
	var state utils.Message
	for index := range portArr {
		if portArr[index] == port {
			cli, err := unicast.NewTCPClient(port, connections)
			if err != nil {
				fmt.Println(err)
			}
			cliArr[index] = cli
			err = cli.RunCli()
			if err != nil {
				fmt.Println(err)
			}
			state, err = cli.FetchInitialState()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	for index := range cliArr {
		// 1234 will always be the first cli here
		// Need to make it so we send our port's data over first instead of always 1234
		time.Sleep(20 * time.Second)
		if portArr[index] == port {
			err5 := cliArr[index].SendMessageToServer(state)
			if err5 != nil {
				fmt.Println(err5)
			}
			continue
		}
		cli, err := unicast.NewTCPClient(portArr[index], connections)
		if err != nil {
			fmt.Println(err)
		}
		cliArr[index] = cli
		
		err = cli.RunCli()
		if err != nil {
			fmt.Println(err)
		}

		err2 := cli.SendMessageToServer(state)
		if err2 != nil {
			fmt.Println(err2)
			break
		}
	}

	var nodes utils.NodeNums
	var messageQueue utils.Messages
	nodes, err = utils.GetNodeNums()
	if err != nil {
		fmt.Println(err)
	}
	wg.Add(1)
	go func() {
		fmt.Println("hi")
		for len(messageQueue.Messages) < (nodes.TotalNodes - nodes.FaultyNodes) {
			message := <- messagesChannel
			messageQueue.Messages = append(messageQueue.Messages, message)
			fmt.Println("This is the messages queue", messageQueue)
		}
		round := 1
		receivedNodes := 0
		isDone := false
		for !isDone {
			for _, val := range messageQueue.Messages {
				if (val.Round == round) {
					receivedNodes++;
				}
				if (receivedNodes > nodes.TotalNodes - nodes.FaultyNodes) {
					fmt.Println("calculating avg for andy using: ", messageQueue)
					avg, err := utils.CalculateAverage(messageQueue, round)
					if err != nil {
						fmt.Println(err)
					}
					for _, client := range cliArr {
						fmt.Println("sending here")
						client.SendMessageToServer(avg)
					}
					isDone = !isDone
					round++
					break
				}
			}
		}
	}()
	// 
	wg.Wait()

}

