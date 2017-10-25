package main

// Channels

type Channels struct {
	doneCommunicatingWithClient chan bool
	doneCommunicatingWithServer chan bool
	firstStepFromClientToServer_int chan FirstStep_from_Client_to_Server_int
	secondStepFromServerToClient_int chan SecondStep_from_Server_to_Client_int
}

func NewChannels() *Channels {
	c := new(Channels)
	c.doneCommunicatingWithClient = make(chan bool)
	c.doneCommunicatingWithServer = make(chan bool)
	c.firstStepFromClientToServer_int = make(chan FirstStep_from_Client_to_Server_int)
	c.secondStepFromServerToClient_int = make(chan SecondStep_from_Server_to_Client_int)
	return c
}

