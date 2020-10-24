package utils


import
(
	"os"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"math"
	"fmt"
	"math/rand"
	"time"
	"strings"
)

/*
	Connections: []Connection
	IP: IP Address to connect to
*/
type Connections struct {
	Connections []Connection `json:"connections"`
	IP string `json:"IP"`
	Delays Delay `json:"Delay"`
	Consensus bool `json:"Consensus"`
	Round int `json:"Round"`
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
	minDelay: lower bound for delay
	maxDelay: upper bound for delay
*/
type Delay struct {
	MinDelay int `json:"minDelay"`
	MaxDelay int `json:"maxDelay"`
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
func GetNodeNums(round int) (NodeNums, error) {
	var connections Connections
	var nodes NodeNums
	totalNodes := 0
	faultyNodes := 0
	roundString := strconv.Itoa(round)
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return nodes, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	jsonFile.Close()
	for _, val := range connections.Connections {
		if (!strings.Contains(val.Status, roundString)) {
			totalNodes++
		} else {
			totalNodes++
			faultyNodes++
		}
	}
	nodes = createNodesObj(totalNodes, faultyNodes)
	// fmt.Println("Nodes:", nodes)
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

/*
	@function: CheckForConsensus
	@description: checks for consensus and returns true when we reach approximate consensus
	@exported: True
	@params: Messages
	@returns: bool, error
*/
func CheckForConsensus(messageQueue Messages) (bool, error) {
	messageArr := messageQueue.Messages
	firstVal, err := strconv.ParseFloat(messageArr[0].State, 64)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	for i := 1; i < len(messageArr); i++ {
		curVal, err := strconv.ParseFloat(messageArr[i].State, 64)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		if (math.Abs(curVal - firstVal) > .001) {
			return false, nil
		}
	}
	return true, nil
}


/*
	@function: createDelayStruct
	@description: creates a delay struct
	@exported: False
	@params: int, int
	@returns: Delay
*/
func createDelayStruct(minDelay, maxDelay int) Delay {
	var delay Delay
	delay.MinDelay = minDelay
	delay.MaxDelay = maxDelay
	return delay
}

/*
	@function: GetDelayParams
	@description: gets delay parameters from the json and creates a delay struct to return
	@exported: True
	@params: N/A
	@returns: Delay, error
*/
func GetDelayParams() (Delay, error) {
	var connections Connections
	var delayStruct Delay
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return delayStruct, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	jsonFile.Close()
	minDelay := connections.Delays.MinDelay
	maxDelay := connections.Delays.MaxDelay
	delayStruct = createDelayStruct(minDelay, maxDelay)
	return delayStruct, nil
}

/*
	@function: GenerateDelay
	@description: takes in a delay struct and generates a delay using time.sleep
	@exported: True
	@params: Delay
	@returns: N/A
*/
func GenerateDelay(delayStruct Delay) {
	rand.Seed(time.Now().UnixNano())
	delayTime := rand.Intn(delayStruct.MaxDelay - delayStruct.MinDelay + 1) + delayStruct.MinDelay
	time.Sleep(time.Duration(delayTime) * time.Millisecond)
}

/*
	@function: SetJSONRound
	@description: Sets the Round field in JSON to the parameter
	@exported: True
	@params: int
	@returns: error
*/
func SetJSONRound(round int) error {
	var connections Connections
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	connections.Round = round
	jsonData, err := json.Marshal(connections)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile("config.json", jsonData, 0777)
	if err != nil {
		fmt.Println(err)
	}
	jsonFile.Close()
	return nil
}


/*
	@function: GetJSONRound
	@description: Gets the round from JSON
	@exported: True
	@params: N/A
	@returns: int, error
*/
func GetJSONRound() (int, error) {
	var connections Connections
	var round int
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return round, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	round = connections.Round 
	return round, nil
}

/*
	@function: GetJSONConsensus
	@description: Gets the consensus from JSON
	@exported: True
	@params: N/A
	@returns: bool, error
*/
func GetJSONConsensus() (bool, error) {
	var connections Connections
	var consensus bool
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return consensus,err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	consensus = connections.Consensus
	jsonFile.Close()
	return consensus, nil
}


/*
	@function: SetJSONConsensus
	@description: Sets the consensus field in JSON to parameter
	@exported: True
	@params: consensus
	@returns: error
*/
func SetJSONConsensus(consensus bool) error {
	var connections Connections
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	connections.Consensus = consensus
	jsonData, err := json.Marshal(connections)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile("config.json", jsonData, 0777)
	if err != nil {
		fmt.Println(err)
	}
	jsonFile.Close()
	return nil
}