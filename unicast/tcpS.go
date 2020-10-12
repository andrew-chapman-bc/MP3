package unicast

import (
	"fmt"
	"net"
	"os"
	"io/ioutil"	
	"encoding/gob"
	"errors"
	"../utils"
	"sync"
	"bufio"
)


// Server holds the structure of our TCP server implemenation
type Server struct {
    port   string
	server net.Listener
	Connections utils.Connections
}




/*
	@function: NewTCPServer
	@description: Creates a Server Instance which can then be used in the main function
	@exported: True
	@family: N/A
	@params: string
	@returns: {*Server}, error
*/
func NewTCPServer(port string) (*Server, error) {
	server := Server{port: port}

	// if port is empty -> throw error
	if port == "" {
		return &server, errors.New("Error: Port not found")
	}

    return &server, nil
}

/*
	@function: RunServ
	@description: Starts the TCP server and calls handle connections
	@exported: True
	@family: Server
	@params: chan string, chan bool, waitgroup
	@returns: error
*/
func (serv *Server) RunServ() (err error) {
	
	serv.server, err = net.Listen("tcp", ":" + serv.port)
    if err != nil {
        return err
	}
	fmt.Println("Listening to the port:", serv.port)
	
	defer serv.server.Close()

    for {
		serv.handleConnections(serv.server)
		// break here when calculation is good 
    }
    return
}
/*
	@function: handleConnections
	@description: calls the Accept function in a loop and calls another handleConnection goroutine which decodes data and sends it to the specified client
	@exported: false
	@family: Server
	@params: map[string]net.Conn, chan bool, WaitGroup
	@returns: error
*/
func (serv *Server) handleConnections() (err error) {
	var messages utils.Messages
	ownState, err := fetchInitialState()
	if err != nil {
		return err
	}
	messages.Messages[0] := ownState
	messages.MessagesArr[0] := messages.Messages
	for {
		conn, err := serv.server.Accept()
		
        if err != nil || conn == nil {
            err = errors.New("Network Error: Could not accept connection")
            break
		}

        go serv.handleConnection(conn, messages)
	}
	
    return
}

/*
	@function: handleConnection
	@description: a goroutine which unserializes JSON data and then calls the sendMessageToClient function
	@exported: false
	@family: Server
	@params: net.Conn, map[string]net.Conn
	@returns: error
*/
[
	[
		{State Round},
		{State, Round},
		{State Round},
		{State Round}
	],
	[
		{State Round},
		{State, Round},
		{State Round},
		{State Round}
	] 
]
func (serv *Server) handleConnection(conn net.Conn, messages utils.Messages) (err error) {
	fmt.Println("ok ok")
	nodes, err := utils.GetNodeNums()
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(conn)
	var mess utils.Message
    for {
		fmt.Println(mess)
		err := gob.Decode(&mess)
		if err != nil {
			fmt.Println(err)
			return err
		}
		
		messages.MessagesArr[mess.Round-1].append(messages.MessagesArr[mess.Round-1], mess)
		for i, val := range messages.MessagesArr {
			if len(val) >= (nodes.totalNodes - nodes.faultyNodes) {
				// send this specific val {Message} through channel after it is calculated
			}
		}


    }
}



/*
	@function: readJSONForServer
	@description: Reads the JSON and returns a struct which contains 
		the type, port, username and IP
	@exported: False
	@family: Server
	@params: string
	@returns: Connections
*/
func (serv *Server) fetchInitialState() (utils.Message, error) {
	jsonFile, err := os.Open("connections.json")
	var connections utils.Connections
	var selfState utils.Message
	if err != nil {
		fmt.Println(err)
		return selfState, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	for i := 0; i < len(connections.Connections); i++ {
		if (connections.Connections[i].Port == serv.port ) {
			selfState = CreateMessage(connections.Connections[i].State, 1)
			return selfState, nil
		}
	}

	return selfState, errors.New("Could not find own state?")
}

