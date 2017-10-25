package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func runServer_rec1(Server_rec1_1 *Server_rec1_1) (interface{}, error) {
	Server_rec1_2, _, err1 := Server_rec1_1.Receive_firstStep_int()
	if err1 != nil {
		log.Fatal(err1)
	}
	var sending_2_1_int int
	Server_rec1_1_new, err2 := Server_rec1_2.Send_secondStep_int(sending_2_1_int)
	if err2 != nil {
		log.Fatal(err2)
	}
	return Server_rec1_1_new, nil
}

func runServer(wg *sync.WaitGroup, Server_rec1_1 *Server_rec1_1) (interface{}, error) {
	defer wg.Done()
	retVal, err1 := LoopServer_rec1(Server_rec1_1)
	if err1 != nil {
		log.Fatal(err1)
	}
	return retVal, nil
}

func StartServer(chans *Channels) *Server_rec1_1 {
	start := &Server_rec1_1{Channels: chans}
	return start
}

func LoopServer_rec1(Server_rec1_1_old *Server_rec1_1) (interface{}, error) {
	var retVal interface{}
	var err error
	looping := true
	for looping {
		retVal, err = runServer_rec1(Server_rec1_1_old)
		if err != nil {
			log.Fatal(err)
		}
		switch t := retVal.(type) {
		case *Server_rec1_1:
			Server_rec1_1_old = t
		default:
			looping = false
		}
	}
	return retVal, nil
}

func runClient_rec1(Client_rec1_1 *Client_rec1_1) (interface{}, error) {
	var sending_1_1_int int
	Client_rec1_2, err1 := Client_rec1_1.Send_firstStep_int(sending_1_1_int)
	if err1 != nil {
		log.Fatal(err1)
	}
	Client_rec1_1_new, _, err2 := Client_rec1_2.Receive_secondStep_int()
	if err2 != nil {
		log.Fatal(err2)
	}
	return Client_rec1_1_new, nil
}

func runClient(wg *sync.WaitGroup, Client_rec1_1 *Client_rec1_1) (interface{}, error) {
	defer wg.Done()
	retVal, err1 := LoopClient_rec1(Client_rec1_1)
	if err1 != nil {
		log.Fatal(err1)
	}
	return retVal, nil
}

func StartClient(chans *Channels) *Client_rec1_1 {
	start := &Client_rec1_1{Channels: chans}
	return start
}

func LoopClient_rec1(Client_rec1_1_old *Client_rec1_1) (interface{}, error) {
	var retVal interface{}
	var err error
	var counter int
	looping := true
	for looping {
		if counter == 1000000 {
			return "", nil
		} else {
			counter++
		}
		retVal, err = runClient_rec1(Client_rec1_1_old)
		if err != nil {
			log.Fatal(err)
		}
		switch t := retVal.(type) {
		case *Client_rec1_1:
			Client_rec1_1_old = t
		default:
			looping = false
		}
	}
	return retVal, nil
}

func main() {
	startTime := time.Now()
	chans := NewChannels()
	startStructServer := StartServer(chans)
	startStructClient := StartClient(chans)
	var newWg sync.WaitGroup
	go runServer(&newWg, startStructServer)
	newWg.Add(1)
	go runClient(&newWg, startStructClient)
	newWg.Wait()
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)
	fmt.Println("Execution time = ", totalTime)
}
