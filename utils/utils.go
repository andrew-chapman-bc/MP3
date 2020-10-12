package utils


import (
	"os"
	"encoding/json"
	"errors"
	"io/ioutil"
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
	Type: "Server"/"Client" whether it's server or client
	Port: "1234", etc. Port attached to username
	Username: name of connection
	IP: IP address to connect to
*/
type Connection struct {
	State string `json:"State"`
	Port string `json:"Port"`
	Status string `json:"Status"`
}

type Message struct {
	State string
	Round int
}

type Messages struct {
	Messages []Message
	MessagesArr []Messages
}

type NodeNums struct {
	totalNodes int
	faultyNodes int
}

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

func CreateMessage(state string, round int) Message {
	var message Message
	message.State = state
	message.Round = round
	return message
}

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

func createNodesObj(total, faulty int) NodeNums {
	var nodes NodeNums
	nodes.totalNodes = total
	nodes.faultyNodes = faulty
}

func calculateAverage(messages Messages, index int ) (Message, error) {
	total := 0
	
}