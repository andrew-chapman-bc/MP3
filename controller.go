package main

import (
	"./utils"
	"fmt"
	"time"
	"os/exec"
	"encoding/json"
	"io/ioutil"
)

func main() {
	// get start time of program
	startTime := time.Now()
	consensus := false
	var err error
	var connections utils.Connections
	connections, err = utils.GetConnections()
	if err != nil {
		fmt.Println(err)
	}


	// continuously loop through checking JSON for consensus
	for !consensus {
		consensus, err = utils.GetJSONConsensus()
		if err != nil {
			fmt.Println(err)
		}
	}

	// Once we've reached consensus get time and round
	jsonRound, err2 := utils.GetJSONRound()
	if err2 != nil {
		fmt.Println(err2)
	}

	// Print out relevant data
	elapsed := time.Since(startTime)
	fmt.Println("The round of the JSON is:", jsonRound)
	fmt.Println("The time elapsed is: ", elapsed)

	// Execute stop script to kill TCP connections
	_, err = exec.Command("/bin/sh", "./killAndy.sh").Output()
	if err != nil {
		fmt.Println(err)
	}

	// just doing this so we don't have to manually edit the json file every time
	jsonData, err1 := json.Marshal(connections)
	if err1 != nil {
		fmt.Println(err1)
	}
	err1 = ioutil.WriteFile("config.json", jsonData, 0777)
	if err1 != nil {
		fmt.Println(err)
	}
	

}