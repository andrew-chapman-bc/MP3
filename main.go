package main

import (
	"./unicast"
	"bufio"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"strings"
	"sync"
	"errors"
)





func getCmdLine() string {
	// Use argparse library to get accurate command line data
	parser := argparse.NewParser("", "Concurrent TCP Channels")
	s := parser.String("s", "string", &argparse.Options{Required: true, Help: "String to print"})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}
	return *s
}

func main() {
	
	s := getCmdLine()
	connections, err := utils.GetConnections()
	if err != nil {
		fmt.Println("Error reading json", err)
	}
	var serv *unicast.Server
	port := cmdLineArr[1]
	serv, err := unicast.NewTCPServer(port)
	if err != nil {
		fmt.Println(err)
	}


}

