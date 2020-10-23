# MP3
Approx Consensus
--- 
# To Run

```bash
./startClients.sh
``` 
This will start up the controller and the clients
---

# Structure and design

Controller:

Clients:
Conncurrent TCP Server
Each client has their own conncurrent TCP server
There us one struct desined to make passing data easier and readable
The Server struct holds the structure of our TCP server implemenation
The port it is listening to, the server, and all the client connections
```
type Server struct {
  port   string
	server net.Listener
	Connections utils.Connections
}
```



# json file
The json file has the following format 
-----------------------------------------------------------------------------------------------
```    
{"connections":
[{"State":".1234",
   "Port":"1234",
    "Status":"5"},
{"State":".4567","Port":"4567","Status":"2"},
{"State":".8543","Port":"8543","Status":"3"},
{"State":".1432","Port":"1432","Status":"4"}],
"IP":"127.0.0.1",
"Delay":{"minDelay": 1000, "maxDelay": 5000},
"Consensus":false,"Round":1}
```
.... .... .......
-----------------------------------------------------------------------------------------------
To read the json file, there are two functions.
One function for the server reading, and one for the client


To add more connection, simply open a new terminal and run the program

For example:
-----------------------------------------------------------------------------------------------  
```  
{
    "connections": [
        {
            "Type": "server",
            "Port": "1234",
            "Username": "Matt"
        },
        {
            "Type": "client",
            "Port": "1234",
            "Username": "Andy"
        },
        {
            "Type": "client",
            "Port": "1234",
            "Username": "Lewis"
        }
    ],
    "IP": "127.0.0.1"
}


```
-----------------------------------------------------------------------------------------------

Goes to 

-----------------------------------------------------------------------------------------------   
``` 
{
    "connections": [
        {
            "Type": "server",
            "Port": "1234",
            "Username": "Matt"
        },
        {
            "Type": "client",
            "Port": "1234",
            "Username": "Andy"
        },
        {
            "Type": "client",
            "Port": "1234",
            "Username": "Lewis"
        },
        {
            "Type": "client",
            "Port": "1234",
            "Username": "Darius"
        }
    ],
    "IP": "127.0.0.1"
}

```
-----------------------------------------------------------------------------------------------

If you run the program again with the username "Darius"


# Input
The user input is broken up into three sections, : 
1. "Send"
2. Username 
3. Message

The program reads each section as follows: 
1. Disregard this keyword
2. Store the username into Message struct 
3. Store the message into Message struct

We are communicating between the server and client using channels 

# Exit Condition 

If the user enters "EXIT" the program will terminate its connection
The user will no longer be able to send/recieve messages

# Processes
The processes can be found in the TCP directory




### Shortcomings and Potential Improvemnts 
As of right now, the program does not run
