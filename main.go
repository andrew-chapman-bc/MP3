package main

import 
(
	"./unicast"
	"./utils"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"sync"
	"time"
	"strings"
	"strconv"
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

// Have to wait for server to connect before connecting to clients
func waitForServerToLoad(isLoadedChan chan bool) bool {
	isLoaded := false
	for !isLoaded {
		isLoaded = <- isLoadedChan
	}
	return true
}


func main() {
	var wg sync.WaitGroup
	// Create message channel to move messages around
	// Create serverChan to make sure we don't dial before listen
	messagesChannel := make(chan utils.Message)
	serverFinishedChan := make(chan bool)

	// Get port from command line and connections (json)
	s := getCmdLine()
	connections, err := utils.GetConnections()
	if err != nil {
		fmt.Println("Error reading json", err)
	}
	
	port := s
	
	// Create a server (Listen on port)
	serv, err := unicast.NewTCPServer(port, connections)
	if err != nil {
		fmt.Println(err)
	}
	
	wg.Add(1)

	// RunServ handles our messages which then get sent back through messagesChan
	go func() {
		err1 := serv.RunServ(messagesChannel, serverFinishedChan)
		if err1 != nil {
			fmt.Println(err1)
		}
		defer wg.Done()
	}()
	
	// Get all of the ports from the config.JSON
	portArr := utils.GetConnectionsPorts(connections)
	portArrLen := len(portArr)
	
	// Make sure we don't dial before we listen
	waitForServerToLoad(serverFinishedChan)

	cliArr := make([]*unicast.Client, portArrLen)
	var state utils.Message
	// Create our own client
	for index := range portArr {
		if portArr[index] == port {
			cli, err := unicast.NewTCPClient(port, connections)
			if err != nil {
				fmt.Println(err)
			}
			cliArr[index] = cli
			// Dials in
			err = cli.RunCli()
			if err != nil {
				fmt.Println(err)
			}
			// Get the initial state of whatever port we are
			state, err = cli.FetchInitialState()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	faultyRounds := ""
	for i := 0; i < len(connections.Connections); i++ {
		if (connections.Connections[i].Port == port) {
			faultyRounds = connections.Connections[i].Status
		}
	}
	
	// Send messages to all of the clients in the array
	for index := range cliArr {
		time.Sleep(3 * time.Second)
		
		// Send our own state to ourselves first
		if portArr[index] == port && !strings.Contains(faultyRounds, "1") {
			err5 := cliArr[index].SendMessageToServer(state)
			if err5 != nil {
				fmt.Println(err5)
			}
			continue
		}
		// Create the rest of the clients
		cli, err := unicast.NewTCPClient(portArr[index], connections)
		if err != nil {
			fmt.Println(err)
		}
		cliArr[index] = cli
		
		err = cli.RunCli()
		if err != nil {
			fmt.Println(err)
		}
		// check if node is faulty
		if (!strings.Contains(faultyRounds, "1")) {
			err2 := cli.SendMessageToServer(state)
			if err2 != nil {
				fmt.Println(err2)
				break
			}
		}
	}

	var nodes utils.NodeNums
	var messageQueue utils.Messages
	if err != nil {
		fmt.Println(err)
	}
	wg.Add(1)
	round := 1
	go func() {
		receivedNodes := 0
		isDone := false
		consensusReached := false
		// isDone checks for the specific round, consensus checks for consensus
		for !consensusReached {
			nodes, err = utils.GetNodeNums(round)
			// if n-f
			for len(messageQueue.Messages) < (nodes.TotalNodes - nodes.FaultyNodes) {
				message := <- messagesChannel
				if (message.Round == round) {
					messageQueue.Messages = append(messageQueue.Messages, message)
				}
				// fmt.Println("This is the messages queue", messageQueue)
			}
			fmt.Println(messageQueue)
			for !isDone {
				for _, val := range messageQueue.Messages {
					if (val.Round == round) {
						// increment when message is for current round
						receivedNodes++;
					}
					if (receivedNodes > nodes.TotalNodes - nodes.FaultyNodes) {
						// fmt.Println("calculating avg for andy using: ", messageQueue)
						avg, err := utils.CalculateAverage(messageQueue, round)
						if err != nil {
							fmt.Println(err)
						}

						round++
						roundString := strconv.Itoa(round)

						// Check to make sure node isn't faulty 
						if !strings.Contains(faultyRounds, roundString) {
							for _, client := range cliArr {
								client.SendMessageToServer(avg)
							}
						}

						isDone = true

						err = utils.SetJSONRound(round)
						if err != nil {
							fmt.Println(err)
						}
						
						break
					}
				}
			}

			isDone = false
			receivedNodes = 0

			// Check to see if we reached consensus
			consensusReached, err = utils.CheckForConsensus(messageQueue)
			if err != nil {
				fmt.Println(err)
			}

			// empty message queue
			messageQueue.Messages = nil
		}
		// Set the Consensus in JSON to true so our controller can read the value to run stop script
		utils.SetJSONConsensus(true)
	}()
	wg.Wait()

}


