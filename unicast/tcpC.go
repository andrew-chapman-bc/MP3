package unicast

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// UserInput holds user input such as message, destination and source
type UserInput struct {
	Destination string
	Message     string
	Source 		string
}

// Delay keeps track of delay bounds from config
type Delay struct {
	minDelay string
	maxDelay string
}
// Connection holds the ip/port and source
type Connection struct {
	ip 		string
	port 	string
	source 	string
}

/*
	@function: ScanConfigForClient
	@description: Scans the config file using the user input destination and retrieves the ip/port that will later be used to connect to the TCP server
	@exported: True
	@params: {userInput} 
	@returns: {Connection}
*/
func ScanConfigForClient(userInput UserInput) Connection {

	destination := userInput.Destination
	
	// Open up config file
	// TODO: create a variable which holds the destination of config file instead of hardcoding here
	config, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(config)
	scanner.Split(bufio.ScanLines)
	var connection Connection
	counter := 0
	for {
		success := scanner.Scan()
		if success == false {
			err = scanner.Err()
			if err == nil {
				break
			} else {
				log.Fatal(err)
				break
			}
		}
		if counter != 0 {
			// TODO: should do some more error handling here to make sure they are accurate ports/ips in the config
			configArray := strings.Fields(scanner.Text())
			if configArray[0] == destination {
				connection.ip = configArray[1]
				connection.port = configArray[2]
				connection.source = userInput.Source
			}
		}
		counter++
	}
	return connection
} 

/*
	@function: connectToTCPServer
	@description:	Connects to the TCP server with the ip/port obtained from config file as a parameter and 
					returns the connection to the server which will later be used to write to the server
	@exported: false
	@params: string 
	@returns: net.Conn, err
*/
func connectToTCPServer(connect string) (net.Conn, error) {
	// Dial in to the TCP Server, return the connection to it
	c, err := net.Dial("tcp", connect)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	return c, err
} 

/*
	@function: SendMessage
	@description: 	SendMessage sends the message from TCPClient to TCPServer by connecting to the server and 
					using the Fprintf function to send the message.
	@exported: True
	@params: {UserInput}, {Connection}
	@returns: N/A
*/
func SendMessage( messageParams UserInput, connection Connection ) {
	connectionString := connection.ip + ":" + connection.port
	c, err := connectToTCPServer(connectionString)
	if (err != nil) {
		fmt.Println("Network Error: ", err)
	}
	
	if (err != nil) {
		fmt.Println("Error: ", err)
	}
	
	// Sending the message to TCP Server
	// Easier to send this over as strings since it is only one message, we want the source to know where it comes from
	fmt.Fprintf(c, messageParams.Message + " " + messageParams.Source + "\n")
	timeOfSend := time.Now().Format("02 Jan 06 15:04:05.000 MST")
	fmt.Println("Sent message " + messageParams.Message + " to destination " + messageParams.Destination + " system time is: " + timeOfSend)
	
} 

