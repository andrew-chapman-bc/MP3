# MP3
Approx Consensus
--- 
# To Run

```bash
./startClients.sh
``` 
This will start up the controller and the clients
The output will be one round of testing
It will print out rounds to reach consensus 
Will also print out time it took to run

To change the number of nodes, inital states, faulty states, or delay parameter
Change the respective fields in the config.json file

---

# Structure and design

Use gob to encode/decode messages over TCP ports
Use JSON for config file 

Controller:
The controller waits for the JSON file to show consensus
When it starts, it takes the time
Once consensus is reached, 
it will print number of rounds it took and time the program took to run
It also attempts to reset JSON to original state.
There is an error where it might delete the delay parameters

Utils:
Filled with all the helper functions
Called from every other file in project

Clients:
Main
Gets every connection
Creates a conncurrent TCP server
Create client side
Send out messages to other connections
Receive messages from other connections
Check for consensus

Client
Creates iniital state
Dials into otehr servers
Sends messages

Conncurrent TCP Server
Listens for messages to come in 
Sends them over channel back to main



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
The Inputs are given in the config file and are preset

# Delay
The Delay is given in the JSON file
It is a range that the program picks some time in that range to delay its message

We are communicating between the server and client using channels 

# Exit Condition 
once consensus is reached, the controller runs a kill script to end all processes
Consensus is known when the consensus field in JSON turns true

# Processes
The processes can be found in the TCP directory




### Shortcomings and Potential Improvemnts 
One way to get better data for the analysis would be to instead of printing out the time and rounds
to write them to a file which could be easily read and converted to box plot i.e. csv

Our style of TCP would be more simialr to a UDP style. Instead of reconnecting multiple times, 
we should set up one connnection and keep it the same for the whole process
