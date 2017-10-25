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

type AggregatorQueasyJetGob struct {
	Name string
	CheckAvailabilityAndPrice2_from_Aggregator_to_QueasyJet_string_int CheckAvailabilityAndPrice2_from_Aggregator_to_QueasyJet_string_int
	ConfirmAvailabilityAndPrice2_from_QueasyJet_to_Aggregator_bool_int ConfirmAvailabilityAndPrice2_from_QueasyJet_to_Aggregator_bool_int
}

type AggregatorBrutishAirwaysGob struct {
	Name string
	CheckAvailabilityAndPrice1_from_Aggregator_to_BrutishAirways_string_int CheckAvailabilityAndPrice1_from_Aggregator_to_BrutishAirways_string_int
	ConfirmAvailabilityAndPrice1_from_BrutishAirways_to_Aggregator_bool_int ConfirmAvailabilityAndPrice1_from_BrutishAirways_to_Aggregator_bool_int
}

func SendToClient(chans *Channels, trans *Transmitter) {
	defer func() {
		trans.Connection.Close()
	}()
	identifier := Identifier{Id: "Aggregator"} // Comment out line if connecting to Client as server; uncomment line if connecting to Client as client.
	trans.Encoder.Encode(identifier) // Comment out line if connecting to Client as server; uncomment line if connecting to Client as client.
	for {
		select {
		case out := <-chans.GreetClientFromAggregatorToClient_Empty:
			g := AggregatorClientGob{Name: "GreetClient_from_Aggregator_to_Client_Empty", GreetClient_from_Aggregator_to_Client_Empty: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case out := <-chans.ProvideFlightInformationFromAggregatorToClient_bool_string:
			g := AggregatorClientGob{Name: "ProvideFlightInformation_from_Aggregator_to_Client_bool_string", ProvideFlightInformation_from_Aggregator_to_Client_bool_string: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case out := <-chans.RequestPaymentInfoFromAggregatorToClient_Empty:
			g := AggregatorClientGob{Name: "RequestPaymentInfo_from_Aggregator_to_Client_Empty", RequestPaymentInfo_from_Aggregator_to_Client_Empty: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case out := <-chans.ConfirmPaymentFromAggregatorToClient_bool:
			g := AggregatorClientGob{Name: "ConfirmPayment_from_Aggregator_to_Client_bool", ConfirmPayment_from_Aggregator_to_Client_bool: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case <- chans.doneCommunicatingWithClient:
			return
		default:
			// Keep looping
		}
	}
}

func SendToQueasyJet(chans *Channels, trans *Transmitter) {
	defer func() {
		trans.Connection.Close()
	}()
	identifier := Identifier{Id: "Aggregator"} // Comment out line if connecting to QueasyJet as server; uncomment line if connecting to QueasyJet as client.
	trans.Encoder.Encode(identifier) // Comment out line if connecting to QueasyJet as server; uncomment line if connecting to QueasyJet as client.
	for {
		select {
		case out := <-chans.CheckAvailabilityAndPrice2FromAggregatorToQueasyJet_string_int:
			g := AggregatorQueasyJetGob{Name: "CheckAvailabilityAndPrice2_from_Aggregator_to_QueasyJet_string_int", CheckAvailabilityAndPrice2_from_Aggregator_to_QueasyJet_string_int: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case <- chans.doneCommunicatingWithQueasyJet:
			return
		default:
			// Keep looping
		}
	}
}

func SendToBrutishAirways(chans *Channels, trans *Transmitter) {
	defer func() {
		trans.Connection.Close()
	}()
	identifier := Identifier{Id: "Aggregator"} // Comment out line if connecting to BrutishAirways as server; uncomment line if connecting to BrutishAirways as client.
	trans.Encoder.Encode(identifier) // Comment out line if connecting to BrutishAirways as server; uncomment line if connecting to BrutishAirways as client.
	for {
		select {
		case out := <-chans.CheckAvailabilityAndPrice1FromAggregatorToBrutishAirways_string_int:
			g := AggregatorBrutishAirwaysGob{Name: "CheckAvailabilityAndPrice1_from_Aggregator_to_BrutishAirways_string_int", CheckAvailabilityAndPrice1_from_Aggregator_to_BrutishAirways_string_int: out}
			err := trans.Encoder.Encode(g)
			if err != nil {
				log.Fatal(err)
			}
		case <- chans.doneCommunicatingWithBrutishAirways:
			return
		default:
			// Keep looping
		}
	}
}

func ReceiveFromClient(chans *Channels, trans *Transmitter) {
	var in AggregatorClientGob
	for {
		err := trans.Decoder.Decode(&in)
		if err != nil {
			log.Fatal(err)
		}
		if in.Name == "GreetAggregator_from_Client_to_Aggregator_Empty" {
			chans.GreetAggregatorFromClientToAggregator_Empty <- in.GreetAggregator_from_Client_to_Aggregator_Empty
		} else if in.Name == "RequestItinerary_from_Client_to_Aggregator_string_int" {
			chans.RequestItineraryFromClientToAggregator_string_int <- in.RequestItinerary_from_Client_to_Aggregator_string_int
		} else if in.Name == "TryAgain_from_Client_to_Aggregator_Empty" {
			chans.TryAgainFromClientToAggregator_Empty <- in.TryAgain_from_Client_to_Aggregator_Empty
		} else if in.Name == "RejectAndLeave_from_Client_to_Aggregator_Empty" {
			chans.RejectAndLeaveFromClientToAggregator_Empty <- in.RejectAndLeave_from_Client_to_Aggregator_Empty
		} else if in.Name == "Accept_from_Client_to_Aggregator_Empty" {
			chans.AcceptFromClientToAggregator_Empty <- in.Accept_from_Client_to_Aggregator_Empty
		} else if in.Name == "ProvidePaymentInto_from_Client_to_Aggregator_string" {
			chans.ProvidePaymentIntoFromClientToAggregator_string <- in.ProvidePaymentInto_from_Client_to_Aggregator_string
		} else {
			log.Fatal("ReceiveFromServer() received unknown gob: ", in)
		}
	}
}

func ReceiveFromQueasyJet(chans *Channels, trans *Transmitter) {
	var in AggregatorQueasyJetGob
	for {
		err := trans.Decoder.Decode(&in)
		if err != nil {
			log.Fatal(err)
		}
		if in.Name == "ConfirmAvailabilityAndPrice2_from_QueasyJet_to_Aggregator_bool_int" {
			chans.ConfirmAvailabilityAndPrice2FromQueasyJetToAggregator_bool_int <- in.ConfirmAvailabilityAndPrice2_from_QueasyJet_to_Aggregator_bool_int
		} else {
			log.Fatal("ReceiveFromServer() received unknown gob: ", in)
		}
	}
}

func ReceiveFromBrutishAirways(chans *Channels, trans *Transmitter) {
	var in AggregatorBrutishAirwaysGob
	for {
		err := trans.Decoder.Decode(&in)
		if err != nil {
			log.Fatal(err)
		}
		if in.Name == "ConfirmAvailabilityAndPrice1_from_BrutishAirways_to_Aggregator_bool_int" {
			chans.ConfirmAvailabilityAndPrice1FromBrutishAirwaysToAggregator_bool_int <- in.ConfirmAvailabilityAndPrice1_from_BrutishAirways_to_Aggregator_bool_int
		} else {
			log.Fatal("ReceiveFromServer() received unknown gob: ", in)
		}
	}
}

func ConnectToClientAsClient(chans *Channels, conType string, serverAddress string) {
	trans := NewTransmitter(conType, serverAddress)
	go SendToClient(chans, trans)
	go ReceiveFromClient(chans, trans)
}

func ConnectToQueasyJetAsClient(chans *Channels, conType string, serverAddress string) {
	trans := NewTransmitter(conType, serverAddress)
	go SendToQueasyJet(chans, trans)
	go ReceiveFromQueasyJet(chans, trans)
}

func ConnectToBrutishAirwaysAsClient(chans *Channels, conType string, serverAddress string) {
	trans := NewTransmitter(conType, serverAddress)
	go SendToBrutishAirways(chans, trans)
	go ReceiveFromBrutishAirways(chans, trans)
}

func CloseConnectAsClientWithClient(chans *Channels) {
	chans.doneCommunicatingWithClient<- true
}

func CloseConnectAsClientWithQueasyJet(chans *Channels) {
	chans.doneCommunicatingWithQueasyJet<- true
}

func CloseConnectAsClientWithBrutishAirways(chans *Channels) {
	chans.doneCommunicatingWithBrutishAirways<- true
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
	if identifier.Id == "Client" {
		HandleClientConnectionAsServer(trans, chans)
	} else if identifier.Id == "QueasyJet" {
		HandleQueasyJetConnectionAsServer(trans, chans)
	} else if identifier.Id == "BrutishAirways" {
		HandleBrutishAirwaysConnectionAsServer(trans, chans)
	} else {
		log.Fatal("HandleConnection received unknown identifier: ", identifier)
	}
}

func HandleClientConnectionAsServer(trans *Transmitter, chans *Channels) {
	go SendToClient(chans, trans)
	go ReceiveFromClient(chans, trans)
}

func HandleQueasyJetConnectionAsServer(trans *Transmitter, chans *Channels) {
	go SendToQueasyJet(chans, trans)
	go ReceiveFromQueasyJet(chans, trans)
}

func HandleBrutishAirwaysConnectionAsServer(trans *Transmitter, chans *Channels) {
	go SendToBrutishAirways(chans, trans)
	go ReceiveFromBrutishAirways(chans, trans)
}

func SetupNetworkConnections(chans *Channels, connType string, address string,port string) {
	//go AcceptConnections(connType, port, chans) // Uncomment to accept connections as Server
	ConnectToClientAsClient(chans, connType, address + port) // Comment out to stop connecting as client
	ConnectToQueasyJetAsClient(chans, connType, address + port) // Comment out to stop connecting as client
	ConnectToBrutishAirwaysAsClient(chans, connType, address + port) // Comment out to stop connecting as client
}

func CloseNetworkConnections(chans *Channels) {
	chans.doneCommunicatingWithClient <- true
	<-chans.doneCommunicatingWithClient
	chans.doneCommunicatingWithQueasyJet <- true
	<-chans.doneCommunicatingWithQueasyJet
	chans.doneCommunicatingWithBrutishAirways <- true
	<-chans.doneCommunicatingWithBrutishAirways
}

