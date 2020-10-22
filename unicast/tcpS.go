package unicast

import (
	"../utils"
	"encoding/gob"
	//"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	//"os"
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
func NewTCPServer(port string, connections utils.Connections) (*Server, error) {
	server := Server{port: port, Connections: connections}

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
func (serv *Server) RunServ(messageChannel chan utils.Message) ( err error) {
	serv.server, err = net.Listen("tcp", ":" + serv.port)
    if err != nil {
		fmt.Println("we did not connect")
        return err
	}
	fmt.Println("Listening to the port:", serv.port)
	
	//defer serv.server.Close()

    for {
		serv.handleConnections(serv.server, messageChannel)
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
func (serv *Server) handleConnections(conn net.Listener, messageChannel chan utils.Message) (err error) {
	var messagesArr utils.Messages
	for {
		conn, err := serv.server.Accept()
        if err != nil || conn == nil {
			err = errors.New("Network Error: Could not accept connection")
			fmt.Println(err)
            break
		}

        go serv.handleConnection(conn, messagesArr, messageChannel)
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
/*
[
	[
		{ownState 1Round},
		
	],
	[
		{State Round2},
		{State, Round2},
		{State Round2},
		{State Round2}
	] 
]
*/
func (serv *Server) handleConnection(conn net.Conn, messagesArr utils.Messages, messageChannel chan utils.Message) (err error) {
	fmt.Println("ok ok")
	nodes, err := utils.GetNodeNums()
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(conn)
	var mess utils.Message
	var messages utils.Messages
	counter := 0
    for (err != io.EOF) {
		err = dec.Decode(&mess)
		fmt.Println("this is the message", mess)
		// .1234, 1
		if err != nil {
			fmt.Println("hit this: ", err)
			return err
		}
		messages.Messages = append(messages.Messages, mess)
		fmt.Println("messages array is", messagesArr)
		//TODO: stuff
	}
	return
}


