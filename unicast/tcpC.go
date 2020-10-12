package unicast

import (
	"fmt"
	"errors"
	"net"
	"os"
	"io/ioutil"
	"encoding/json"
	"../util"
	"strings"
	"sync"
)


// Client holds the structure of our TCP Client implementation
type Client struct {
	Username string
	client net.Conn
}

/*
	@function: NewTCPClient
	@description: Creates a Client instance which can be used in the main function
	@exported: True
	@family: N/A
	@params: string
	@returns: {*Client}, error
*/
func NewTCPClient(username string) (*Client, error) {
	client := Client{Username: username}
	// if username is empty -> throw error
	if username == "" {
		fmt.Println("error here 1")
		return nil, errors.New("Error: Address not found")
	}

	return &client, nil
}


func (cli *Client) sendUserToServer() (err error){
	fmt.Fprintf(cli.client, cli.Username + "\n")
	return
}

/*
	@function: RunCli
	@description: Starts the TCP client which calls the function to send message to server
	@exported: True
	@family: Client
	@params: chan {Message}
	@returns: error
*/
func (cli *Client) RunCli(messageChan chan util.Message, ) (err error) {
	
	connection, err := cli.readJSONForClient(cli.Username)
	if err != nil {
		return err
	}
	fmt.Println(connection.Port)
	cli.client, err = net.Dial("tcp", connection.Port)
	if err != nil {
		return err
	}

	for {
		messageData := <- messageChan
		if messageData.Message == "EXIT" {
			wg.Done()
			break
		} else if messageData.Receiver == "Client Not Connected" {
			fmt.Println("Client not connected")
		}
		fmt.Println(messageData)
		go cli.sendMessageToServer(cli.client, messageData)
	}
	wg.Done()
	fmt.Println("outside the for loop here now")
	wg.Wait()
	return

}

/*
	@function: sendMessageToServer
	@description: Reads the message channel and serializes the data to send over to server
	@exported: false
	@family: Client
	@params: net.Conn, chan {Message}
	@returns: error
*/
func (cli *Client) sendMessageToServer(conn net.Conn, messageData util.Message) (err error) {
	

	jsonData, err := json.Marshal(messageData)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(conn)
	encoder.Encode(jsonData)
	fmt.Println("data sent!", messageData)
	return
}

/*
	@function: listenForMessage
	@description: Listens for a message from the server and deserializes it 
	@exported: false
	@family: Client
	@params: net.Conn, chan {Message}
	@returns: error
*/
func (cli *Client) listenForMessage(conn net.Conn, messageChan chan util.Message) (err error) {
	for {
		decoder := json.NewDecoder(conn)
		var mess util.Message
		decoder.Decode(&mess)

		if mess.Message == "error" {
			messageChan <- util.Message{"Client Not Connected", "", ""}
			return errors.New("Person not connected yet")
		} else if mess.Message == "EXIT" {
			conn.Close()
			os.Exit(0)
			messageChan <- util.Message{"","EXIT",""}
		} else if mess.Message != "" {
			fmt.Printf("Received the message from" + strings.TrimSpace(mess.Sender) + "\n") 
			fmt.Printf("Message:" + strings.TrimSpace(mess.Message))
		}
	}
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

