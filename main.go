package main

import (
	"./unicast"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"io/ioutil"
	"os"
	"strconv"
)

// go run main.go 1234
//.5

/*
	@function: getInput
	@description: gets the input entered through I/O and packages it into an array that will be used to create a {UserInput}
	@exported: False
	@params: N/A
	@returns: []string
*/
/*
func getInput() float64 {
	fmt.Println("Enter input >> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	s, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return s
	}
	return 0
}
*/
/*
	@function: parseInput
	@description: Parses the UserInput into a {UserInput} and calls ScanConfig() to parse the parameters of TCP connection into a {Connection}
	@exported: False
	@params: N/A
	@returns: {UserInput}, {Connection}
*/

/*
func parseInput(source *string) (unicast.UserInput, unicast.Connection) {
	input := getInput()
	inputStruct := unicast.CreateUserInputStruct(input, *source)
	connection := unicast.ScanConfigForClient(inputStruct)
	return inputStruct, connection
}
*/

/*
	@function: openTCPServerConnections
	@description: Opens all of the ports defined in the config file using ScanConfigForServer() to get an array of ports 
					and ConnectToTCPClient() to open them
	@exported: False
	@params: {WaitGroup}
	@returns: N/A
*/
func openTCPServerConnections(source *string, valueChan chan unicast.UserInput) error {
	// Need to send the source string in here so we know what port to look for
	// openPort, err := unicast.ScanConfigForServer(*source)
	fmt.Println("we are here in the openTCP", *source)
	if *source == "" {
		return errors.New("Source string is incorrect")
	}
	unicast.ConnectToTCPClient(*source, valueChan)
	return nil
}

func parseJSON(source *string) (unicast.UserInput, unicast.Connections) {
	config, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	var connections unicast.Connections
	byteValue, err := ioutil.ReadAll(config)
	json.Unmarshal(byteValue, &connections)
	var initialNode unicast.UserInput
	for i := 0; i < len(connections.Connections); i++ {
		if connections.Connections[i].Port == *source {
			s, err := strconv.ParseFloat(connections.Connections[i].State, 64)
			if err != nil {
				fmt.Println("Error in initial Node")
			}
			initialNode.State = s
		}
	}
	initialNode.Source = *source
	return initialNode, connections
}

/*
	@function: unicast_send
	@description: function used as a goroutine to call SendMessage() to pass data from client to server, utilizes waitgroup
	@exported: False
	@params: {UserInput}, {Connection}, {WaitGroup}
	@returns: N/A
*/
func unicastSend(inputStruct unicast.UserInput, connection unicast.Connections) {
	//defer wg.Done()
	// Send the message using UserInput struct and Connection struct to easily pass around data
	unicast.SendMessage(inputStruct, connection)
}

func main() {
	// Use argparse library to get accurate command line data
	parser := argparse.NewParser("", "Concurrent TCP Channels")
	i := parser.Int("i", "int", &argparse.Options{Required: true, Help: "Source destination/identifiers"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	s := strconv.Itoa(*i)

	// s is the source

	valueChannel := make(chan unicast.UserInput)

	// pass in the source to listen to it
	go openTCPServerConnections(&s, valueChannel)

	inputStruct, connection := parseJSON(&s)

	unicastSend(inputStruct, connection)


}