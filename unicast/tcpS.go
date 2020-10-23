package unicast

import (
	"../utils"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
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
	@params: string, connections
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
	@params: chan Message, chan bool
	@returns: error
*/
func (serv *Server) RunServ(messageChannel chan utils.Message, serverFinished chan bool) ( err error) {
	serv.server, err = net.Listen("tcp", ":" + serv.port)
    if err != nil {
		fmt.Println("we did not connect")
        return err
	}
	serverFinished <- true
	// fmt.Println("Listening to the port:", serv.port)
	
	//defer serv.server.Close()

    for {
		serv.handleConnections(serv.server, messageChannel)
    }
    return
}

/*
	@function: handleConnections
	@description: calls the Accept function in a loop and calls another handleConnection goroutine which decodes data via Gob
	@exported: false
	@family: Server
	@params: net.Listener, chan Message
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
	@description: a goroutine which decodes data via gob and sends it over a channel
		which then is used to get the new state
	@exported: false
	@family: Server
	@params: net.Conn, Messages, chan Message
	@returns: error
*/
func (serv *Server) handleConnection(conn net.Conn, messagesArr utils.Messages, messageChannel chan utils.Message) (err error) {
	// fmt.Println("ok ok")
	var mess utils.Message
    for (err != io.EOF) {
		dec := gob.NewDecoder(conn)
		err = dec.Decode(&mess)
		// fmt.Println("Received message:", mess)

		if err != nil {
			fmt.Println(err)
			return err
		}

		messageChannel <- mess
	}
	return
}


