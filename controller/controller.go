package controller

import (
	"../utils"
	"fmt"
	"time"
	"os/exec"
)

func main() {
	startTime := time.Now()
	consensus := false
	var err error

	for !consensus {
		consensus, err = utils.GetJSONConsensus()
		if err != nil {
			fmt.Println(err)
		}
	}
	jsonRound, err2 := utils.GetJSONRound()
	if err2 != nil {
		fmt.Println(err2)
	}
	elapsed := time.Since(startTime)
	fmt.Println("The round of the JSON is:", jsonRound)
	fmt.Println("The time elapsed is: ", elapsed)
	_, err = exec.Command("/bin/sh", "./killAndy.sh").Output()
	if err != nil {
		fmt.Println(err)
	}


	

	// just doing this so we don't have to manually edit the json file every time
	utils.SetJSONRound(1)
	utils.SetJSONConsensus(false)
}