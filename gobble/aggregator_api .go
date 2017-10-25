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
	CheckAvailabilityAndPrice2FromAggregatorToQueasyJet_string_int chan CheckAvailabilityAndPrice2_from_Aggregator_to_QueasyJet_string_int
	ConfirmAvailabilityAndPrice2FromQueasyJetToAggregator_bool_int chan ConfirmAvailabilityAndPrice2_from_QueasyJet_to_Aggregator_bool_int
	CheckAvailabilityAndPrice1FromAggregatorToBrutishAirways_string_int chan CheckAvailabilityAndPrice1_from_Aggregator_to_BrutishAirways_string_int
	ConfirmAvailabilityAndPrice1FromBrutishAirwaysToAggregator_bool_int chan ConfirmAvailabilityAndPrice1_from_BrutishAirways_to_Aggregator_bool_int
	TryAgainFromClientToAggregator_Empty chan TryAgain_from_Client_to_Aggregator_Empty
	RejectAndLeaveFromClientToAggregator_Empty chan RejectAndLeave_from_Client_to_Aggregator_Empty
	AcceptFromClientToAggregator_Empty chan Accept_from_Client_to_Aggregator_Empty
	RequestPaymentInfoFromAggregatorToClient_Empty chan RequestPaymentInfo_from_Aggregator_to_Client_Empty
	ProvidePaymentIntoFromClientToAggregator_string chan ProvidePaymentInto_from_Client_to_Aggregator_string
	ConfirmPaymentFromAggregatorToClient_bool chan ConfirmPayment_from_Aggregator_to_Client_bool
	done_rec1_par1_A chan bool
	done_rec1_par1_B chan bool
	TryAgainFromAggregatorToAggregator_Empty chan TryAgain_from_Client_to_Aggregator_Empty
	RejectAndLeaveFromAggregatorToAggregator_Empty chan RejectAndLeave_from_Client_to_Aggregator_Empty
	AcceptFromAggregatorToAggregator_Empty chan Accept_from_Client_to_Aggregator_Empty
	doneCommunicatingWithClient chan bool
	doneCommunicatingWithQueasyJet chan bool
	doneCommunicatingWithBrutishAirways chan bool
}

func NewChannels() *Channels {
	c := new(Channels)
	c.GreetAggregatorFromClientToAggregator_Empty = make(chan GreetAggregator_from_Client_to_Aggregator_Empty)
	c.GreetClientFromAggregatorToClient_Empty = make(chan GreetClient_from_Aggregator_to_Client_Empty)
	c.RequestItineraryFromClientToAggregator_string_int = make(chan RequestItinerary_from_Client_to_Aggregator_string_int)
	c.ProvideFlightInformationFromAggregatorToClient_bool_string = make(chan ProvideFlightInformation_from_Aggregator_to_Client_bool_string)
	c.CheckAvailabilityAndPrice2FromAggregatorToQueasyJet_string_int = make(chan CheckAvailabilityAndPrice2_from_Aggregator_to_QueasyJet_string_int)
	c.ConfirmAvailabilityAndPrice2FromQueasyJetToAggregator_bool_int = make(chan ConfirmAvailabilityAndPrice2_from_QueasyJet_to_Aggregator_bool_int)
	c.CheckAvailabilityAndPrice1FromAggregatorToBrutishAirways_string_int = make(chan CheckAvailabilityAndPrice1_from_Aggregator_to_BrutishAirways_string_int)
	c.ConfirmAvailabilityAndPrice1FromBrutishAirwaysToAggregator_bool_int = make(chan ConfirmAvailabilityAndPrice1_from_BrutishAirways_to_Aggregator_bool_int)
	c.TryAgainFromClientToAggregator_Empty = make(chan TryAgain_from_Client_to_Aggregator_Empty)
	c.RejectAndLeaveFromClientToAggregator_Empty = make(chan RejectAndLeave_from_Client_to_Aggregator_Empty)
	c.AcceptFromClientToAggregator_Empty = make(chan Accept_from_Client_to_Aggregator_Empty)
	c.RequestPaymentInfoFromAggregatorToClient_Empty = make(chan RequestPaymentInfo_from_Aggregator_to_Client_Empty)
	c.ProvidePaymentIntoFromClientToAggregator_string = make(chan ProvidePaymentInto_from_Client_to_Aggregator_string)
	c.ConfirmPaymentFromAggregatorToClient_bool = make(chan ConfirmPayment_from_Aggregator_to_Client_bool)
	c.doneCommunicatingWithClient = make(chan bool)
	c.doneCommunicatingWithQueasyJet = make(chan bool)
	c.doneCommunicatingWithBrutishAirways = make(chan bool)
	c.done_rec1_par1_A = make(chan bool, 1)
	c.done_rec1_par1_B = make(chan bool, 1)
	c.TryAgainFromAggregatorToAggregator_Empty = make(chan TryAgain_from_Client_to_Aggregator_Empty, 1)
	c.RejectAndLeaveFromAggregatorToAggregator_Empty = make(chan RejectAndLeave_from_Client_to_Aggregator_Empty, 1)
	c.AcceptFromAggregatorToAggregator_Empty = make(chan Accept_from_Client_to_Aggregator_Empty, 1)
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

type CheckAvailabilityAndPrice2_from_Aggregator_to_QueasyJet_string_int struct {
	Param1 string
	Param2 int
}

type ConfirmAvailabilityAndPrice2_from_QueasyJet_to_Aggregator_bool_int struct {
	Param1 bool
	Param2 int
}

type CheckAvailabilityAndPrice1_from_Aggregator_to_BrutishAirways_string_int struct {
	Param1 string
	Param2 int
}

type ConfirmAvailabilityAndPrice1_from_BrutishAirways_to_Aggregator_bool_int struct {
	Param1 bool
	Param2 int
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

type Aggregator1 struct {
	Channels *Channels
	Used bool
}

type Aggregator2 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_1 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_2 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_par1_A1 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_par1_A2 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_par1_B1 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_par1_B2 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_choice1 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_choice1_C2 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_choice1_C3 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_choice1_C4 struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_par1_start struct {
	Channels *Channels
	Used bool
}

type Aggregator_rec1_par1_end struct {
	Channels *Channels
	Used bool
}

//Methods

func (self *Aggregator1) Receive_GreetAggregator() (*Aggregator2, error) {
	defer func() { self.Used = true }()
	retVal := &Aggregator2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator1; method: Receive()")
	}
	<- self.Channels.GreetAggregatorFromClientToAggregator_Empty
	return retVal, nil
}

func (self *Aggregator2) Send_GreetClient() (*Aggregator_rec1_1, error) {
	defer func() { self.Used = true }()
	sendVal := GreetClient_from_Aggregator_to_Client_Empty{}
	retVal := &Aggregator_rec1_1{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator2; method: Send()")
	}
	self.Channels.GreetClientFromAggregatorToClient_Empty <- sendVal
	return retVal, nil
}

func (self *Aggregator_rec1_1) Receive_RequestItinerary_string_int() (*Aggregator_rec1_par1_start, string, int, error) {
	defer func() { self.Used = true }()
	var in_1_string string
	var in_2_int int
	retVal := &Aggregator_rec1_par1_start{Channels: self.Channels}
	if self.Used {
		return retVal, in_1_string, in_2_int, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_1; method: Receive_string_int()")
	}
	in := <- self.Channels.RequestItineraryFromClientToAggregator_string_int
	in_1_string = in.Param1
	in_2_int = in.Param2
	return retVal, in_1_string, in_2_int, nil
}

func (self *Aggregator_rec1_2) Send_ProvideFlightInformation_bool_string(param1 bool, param2 string) (*Aggregator_rec1_choice1, error) {
	defer func() { self.Used = true }()
	sendVal := ProvideFlightInformation_from_Aggregator_to_Client_bool_string{Param1: param1, Param2: param2}
	retVal := &Aggregator_rec1_choice1{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_2; method: Send_bool_string()")
	}
	self.Channels.ProvideFlightInformationFromAggregatorToClient_bool_string <- sendVal
	return retVal, nil
}

func (self *Aggregator_rec1_par1_A1) Send_CheckAvailabilityAndPrice2_string_int(param1 string, param2 int) (*Aggregator_rec1_par1_A2, error) {
	defer func() { self.Used = true }()
	sendVal := CheckAvailabilityAndPrice2_from_Aggregator_to_QueasyJet_string_int{Param1: param1, Param2: param2}
	retVal := &Aggregator_rec1_par1_A2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_par1_A1; method: Send_string_int()")
	}
	self.Channels.CheckAvailabilityAndPrice2FromAggregatorToQueasyJet_string_int <- sendVal
	return retVal, nil
}

func (self *Aggregator_rec1_par1_A2) Receive_ConfirmAvailabilityAndPrice2_bool_int() (*Aggregator_rec1_par1_start, bool, int, error) {
	defer func() { self.Used = true }()
	var in_1_bool bool
	var in_2_int int
	retVal := &Aggregator_rec1_par1_start{Channels: self.Channels}
	if self.Used {
		return retVal, in_1_bool, in_2_int, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_par1_A2; method: Receive_bool_int()")
	}
	in := <- self.Channels.ConfirmAvailabilityAndPrice2FromQueasyJetToAggregator_bool_int
	in_1_bool = in.Param1
	in_2_int = in.Param2
	self.Channels.done_rec1_par1_A <- true
	return retVal, in_1_bool, in_2_int, nil
}

func (self *Aggregator_rec1_par1_B1) Send_CheckAvailabilityAndPrice1_string_int(param1 string, param2 int) (*Aggregator_rec1_par1_B2, error) {
	defer func() { self.Used = true }()
	sendVal := CheckAvailabilityAndPrice1_from_Aggregator_to_BrutishAirways_string_int{Param1: param1, Param2: param2}
	retVal := &Aggregator_rec1_par1_B2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_par1_B1; method: Send_string_int()")
	}
	self.Channels.CheckAvailabilityAndPrice1FromAggregatorToBrutishAirways_string_int <- sendVal
	return retVal, nil
}

func (self *Aggregator_rec1_par1_B2) Receive_ConfirmAvailabilityAndPrice1_bool_int() (*Aggregator_rec1_par1_start, bool, int, error) {
	defer func() { self.Used = true }()
	var in_1_bool bool
	var in_2_int int
	retVal := &Aggregator_rec1_par1_start{Channels: self.Channels}
	if self.Used {
		return retVal, in_1_bool, in_2_int, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_par1_B2; method: Receive_bool_int()")
	}
	in := <- self.Channels.ConfirmAvailabilityAndPrice1FromBrutishAirwaysToAggregator_bool_int
	in_1_bool = in.Param1
	in_2_int = in.Param2
	self.Channels.done_rec1_par1_B <- true
	return retVal, in_1_bool, in_2_int, nil
}

func (self *Aggregator_rec1_choice1) Receive_TryAgain() (*Aggregator_rec1_1, error) {
	defer func() { self.Used = true }()
	retVal := &Aggregator_rec1_1{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_choice1; method: Receive()")
	}
	<- self.Channels.TryAgainFromAggregatorToAggregator_Empty
	return retVal, nil
}

func (self *Aggregator_rec1_choice1) Receive_RejectAndLeave() (error) {
	defer func() { self.Used = true }()
	if self.Used {
		return errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_choice1; method: Receive()")
	}
	<- self.Channels.RejectAndLeaveFromAggregatorToAggregator_Empty
	return nil
}

func (self *Aggregator_rec1_choice1) Receive_Accept() (*Aggregator_rec1_choice1_C2, error) {
	defer func() { self.Used = true }()
	retVal := &Aggregator_rec1_choice1_C2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_choice1; method: Receive()")
	}
	<- self.Channels.AcceptFromAggregatorToAggregator_Empty
	return retVal, nil
}

func (self *Aggregator_rec1_choice1_C2) Send_RequestPaymentInfo() (*Aggregator_rec1_choice1_C3, error) {
	defer func() { self.Used = true }()
	sendVal := RequestPaymentInfo_from_Aggregator_to_Client_Empty{}
	retVal := &Aggregator_rec1_choice1_C3{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_choice1_C2; method: Send()")
	}
	self.Channels.RequestPaymentInfoFromAggregatorToClient_Empty <- sendVal
	return retVal, nil
}

func (self *Aggregator_rec1_choice1_C3) Receive_ProvidePaymentInto_string() (*Aggregator_rec1_choice1_C4, string, error) {
	defer func() { self.Used = true }()
	var in_1_string string
	retVal := &Aggregator_rec1_choice1_C4{Channels: self.Channels}
	if self.Used {
		return retVal, in_1_string, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_choice1_C3; method: Receive_string()")
	}
	in := <- self.Channels.ProvidePaymentIntoFromClientToAggregator_string
	in_1_string = in.Param1
	return retVal, in_1_string, nil
}

func (self *Aggregator_rec1_choice1_C4) Send_ConfirmPayment_bool(param1 bool) (error) {
	defer func() { self.Used = true }()
	sendVal := ConfirmPayment_from_Aggregator_to_Client_bool{Param1: param1}
	if self.Used {
		return errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_choice1_C4; method: Send_bool()")
	}
	self.Channels.ConfirmPaymentFromAggregatorToClient_bool <- sendVal
	return nil
}

func (self *Aggregator_rec1_par1_start) StartPar() (*Aggregator_rec1_par1_end, *Aggregator_rec1_par1_A1, *Aggregator_rec1_par1_B1, error) {
	defer func() { self.Used = true }()
	Aggregator_rec1_par1_end := &Aggregator_rec1_par1_end{Channels: self.Channels}
	Aggregator_rec1_par1_A1 := &Aggregator_rec1_par1_A1{Channels: self.Channels}
	Aggregator_rec1_par1_B1 := &Aggregator_rec1_par1_B1{Channels: self.Channels}
	if self.Used {
		return Aggregator_rec1_par1_end, Aggregator_rec1_par1_A1, Aggregator_rec1_par1_B1, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_par1_start; method: StartPar()")
	}
	return Aggregator_rec1_par1_end, Aggregator_rec1_par1_A1, Aggregator_rec1_par1_B1, nil
}

func (self *Aggregator_rec1_par1_end) EndPar() (*Aggregator_rec1_2, error) {
	defer func() { self.Used = true }()
	retVal := &Aggregator_rec1_2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Aggregator_rec1_par1_start; method: EndPar()")
	}
	par_A_done := false
	select{
	case <- self.Channels.done_rec1_par1_A:
		par_A_done = true
	default:
	}
	par_B_done := false
	select{
	case <- self.Channels.done_rec1_par1_B:
		par_B_done = true
	default:
	}
	if !par_A_done {
		return retVal, errors.New("Dynamic session type checking error: attempted call to Seller_rec1_par1_end.EndPar() prior to completion of parallel process Seller_rec1_par1_A")
	}
	if !par_B_done {
		return retVal, errors.New("Dynamic session type checking error: attempted call to Seller_rec1_par1_end.EndPar() prior to completion of parallel process Seller_rec1_par1_B")
	}
	return retVal, nil
}

