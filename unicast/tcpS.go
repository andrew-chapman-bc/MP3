package unicast

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"time"
	"os"
	"log"
	"errors"
	"math/rand"
	"strconv"
	"encoding/gob"
	"ioutil"
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
	nodeArray := []int
	connections := new(Connections)
	byteValue, err := ioutil.ReadAll(config)
	json.Unmarshal(byteValue, &connections)
	faultyCounter := 0
	nodeCounter := len(connections.Connections)
	for i := 0; i < len(connections.Connections); i++ {
		if connections.Connections[i].Status == "faulty" {
			faultyCounter += 1
		}
	}
	nodeArray.append(nodeArray, nodeCounter, faultyCounter)
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
func handleConnection(c net.Conn) error {
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
	decoder := gob.NewDecoder(c)
	for {
		// state round source
        if err != nil {
            fmt.Println(err)
            return
		}
		// generate the network delay on the receive side, must do it here and not in the sendmessage function because we are using goroutines
		generateDelay(delay)

		// [ [State, round, source], [State, round, source] ]
		input := new(UserInput)
		decoder.Decode(&input)
		roundCounter := 1
		counter := 0
		inputArray = append(inputArray, input)
		// function to check if within .0001
		
		nodesWaitedFor := nodeArray[0] - nodeArray[1]
		for index, val := range inputArray {
			if (val.Round == roundCounter) {
				counter += 1
			}
			if (counter == nodesWaitedFor) {
				newValue, err := getAvgValues(inputArray, roundCounter)
				// roundCounter++
				break
			}
		}
		
	}
}

func getAvgValues(inputArr []UserInput, roundCounter int) (float64, error) {
	for index, val := range inputArr {
		
	}
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
func ConnectToTCPClient(PORT string, valueChan chan float64) {
	// listen/connect to the tcp client
	l, err := net.Listen("tcp4", ":" + PORT)
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(c, valueChan)
	}
}