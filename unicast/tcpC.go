package unicast

import (
	"../utils"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"io/ioutil"
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
	@params: string, connections
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
	@description: Starts the TCP client (Dials)
	@exported: True
	@family: Client
	@params: N/A
	@returns: error
*/
func (cli *Client) RunCli() (err error) {
	cli.client, err = net.Dial("tcp", cli.Connections.IP + ":" + cli.port)
	if err != nil {
		return err
	}

	return nil

}

/*
	@function: SendMessageToServer
	@description: Reads the message channel and sends the data over to server using GOB
	@exported: True
	@family: Client
	@params: chan {Message}
	@returns: error
*/
func (cli *Client) SendMessageToServer(messageData utils.Message) (err error) {
	delay, err := utils.GetDelayParams()
	if err != nil {
		fmt.Println(err)
		return err
	}
	utils.GenerateDelay(delay)
	encoder := gob.NewEncoder(cli.client)
	encoder.Encode(messageData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println("Sent data to ", cli.port, messageData)
	return
}

/*
	@function: FetchInitialState
	@description: Returns the initial state of the active node
	@exported: True
	@family: Client
	@params: N/A
	@returns: Message, error
*/
func (cli *Client) FetchInitialState() (utils.Message, error) {
	jsonFile, err := os.Open("config.json")
	var connections utils.Connections
	var selfState utils.Message
	if err != nil {
		fmt.Println(err)
		return selfState, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	for i := 0; i < len(connections.Connections); i++ {
		if (connections.Connections[i].Port == cli.port ) {
			selfState = utils.CreateMessage(connections.Connections[i].State, 1)
			return selfState, nil
		}
	}

	return selfState, errors.New("Could not find own state?")
}

