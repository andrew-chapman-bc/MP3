# MP3: Approx Consensus
--- 

https://docs.google.com/spreadsheets/d/1p_94uyj_fO1ug_x-gX9nW4s4NkYB5qco-zsgbmmXE6o/edit#gid=0

--- 
As we added more nodes, it had the most significant impact on the time it took to achieve consensus.  Furthemore, latency played the next largest factor, however it was random so it also had one of the largest ranges.  The range of values did not seem to make that large of an impact when we were running these tests.  During our trials, it became obvious that the most common amount of rounds it took for all nodes to reach approximate consensus was 3.
# To Run

```bash
./startClients.sh
``` 
This will start up the controller and the clients.
The output will be one round of testing.
It will print out rounds to reach consensus 
Will also print out time it took to run.

To change the number of nodes, inital states, faulty states, or delay parameter
Change the respective fields in the config.json file

---

# Structure and design

* Use gob to encode/decode messages over TCP ports
* Use JSON for config file 

## Controller
The controller waits for the JSON file to show consensus
When it starts, it creates a variable for the time
Once consensus is reached, 
it will print number of rounds it took and time the program took to run
It also attempts to reset JSON to original state.
There is an error where it might delete the delay parameters

### Packages
```  
import
(
    "./utils"
	"fmt"
	"time"
	"os/exec"
	"encoding/json"
	"io/ioutil"
)
``` 
* utils -> Custom Utility package 
* json -> Encode and decode JSON
* ioutil -> Read from files
* fmt -> Standard for printing, etc.
* time -> Getting timestamps
* exec -> Executing shell commands 

## Utils
Holds various functions used for utility, such as reading and writing JSON, creating structs, getting amount of nodes, calculating the consensus and delays for latency simulation.

### Packages
```  
import
(
	"os"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"math"
	"fmt"
	"math/rand"
	"time"
	"strings"
)
```  

* os -> Open files
* json -> Encode and decode JSON
* errors -> Error handling
* ioutil -> Read from files
* strconv -> Convert from various types
* math -> Absolute value for consensus logic
* fmt -> Standard for printing, etc.
* rand -> Getting random latency
* time -> Getting timestamps
* strings -> String manipulation / Checking if strings contain substrings
### Structs
```  
/*
    This Contains our complete config.JSON
    We initially call a function to create an instance of this in
    main.go, then pass it to create a TCPClient and TCPServer struct
    which have it as a field
*/
type Connections struct {
	Connections []Connection `json:"connections"`
	IP string `json:"IP"`
	Delays Delay `json:"Delay"`
	Consensus bool `json:"Consensus"`
	Round int `json:"Round"`
} 

/*
	State: State of the Node
	Port: "1234", etc. Port attached to Node
	Status: Whether or not the node is faulty
    This holds each node's values
*/
type Connection struct {
	State string `json:"State"`
	Port string `json:"Port"`
	Status string `json:"Status"`
}

/*
	State: The state of the message
	Round: What round the message is sent in
    This struct contains our message that we are sending across GOB
*/
type Message struct {
	State string
	Round int
}
/*
	Messages: An array of message
*/
type Messages struct {
	Messages []Message
}
/*
	TotalNodes: Total nodes in our distributed system
	FaultyNodes: Total amt of faulty nodes in our sys
	(We get these from the JSON config file)
    We need these for n-f logic around the program
*/
type NodeNums struct {
	TotalNodes int
	FaultyNodes int
}
/*
	minDelay: lower bound for delay
	maxDelay: upper bound for delay
    Used to simulate latency
*/
type Delay struct {
	MinDelay int `json:"MinDelay"`
	MaxDelay int `json:"MaxDelay"`
}
```  


## TCP Client
* Creates an instance of a TCP Client to a port
* Dials to a port
* Sends a message to whichever server the client is connected to using GOB
* Fetches initial state of the Node
### Packages
```  
import 
(
	"../utils"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"io/ioutil"
)
``` 
* utils -> Custom Utility package 
* gob -> Encode and decode GOB
* json -> Encode and decode JSON
* errors -> Error handling
* fmt -> Standard for printing, etc.
* net -> Dialing function for TCP
* os -> Open files
* ioutil -> Read from files
### Structs
```  
/* 
//  Port: Port the Client (Node) is dialed to
//	Client: Connection to TCP Server
//	Connections: Global JSON which holds program data

//  This struct is used to hold an instance of the dial
//  connection to all the other nodes
*/

type Client struct {
	port string
	client net.Conn
	Connections utils.Connections
}
``` 
## TCP Server
* Creates an instance of a TCP Server to a port
* Listens to a port
* Handles connections using the Accept function from net package
* Goroutine to handle connections from multiple dials
* Listens for messages to come in 
* Sends them over channel back to main

### Packages
```  

import 
(
	"../utils"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"
)
``` 
* utils -> Custom Utility package 
* gob -> Encode and decode GOB
* errors -> Error handling
* fmt -> Standard for printing, etc.
* net -> Listener function for TCP
* io -> Check if end of file
### Structs
```  
/* 
//  port: Port the Server (Node) is dialed to
//  server: Connection to TCP Server
//	Connections: Global JSON which holds program data

//  This struct is used to hold an instance of the listener
//  connection to all the other nodes
*/

type Server struct {
	port string
	server net.Listener
	Connections utils.Connections
}
``` 




## Main
* Gets command line arguments
* Reads JSON to pass to Server and Client instances
* Creates a concurrent TCP server for Node
* Makes sure server is created before we dial to all nodes
* Create TCP client to each Node in the system
* Contains rounds and checks which node will not send a message (faulty)
* Send out messages to other nodes
* Receive messages from other nodes
* Check for approximate consensus and then if reached sets JSON consensus to true

### Packages
```  
import 
(
	"./unicast"
	"./utils"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"sync"
	"time"
	"strings"
	"strconv"
)

``` 
* unicast -> Custom unicast package which holds client and server functionality
* utils -> Custom Utility package 
* fmt -> Standard for printing, etc.
* argparse -> Used to parse user input
* os -> Open files
* sync -> Waitgroups
* time -> Getting timestamps
* strings -> String Manipulation
* strconv -> Converting between Go Types


## JSON
The JSON file has the following format 
To add more nodes, simply add a node's data to the connections array 
-----------------------------------------------------------------------------------------------
```    
{
  "connections": [
    {
      "State":".1234",
      "Port":"1234",
      "Status":"5"
    },
    {
      "State":".1234",
      "Port":"1234",
      "Status":"5"
    },
    {
      "State":".1234",
      "Port":"1234",
      "Status":"5"
    },
    {
      "State":".1234",
      "Port":"1234",
      "Status":"5"
    },
  ],
  "IP":"127.0.0.1",
  "Delay":
    {
      "minDelay": 1000, 
      "maxDelay": 5000
    },
  "Consensus":false,
  "Round":1
}



## Exit Condition 
once consensus is reached, the controller runs a kill script to end all processes
Consensus is known when the consensus field in JSON turns true



### Shortcomings and Potential Improvemnts 
One way to get better data for the analysis would be to instead of printing out the time and rounds
to write them to a file which could be easily read and converted to box plot i.e. csv

Our style of TCP would be more similar to a UDP style. Instead of reconnecting multiple times, 
we should set up one connnection and keep it the same for the whole process.
We also could've abstracted our utility.go file more as there is function clutter.
Additionally, had we had more time we would've found a better way to close connections
where we wouldn't have to use shell scripts.  Additionally, right now our actual approxiate
consensus model is a little wrong.  This is because we don't currently hold all the node's
states, however just the non-faulty nodes.
