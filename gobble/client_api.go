package main

import (
	"errors"
)

//Channels

type Channels struct {
	GreetAggregatorFromClientToAggregator_Empty chan GreetAggregator_from_Client_to_Aggregator_Empty
	GreetClientFromAggregatorToClient_Empty chan GreetClient_from_Aggregator_to_Client_Empty
	RequestItineraryFromClientToAggregator_string_int chan RequestItinerary_from_Client_to_Aggregator_string_int
	ProvideFlightInformationFromAggregatorToClient_bool_string chan ProvideFlightInformation_from_Aggregator_to_Client_bool_string
	TryAgainFromClientToAggregator_Empty chan TryAgain_from_Client_to_Aggregator_Empty
	RejectAndLeaveFromClientToAggregator_Empty chan RejectAndLeave_from_Client_to_Aggregator_Empty
	AcceptFromClientToAggregator_Empty chan Accept_from_Client_to_Aggregator_Empty
	RequestPaymentInfoFromAggregatorToClient_Empty chan RequestPaymentInfo_from_Aggregator_to_Client_Empty
	ProvidePaymentIntoFromClientToAggregator_string chan ProvidePaymentInto_from_Client_to_Aggregator_string
	ConfirmPaymentFromAggregatorToClient_bool chan ConfirmPayment_from_Aggregator_to_Client_bool
	doneCommunicatingWithAggregator chan bool
}

func NewChannels() *Channels {
	c := new(Channels)
	c.GreetAggregatorFromClientToAggregator_Empty = make(chan GreetAggregator_from_Client_to_Aggregator_Empty)
	c.GreetClientFromAggregatorToClient_Empty = make(chan GreetClient_from_Aggregator_to_Client_Empty)
	c.RequestItineraryFromClientToAggregator_string_int = make(chan RequestItinerary_from_Client_to_Aggregator_string_int)
	c.ProvideFlightInformationFromAggregatorToClient_bool_string = make(chan ProvideFlightInformation_from_Aggregator_to_Client_bool_string)
	c.TryAgainFromClientToAggregator_Empty = make(chan TryAgain_from_Client_to_Aggregator_Empty)
	c.RejectAndLeaveFromClientToAggregator_Empty = make(chan RejectAndLeave_from_Client_to_Aggregator_Empty)
	c.AcceptFromClientToAggregator_Empty = make(chan Accept_from_Client_to_Aggregator_Empty)
	c.RequestPaymentInfoFromAggregatorToClient_Empty = make(chan RequestPaymentInfo_from_Aggregator_to_Client_Empty)
	c.ProvidePaymentIntoFromClientToAggregator_string = make(chan ProvidePaymentInto_from_Client_to_Aggregator_string)
	c.ConfirmPaymentFromAggregatorToClient_bool = make(chan ConfirmPayment_from_Aggregator_to_Client_bool)
	c.doneCommunicatingWithAggregator = make(chan bool)
	return c
}

//Structs

type GreetAggregator_from_Client_to_Aggregator_Empty struct {
	Param1  struct{}
}

type GreetClient_from_Aggregator_to_Client_Empty struct {
	Param1  struct{}
}

type RequestItinerary_from_Client_to_Aggregator_string_int struct {
	Param1 string
	Param2 int
}

type ProvideFlightInformation_from_Aggregator_to_Client_bool_string struct {
	Param1 bool
	Param2 string
}

type TryAgain_from_Client_to_Aggregator_Empty struct {
	Param1  struct{}
}

type RejectAndLeave_from_Client_to_Aggregator_Empty struct {
	Param1  struct{}
}

type Accept_from_Client_to_Aggregator_Empty struct {
	Param1  struct{}
}

type RequestPaymentInfo_from_Aggregator_to_Client_Empty struct {
	Param1  struct{}
}

type ProvidePaymentInto_from_Client_to_Aggregator_string struct {
	Param1 string
}

type ConfirmPayment_from_Aggregator_to_Client_bool struct {
	Param1 bool
}

type Client1 struct {
	Channels *Channels
	Used bool
}

type Client2 struct {
	Channels *Channels
	Used bool
}

type Client_rec1_1 struct {
	Channels *Channels
	Used bool
}

type Client_rec1_2 struct {
	Channels *Channels
	Used bool
}

type Client_rec1_choice1 struct {
	Channels *Channels
	Used bool
}

type Client_rec1_choice1_C2 struct {
	Channels *Channels
	Used bool
}

type Client_rec1_choice1_C3 struct {
	Channels *Channels
	Used bool
}

type Client_rec1_choice1_C4 struct {
	Channels *Channels
	Used bool
}

//Methods

func (self *Client1) Send_GreetAggregator() (*Client2, error) {
	defer func() { self.Used = true }()
	sendVal := GreetAggregator_from_Client_to_Aggregator_Empty{}
	retVal := &Client2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client1; method: Send()")
	}
	self.Channels.GreetAggregatorFromClientToAggregator_Empty <- sendVal
	return retVal, nil
}

func (self *Client2) Receive_GreetClient() (*Client_rec1_1, error) {
	defer func() { self.Used = true }()
	retVal := &Client_rec1_1{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client2; method: Receive()")
	}
	<- self.Channels.GreetClientFromAggregatorToClient_Empty
	return retVal, nil
}

func (self *Client_rec1_1) Send_RequestItinerary_string_int(param1 string, param2 int) (*Client_rec1_2, error) {
	defer func() { self.Used = true }()
	sendVal := RequestItinerary_from_Client_to_Aggregator_string_int{Param1: param1, Param2: param2}
	retVal := &Client_rec1_2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_1; method: Send_string_int()")
	}
	self.Channels.RequestItineraryFromClientToAggregator_string_int <- sendVal
	return retVal, nil
}

func (self *Client_rec1_2) Receive_ProvideFlightInformation_bool_string() (*Client_rec1_choice1, bool, string, error) {
	defer func() { self.Used = true }()
	var in_1_bool bool
	var in_2_string string
	retVal := &Client_rec1_choice1{Channels: self.Channels}
	if self.Used {
		return retVal, in_1_bool, in_2_string, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_2; method: Receive_bool_string()")
	}
	in := <- self.Channels.ProvideFlightInformationFromAggregatorToClient_bool_string
	in_1_bool = in.Param1
	in_2_string = in.Param2
	return retVal, in_1_bool, in_2_string, nil
}

func (self *Client_rec1_choice1) Send_TryAgain() (*Client_rec1_1, error) {
	defer func() { self.Used = true }()
	sendVal := TryAgain_from_Client_to_Aggregator_Empty{}
	retVal := &Client_rec1_1{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_choice1; method: Send()")
	}
	self.Channels.TryAgainFromClientToAggregator_Empty <- sendVal
	return retVal, nil
}

func (self *Client_rec1_choice1) Send_RejectAndLeave() (error) {
	defer func() { self.Used = true }()
	sendVal := RejectAndLeave_from_Client_to_Aggregator_Empty{}
	if self.Used {
		return errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_choice1; method: Send()")
	}
	self.Channels.RejectAndLeaveFromClientToAggregator_Empty <- sendVal
	return nil
}

func (self *Client_rec1_choice1) Send_Accept() (*Client_rec1_choice1_C2, error) {
	defer func() { self.Used = true }()
	sendVal := Accept_from_Client_to_Aggregator_Empty{}
	retVal := &Client_rec1_choice1_C2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_choice1; method: Send()")
	}
	self.Channels.AcceptFromClientToAggregator_Empty <- sendVal
	return retVal, nil
}

func (self *Client_rec1_choice1_C2) Receive_RequestPaymentInfo() (*Client_rec1_choice1_C3, error) {
	defer func() { self.Used = true }()
	retVal := &Client_rec1_choice1_C3{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_choice1_C2; method: Receive()")
	}
	<- self.Channels.RequestPaymentInfoFromAggregatorToClient_Empty
	return retVal, nil
}

func (self *Client_rec1_choice1_C3) Send_ProvidePaymentInto_string(param1 string) (*Client_rec1_choice1_C4, error) {
	defer func() { self.Used = true }()
	sendVal := ProvidePaymentInto_from_Client_to_Aggregator_string{Param1: param1}
	retVal := &Client_rec1_choice1_C4{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_choice1_C3; method: Send_string()")
	}
	self.Channels.ProvidePaymentIntoFromClientToAggregator_string <- sendVal
	return retVal, nil
}

func (self *Client_rec1_choice1_C4) Receive_ConfirmPayment_bool() (bool, error) {
	defer func() { self.Used = true }()
	var in_1_bool bool
	if self.Used {
		return in_1_bool, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_choice1_C4; method: Receive_bool()")
	}
	in := <- self.Channels.ConfirmPaymentFromAggregatorToClient_bool
	in_1_bool = in.Param1
	return in_1_bool, nil
}

