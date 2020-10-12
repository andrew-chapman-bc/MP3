package unicast

import (
	"encoding/gob"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
	@function: ScanConfigForServer
	@description: Scans the config file for all of the ports that will be used to open concurrent TCP Servers
	@exported: True
	@params: N/A
	@returns: []string
*/
func scanConfigForFaultyNodes() ([]int, error) {
	config, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	// [n, f]
	nodeArray := []int{}
	connections := new(Connections)
	byteValue, err := ioutil.ReadAll(config)
	json.Unmarshal(byteValue, &connections)
	faultyCounter := 0
	nodeCounter := len(connections.Connections)
	for i := 0; i < len(connections.Connections); i++ {
		if connections.Connections[i].Status == "faulty" {
			faultyCounter++
		}
	}
	nodeArray = append(nodeArray, nodeCounter, faultyCounter)
	return nodeArray, errors.New("Cannot find port")
}


/*
	@function: CreateUserInputStruct
	@description: Uses a destination, message and source string to construct a UserInput struct to send and receive the message across the server/client
	@exported: True
	@params: string, string, string
	@returns: {UserInput}
*/
func CreateUserInputStruct(state float64, source string) UserInput {
	var input UserInput
	input.Round = 0
	input.State = state
	input.Source = source
	return input
}


/*
	@function: handleConnection
	@description: handles connections to the concurrent TCP client and receives messages that are sent over through a goroutine in connectToTCPClient
	@exported: False
	@params: net.Conn
	@returns: N/A
*/
func handleConnection(c net.Conn, valueChan chan UserInput) error {
	fmt.Println("up here")
	nodeArray, err := scanConfigForFaultyNodes()
	if (err != nil) {
		fmt.Println(err)
		return errors.New("There was an error scanning node array")
	}
	// even though we don't support multi-messaging at the moment, no reason to possibly be running this multiple times inside the for loop
	delay, err := getDelayParams()
	if (err != nil) {
		fmt.Println("Error: ", err)
		return errors.New("Error getting delay params")
	}
	inputArray := []UserInput{}
	// [userInput: {State1, round1, 1234}, userInput: {State2, round1, 4567}]
	fmt.Println("here")
	decoder := gob.NewDecoder(c)
	fmt.Println("didn't make it here")
	roundCounter := 1
	for {
		// state round source
        if err != nil {
            fmt.Println(err)
            return err
		}
		// generate the network delay on the receive side, must do it here and not in the sendmessage function because we are using goroutines
		generateDelay(delay)

		// [ [State, round, source], [State, round, source] ]
		var input UserInput
		_ = decoder.Decode(&input)
		counter := 0
		fmt.Println(inputArray)
		inputArray = append(inputArray, input)
		// function to check if within .0001
		
		nodesWaitedFor := nodeArray[0] - nodeArray[1]
		for _, val := range inputArray {
			if val.Round == roundCounter {
				counter += 1
			}
			if counter == nodesWaitedFor {
				newValue, err := getAvgValues(inputArray, roundCounter)
				if err != nil {
					fmt.Println(err)
					return err
				}

				
				var send UserInput
				send.State = newValue
				send.Round = roundCounter
				sendMessageFromServer(send)

				break
			}
		}
		roundCounter++
		
	}
}

func sendMessageFromServer(input UserInput) {
	config, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	var connections Connections
	byteValue, err := ioutil.ReadAll(config)
	json.Unmarshal(byteValue, &connections)
	SendMessage(input, connections)
}

func getAvgValues(inputArr []UserInput, roundCounter int) (float64, error) {
	counter := 0.00
	total := 0.00
	for _, val := range inputArr {
		if val.Round == roundCounter {
			total += val.State
			counter++
		}

	}
	averageVal := total/counter
	return averageVal, nil
}



/*
	@function: getDelayParams
	@description: Scans the config file for the first line to get the delay parameters that will be used to simulate the network delay
	@exported: false
	@params: N/A 
	@returns: Delay, error
*/
func getDelayParams() (Delay, error) {
	config, err := os.Open("config.txt")
	scanner := bufio.NewScanner(config)
	scanner.Split(bufio.ScanLines)
	success := scanner.Scan()
	if success == false {
		err = scanner.Err()
		if err == nil {
			fmt.Println("Scanned first line")
		} else {
			log.Fatal(err)
		}
	}
	delays := strings.Fields(scanner.Text())
	var delayStruct Delay
	delayStruct.minDelay = delays[0]
	delayStruct.maxDelay = delays[1]
	return delayStruct, err
} 

/*
	@function: generateDelay
	@description: Uses the delay parameters obtained from getDelayParams() to generate the delay that will be used in sendMessage function
	@exported: false
	@params: Delay
	@returns: N/A
*/
func generateDelay (delay Delay) {
	rand.Seed(time.Now().UnixNano())
	min, _ := strconv.Atoi(delay.minDelay)
	max, _ := strconv.Atoi(delay.maxDelay)
	delayTime := rand.Intn(max - min + 1) + min
	time.Sleep(time.Duration(delayTime) * time.Millisecond)
} 

/*
	@function: connectToTCPClient
	@description: Opens a concurrent TCP Server and calls the net.Listen function to connect to the TCP Client
	@exported: True
	@params: string
	@returns: N/A
*/
func ConnectToTCPClient(PORT string, valueChan chan UserInput) {
	// listen/connect to the tcp client
	fmt.Println("shit")
	l, err := net.Listen("tcp4", ":" + PORT)
	fmt.Println(PORT)
	fmt.Println("shit2")
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	for {
		fmt.Println("hereeeee")
		c, err := l.Accept()
		fmt.Println("penis")
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(c, valueChan)
	}
}