package unicast

import (
	"fmt"
	"errors"
	"net"
	"os"
	"io/ioutil"
	"encoding/json"
	"../utils"
	"strings"
)


// Client holds the structure of our TCP Client implementation
type Client struct {
	port string
	client net.Conn
	Connections utils.Connections
}

/*
	@function: NewTCPClient
	@description: Creates a Client instance which can be used in the main function
	@exported: True
	@family: N/A
	@params: string
	@returns: {*Client}, error
*/
func NewTCPClient(port string, connections utils.Connections) (*Client, error) {
	client := Client{port: port, Connections: connections}
	// if username is empty -> throw error
	if port == "" {
		fmt.Println("error here 1")
		return nil, errors.New("Error: Port not found")
	}

	return &client, nil
}


/*
	@function: RunCli
	@description: Starts the TCP client which calls the function to send message to server
	@exported: True
	@family: Client
	@params: chan {Message}
	@returns: error
*/
func (cli *Client) RunCli() (err error) {
	fmt.Println(cli.port)
	cli.client, err = net.Dial("tcp", cli.port)
	if err != nil {
		return err
	}

	return nil

}

/*
	@function: sendMessageToServer
	@description: Reads the message channel and serializes the data to send over to server
	@exported: false
	@family: Client
	@params: net.Conn, chan {Message}
	@returns: error
*/
func (cli *Client) sendMessageToServer(messageData utils.Message) (err error) {
	

	jsonData, err := json.Marshal(messageData)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(cli.client)
	encoder.Encode(jsonData)
	fmt.Println("data sent!", messageData)
	return
}


/*
	@function: readJSONForClient
	@description: Reads the JSON File and adds to it if needed, then returns the specific connection that is needed
	@exported: false
	@family: Client
	@params: string
	@returns: {Connection}, error
*/
func (cli *Client) readJSONForClient(userName string) (util.Connection, error) {
	jsonFile, err := os.Open("connections.json")
	var connections util.Connections
	ourConnect := util.Connection{"","",""}
	if err != nil {
		return ourConnect, errors.New("Error opening JSON file on Client Side")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	for i := 0; i < len(connections.Connections); i++ {
		if connections.Connections[i].Username == userName {
			connections.Connections[i].Port = connections.IP + ":" + connections.Connections[i].Port
			return connections.Connections[i], nil
		}
	}

	ourConnect.Port = connections.IP + ":" + connections.Connections[0].Port
	ourConnect.Type = "client"
	ourConnect.Username = userName
	
	connections.Connections = append(connections.Connections, ourConnect)
	
	jsonData, err := json.Marshal(connections)
	if err != nil {
		fmt.Println("Error marshalling JSON")
	}

	ioutil.WriteFile("connections.json", jsonData, os.ModePerm)
	return ourConnect, nil
}

