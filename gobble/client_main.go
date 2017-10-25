package main

import (
	"sync"
	"log"
	"fmt"
)

func runClient_rec1_choice1_A(Client_rec1_choice1 *Client_rec1_choice1) (interface{}, error) {
	Client_rec1_1_new, err1 := Client_rec1_choice1.Send_TryAgain()
	if err1 != nil {
		log.Fatal(err1)
	}
	return Client_rec1_1_new, nil
}

func runClient_rec1_choice1_B(Client_rec1_choice1 *Client_rec1_choice1) (interface{}, error) {
	err1 := Client_rec1_choice1.Send_RejectAndLeave()
	if err1 != nil {
		log.Fatal(err1)
	}
	return "", nil
}

func runClient_rec1_choice1_C(Client_rec1_choice1 *Client_rec1_choice1) (interface{}, error) {
	Client_rec1_choice1_C2, err1 := Client_rec1_choice1.Send_Accept()
	if err1 != nil {
		log.Fatal(err1)
	}
	Client_rec1_choice1_C3, err2 := Client_rec1_choice1_C2.Receive_RequestPaymentInfo()
	if err2 != nil {
		log.Fatal(err2)
	}
	var sending_3_1_string string
	Client_rec1_choice1_C4, err3 := Client_rec1_choice1_C3.Send_ProvidePaymentInto_string(sending_3_1_string)
	if err3 != nil {
		log.Fatal(err3)
	}
	received_4_1_bool, err4 := Client_rec1_choice1_C4.Receive_ConfirmPayment_bool()
	if err4 != nil {
		log.Fatal(err4)
	}
	fmt.Println("Client_rec1_choice1_C4 received value type: bool: ", received_4_1_bool)
	return "", nil
}

func runClient_rec1(Client_rec1_1 *Client_rec1_1) (interface{}, error) {
	var sending_1_1_string string
	var sending_1_2_int int
	Client_rec1_2, err1 := Client_rec1_1.Send_RequestItinerary_string_int(sending_1_1_string, sending_1_2_int)
	if err1 != nil {
		log.Fatal(err1)
	}
	Client_rec1_choice1, received_2_1_bool, received_2_2_string, err2 := Client_rec1_2.Receive_ProvideFlightInformation_bool_string()
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("Client_rec1_2 received value type: bool: ", received_2_1_bool)
	fmt.Println("Client_rec1_2 received value type: string: ", received_2_2_string)
	retVal, err3 := Make_Client_rec1_choice1_Choices(Client_rec1_choice1)
	if err3 != nil {
		log.Fatal(err3)
	}
	return retVal, nil
}

func runClient(wg *sync.WaitGroup, Client1 *Client1) (interface{}, error) {
	defer wg.Done()
	Client2, err1 := Client1.Send_GreetAggregator()
	if err1 != nil {
		log.Fatal(err1)
	}
	Client_rec1_1, err2 := Client2.Receive_GreetClient()
	if err2 != nil {
		log.Fatal(err2)
	}
	retVal, err3 := LoopClient_rec1(Client_rec1_1)
	if err3 != nil {
		log.Fatal(err3)
	}
	return retVal, nil
}

// Functions

func Make_Client_rec1_choice1_Choices(Client_rec1_choice1 *Client_rec1_choice1) (interface{}, error) {
	retVal, err := runClient_rec1_choice1_A(Client_rec1_choice1)
	//retVal, err := runClient_rec1_choice1_B(Client_rec1_choice1) // Uncomment to use
	//retVal, err := runClient_rec1_choice1_C(Client_rec1_choice1) // Uncomment to use
	return retVal, err
}

func StartClient(chans *Channels) (*Client1) {
	start := &Client1{Channels: chans}
	return start
}

func LoopClient_rec1(Client_rec1_1_old *Client_rec1_1) (interface{}, error) {
	var retVal interface{}
	var err error
	looping := true
	for looping {
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
	chans := NewChannels()
	connType := "tcp"
	address := "localhost"
	port := ":8080"
	SetupNetworkConnections(chans, connType, address, port)
	startStruct := StartClient(chans)
	var newWg sync.WaitGroup
	newWg.Add(1)
	go runClient(&newWg, startStruct)
	newWg.Wait()
	CloseNetworkConnections(chans)
}

