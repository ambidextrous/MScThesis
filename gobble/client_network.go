package main

import (
	"log"
	"fmt"
	"net"
	"encoding/gob"
	"time"
)

// Network 

type Transmitter struct {
	Encoder *gob.Encoder
	Decoder *gob.Decoder
	Connection net.Conn
}

func NewTransmitter(conType string, serverAddress string) *Transmitter {
	t := new(Transmitter)
	var connToServer net.Conn
	var err error
	establishingConnection := true
	for establishingConnection {
		connToServer, err = net.Dial(conType, serverAddress)
		if err != nil {
			fmt.Println("Connection error: " + err.Error())
			fmt.Println("Retrying...")
			time.Sleep(time.Second * 3)
		} else {
			establishingConnection = false
			fmt.Println("Successfully established " + conType + " connection with " + serverAddress)
		}
	}
	Encoder := gob.NewEncoder(connToServer)
	Decoder := gob.NewDecoder(connToServer)
	t.Connection = connToServer
	t.Encoder = Encoder
	t.Decoder = Decoder
	return t
}

type Identifier struct {
	Id string
}

type AggregatorClientGob struct {
	Name string
	Accept_from_Client_to_Aggregator_Empty Accept_from_Client_to_Aggregator_Empty
	ConfirmPayment_from_Aggregator_to_Client_bool ConfirmPayment_from_Aggregator_to_Client_bool
	GreetAggregator_from_Client_to_Aggregator_Empty GreetAggregator_from_Client_to_Aggregator_Empty
	GreetClient_from_Aggregator_to_Client_Empty GreetClient_from_Aggregator_to_Client_Empty
	ProvideFlightInformation_from_Aggregator_to_Client_bool_string ProvideFlightInformation_from_Aggregator_to_Client_bool_string
	ProvidePaymentInto_from_Client_to_Aggregator_string ProvidePaymentInto_from_Client_to_Aggregator_string
	RejectAndLeave_from_Client_to_Aggregator_Empty RejectAndLeave_from_Client_to_Aggregator_Empty
	RequestItinerary_from_Client_to_Aggregator_string_int RequestItinerary_from_Client_to_Aggregator_string_int
	RequestPaymentInfo_from_Aggregator_to_Client_Empty RequestPaymentInfo_from_Aggregator_to_Client_Empty
	TryAgain_from_Client_to_Aggregator_Empty TryAgain_from_Client_to_Aggregator_Empty
}

func SendToAggregator(chans *Channels, trans *Transmitter) {
	defer func() {
		trans.Connection.Close()
	}()
	// identifier := Identifier{Id: "Client"} // Comment out line if connecting to Aggregator as server; uncomment line if connecting to Aggregator as client.
	// trans.Encoder.Encode(identifier) // Comment out line if connecting to Aggregator as server; uncomment line if connecting to Aggregator as client.
	for {
		select {
		case out := <-chans.GreetAggregatorFromClientToAggregator_Empty:
			g := AggregatorClientGob{Name: "GreetAggregator_from_Client_to_Aggregator_Empty", GreetAggregator_from_Client_to_Aggregator_Empty: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case out := <-chans.RequestItineraryFromClientToAggregator_string_int:
			g := AggregatorClientGob{Name: "RequestItinerary_from_Client_to_Aggregator_string_int", RequestItinerary_from_Client_to_Aggregator_string_int: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case out := <-chans.TryAgainFromClientToAggregator_Empty:
			g := AggregatorClientGob{Name: "TryAgain_from_Client_to_Aggregator_Empty", TryAgain_from_Client_to_Aggregator_Empty: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case out := <-chans.RejectAndLeaveFromClientToAggregator_Empty:
			g := AggregatorClientGob{Name: "RejectAndLeave_from_Client_to_Aggregator_Empty", RejectAndLeave_from_Client_to_Aggregator_Empty: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case out := <-chans.AcceptFromClientToAggregator_Empty:
			g := AggregatorClientGob{Name: "Accept_from_Client_to_Aggregator_Empty", Accept_from_Client_to_Aggregator_Empty: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case out := <-chans.ProvidePaymentIntoFromClientToAggregator_string:
			g := AggregatorClientGob{Name: "ProvidePaymentInto_from_Client_to_Aggregator_string", ProvidePaymentInto_from_Client_to_Aggregator_string: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case <- chans.doneCommunicatingWithAggregator:
			return
		default:
			// Keep looping
		}
	}
}

func ReceiveFromAggregator(chans *Channels, trans *Transmitter) {
	var in AggregatorClientGob
	for {
		err := trans.Decoder.Decode(&in)
		if err != nil {
			log.Fatal(err)
		}
		if in.Name == "GreetClient_from_Aggregator_to_Client_Empty" {
			chans.GreetClientFromAggregatorToClient_Empty <- in.GreetClient_from_Aggregator_to_Client_Empty
		} else if in.Name == "ProvideFlightInformation_from_Aggregator_to_Client_bool_string" {
			chans.ProvideFlightInformationFromAggregatorToClient_bool_string <- in.ProvideFlightInformation_from_Aggregator_to_Client_bool_string
		} else if in.Name == "RequestPaymentInfo_from_Aggregator_to_Client_Empty" {
			chans.RequestPaymentInfoFromAggregatorToClient_Empty <- in.RequestPaymentInfo_from_Aggregator_to_Client_Empty
		} else if in.Name == "ConfirmPayment_from_Aggregator_to_Client_bool" {
			chans.ConfirmPaymentFromAggregatorToClient_bool <- in.ConfirmPayment_from_Aggregator_to_Client_bool
		} else {
			log.Fatal("ReceiveFromServer() received unknown gob: ", in)
		}
	}
}

func ConnectToAggregatorAsClient(chans *Channels, conType string, serverAddress string) {
	trans := NewTransmitter(conType, serverAddress)
	go SendToAggregator(chans, trans)
	go ReceiveFromAggregator(chans, trans)
}

func CloseConnectAsClientWithAggregator(chans *Channels) {
	chans.doneCommunicatingWithAggregator<- true
}

func AcceptConnections(conType string, port string, chans *Channels) {
	ln, err := net.Listen(conType, port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		HandleConnection(conn, chans)
	}
}

func HandleConnection(conn net.Conn, chans *Channels) {
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)
	trans := &Transmitter{Decoder: decoder, Encoder: encoder, Connection: conn}
	var identifier Identifier
	err := trans.Decoder.Decode(&identifier)
	if err != nil {
		log.Fatal(err)
	}
	if identifier.Id == "Aggregator" {
		HandleAggregatorConnectionAsServer(trans, chans)
	} else {
		log.Fatal("HandleConnection received unknown identifier: ", identifier)
	}
}

func HandleAggregatorConnectionAsServer(trans *Transmitter, chans *Channels) {
	go SendToAggregator(chans, trans)
	go ReceiveFromAggregator(chans, trans)
}

func SetupNetworkConnections(chans *Channels, connType string, address string,port string) {
	go AcceptConnections(connType, port, chans) // Comment out to stop accepting connections as Server
	// ConnectToAggregatorAsClient(chans, connType, address + port) // Uncomment to connect as client
}

func CloseNetworkConnections(chans *Channels) {
	chans.doneCommunicatingWithAggregator <- true
	<-chans.doneCommunicatingWithAggregator
}

