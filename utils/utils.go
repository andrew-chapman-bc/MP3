package utils


import (
	"os"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
)
/*
	Connections: []Connection
	IP: IP Address to connect to
*/
type Connections struct {
	Connections []Connection `json:"connections"`
	IP string `json:"IP"`
} 

/*
	State: State of the Node
	Port: "1234", etc. Port attached to Node
	Status: Whether or not the node is faulty
*/
type Connection struct {
	State string `json:"State"`
	Port string `json:"Port"`
	Status string `json:"Status"`
}

/*
	State: The state of the message
	Round: What round the message is sent in
*/
type Message struct {
	State string
	Round int
}
/*
	Messages: An array of message
*/
type Messages struct {
	Messages []Message
}
/*
	TotalNodes: Total nodes in our distributed system
	FaultyNodes: Total amt of faulty nodes in our sys
	(We get these from the JSON config file)
*/
type NodeNums struct {
	TotalNodes int
	FaultyNodes int
}

/*
	@function: GetConnections
	@description: Reads the config file
	@exported: True
	@params: N/A
	@returns: Connections, error
*/
func GetConnections() (Connections, error) {
	jsonFile, err := os.Open("config.json")
	var connections Connections
	if err != nil {
		return connections, errors.New("Error opening JSON file on Client Side")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	return connections, nil
}


/*
	@function: CreateMessage
	@description: Creates a message struct 
	@exported: True
	@params: string, int
	@returns: Message
*/
func CreateMessage(state string, round int) Message {
	var message Message
	message.State = state
	message.Round = round
	return message
}

/*
	@function: GetNodeNums
	@description: Gets the totalNodes and faultyNodes from JSON
	@exported: True
	@params: N/A
	@returns: NodeNums, error
*/
func GetNodeNums() (NodeNums, error) {
	var connections Connections
	var nodes NodeNums
	totalNodes := 0
	faultyNodes := 0
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return nodes, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	for _, val := range connections.Connections {
		if (val.Status == "") {
			totalNodes++
		} else {
			totalNodes++
			faultyNodes++
		}
	}
	nodes = createNodesObj(totalNodes, faultyNodes)
	return nodes, nil
}

/*
	@function: createNodesObj
	@description: creates a nodes obj (used in GetNodeNums)
	@exported: false
	@params: int, int
	@returns: NodeNums
*/
func createNodesObj(total, faulty int) NodeNums {
	var nodes NodeNums
	nodes.TotalNodes = total
	nodes.FaultyNodes = faulty
	return nodes
}

/*
	@function: CalculateAverage
	@description: Calculates the new state by averaging them and returns a Message struct with the new state + round
	@exported: True
	@params: Messages, int
	@returns: Message, error
*/
func CalculateAverage(messages Messages, index int) (Message, error) {
	total := 0.00
	divisor := 0.00
	var newMess Message
	round := index
	for i := 0; i < len(messages.Messages); i++ {
		stateFloat, err := strconv.ParseFloat(messages.Messages[i].State, 64)
		if err != nil {
			return newMess, err
		}
		total += stateFloat
		divisor++
	}
	newState := total/divisor
	newStateString := strconv.FormatFloat(newState, 'f', 4, 64)
	newRound := round + 1
	newMess = CreateMessage(newStateString, newRound)
	return newMess, nil
}

/*
	@function: GetConnectionsPorts
	@description: Gets all of the ports from the connections which comes from JSON
	@exported: True
	@params: Connections
	@returns: []String
*/
func GetConnectionsPorts(connections Connections) []string {
	var portArr []string
	for _, connection := range connections.Connections {
		portArr = append(portArr, connection.Port)
	}
	return portArr
}
