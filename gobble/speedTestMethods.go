package main

import (
	"errors"
)

func (self *Server_rec1_1) Receive_firstStep_int() (*Server_rec1_2, int, error) {
	defer func() { self.Used = true }()
	var in_1_int int
	retVal := &Server_rec1_2{Channels: self.Channels}
	if self.Used {
		return retVal, in_1_int, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Server_rec1_1; method: Receive_int()")
	}
	in := <-self.Channels.firstStepFromClientToServer_int
	in_1_int = in.Param1
	return retVal, in_1_int, nil
}

func (self *Server_rec1_2) Send_secondStep_int(param1 int) (*Server_rec1_1, error) {
	defer func() { self.Used = true }()
	sendVal := SecondStep_from_Server_to_Client_int{Param1: param1}
	retVal := &Server_rec1_1{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Server_rec1_2; method: Send_int()")
	}
	self.Channels.secondStepFromServerToClient_int <- sendVal
	return retVal, nil
}

func (self *Client_rec1_1) Send_firstStep_int(param1 int) (*Client_rec1_2, error) {
	defer func() { self.Used = true }()
	sendVal := FirstStep_from_Client_to_Server_int{Param1: param1}
	retVal := &Client_rec1_2{Channels: self.Channels}
	if self.Used {
		return retVal, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_1; method: Send_int()")
	}
	self.Channels.firstStepFromClientToServer_int <- sendVal
	return retVal, nil
}

func (self *Client_rec1_2) Receive_secondStep_int() (*Client_rec1_1, int, error) {
	defer func() { self.Used = true }()
	var in_1_int int
	retVal := &Client_rec1_1{Channels: self.Channels}
	if self.Used {
		return retVal, in_1_int, errors.New("Dynamic session type checking error: attempted repeat method call to one or more methods of a given session type struct. Struct: Client_rec1_2; method: Receive_int()")
	}
	in := <-self.Channels.secondStepFromServerToClient_int
	in_1_int = in.Param1
	return retVal, in_1_int, nil
}
