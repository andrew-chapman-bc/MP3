package unicast

import (
	"../utils"
	"encoding/json"
	"errors"
	"fmt"
	"net"
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
func (cli *Client) SendMessageToServer(messageData utils.Message) (err error) {
	

	jsonData, err := json.Marshal(messageData)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(cli.client)
	encoder.Encode(jsonData)
	fmt.Println("data sent!", messageData)
	return
}
