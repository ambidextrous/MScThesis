package main

import (
	"sync"
	"log"
	"fmt"
)

func runAggregator_rec1_par1_A(Aggregator_rec1_par1_A1 *Aggregator_rec1_par1_A1, wg *sync.WaitGroup) (interface{}, error) {
	defer wg.Done()
	var sending_1_1_string string
	var sending_1_2_int int
	Aggregator_rec1_par1_A2, err1 := Aggregator_rec1_par1_A1.Send_CheckAvailabilityAndPrice2_string_int(sending_1_1_string, sending_1_2_int)
	if err1 != nil {
		log.Fatal(err1)
	}
	Aggregator_rec1_par1_start, received_2_1_bool, received_2_2_int, err2 := Aggregator_rec1_par1_A2.Receive_ConfirmAvailabilityAndPrice2_bool_int()
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("Aggregator_rec1_par1_A2 received value type: bool: ", received_2_1_bool)
	fmt.Println("Aggregator_rec1_par1_A2 received value type: int: ", received_2_2_int)
	return Aggregator_rec1_par1_start, nil
}

func runAggregator_rec1_par1_B(Aggregator_rec1_par1_B1 *Aggregator_rec1_par1_B1, wg *sync.WaitGroup) (interface{}, error) {
	defer wg.Done()
	var sending_1_1_string string
	var sending_1_2_int int
	Aggregator_rec1_par1_B2, err1 := Aggregator_rec1_par1_B1.Send_CheckAvailabilityAndPrice1_string_int(sending_1_1_string, sending_1_2_int)
	if err1 != nil {
		log.Fatal(err1)
	}
	Aggregator_rec1_par1_start, received_2_1_bool, received_2_2_int, err2 := Aggregator_rec1_par1_B2.Receive_ConfirmAvailabilityAndPrice1_bool_int()
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println("Aggregator_rec1_par1_B2 received value type: bool: ", received_2_1_bool)
	fmt.Println("Aggregator_rec1_par1_B2 received value type: int: ", received_2_2_int)
	return Aggregator_rec1_par1_start, nil
}

func runAggregator_rec1_choice1_A(Aggregator_rec1_choice1 *Aggregator_rec1_choice1) (interface{}, error) {
	Aggregator_rec1_1_new, err1 := Aggregator_rec1_choice1.Receive_TryAgain()
	if err1 != nil {
		log.Fatal(err1)
	}
	return Aggregator_rec1_1_new, nil
}

func runAggregator_rec1_choice1_B(Aggregator_rec1_choice1 *Aggregator_rec1_choice1) (interface{}, error) {
	err1 := Aggregator_rec1_choice1.Receive_RejectAndLeave()
	if err1 != nil {
		log.Fatal(err1)
	}
	return "", nil
}

func runAggregator_rec1_choice1_C(Aggregator_rec1_choice1 *Aggregator_rec1_choice1) (interface{}, error) {
	Aggregator_rec1_choice1_C2, err1 := Aggregator_rec1_choice1.Receive_Accept()
	if err1 != nil {
		log.Fatal(err1)
	}
	Aggregator_rec1_choice1_C3, err2 := Aggregator_rec1_choice1_C2.Send_RequestPaymentInfo()
	if err2 != nil {
		log.Fatal(err2)
	}
	Aggregator_rec1_choice1_C4, received_3_1_string, err3 := Aggregator_rec1_choice1_C3.Receive_ProvidePaymentInto_string()
	if err3 != nil {
		log.Fatal(err3)
	}
	fmt.Println("Aggregator_rec1_choice1_C3 received value type: string: ", received_3_1_string)
	var sending_4_1_bool bool
	err4 := Aggregator_rec1_choice1_C4.Send_ConfirmPayment_bool(sending_4_1_bool)
	if err4 != nil {
		log.Fatal(err4)
	}
	return "", nil
}

func runAggregator_rec1(Aggregator_rec1_1 *Aggregator_rec1_1) (interface{}, error) {
	Aggregator_rec1_par1_start, received_1_1_string, received_1_2_int, err1 := Aggregator_rec1_1.Receive_RequestItinerary_string_int()
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println("Aggregator_rec1_1 received value type: string: ", received_1_1_string)
	fmt.Println("Aggregator_rec1_1 received value type: int: ", received_1_2_int)
	Aggregator_rec1_2_candidate, err2 := runAggregator_rec1_par1(Aggregator_rec1_par1_start)
	if err2 != nil {
		log.Fatal(err2)
	}
	var Aggregator_rec1_2_step2 *Aggregator_rec1_2
	switch v := Aggregator_rec1_2_candidate.(type) {
	case *Aggregator_rec1_2:
		Aggregator_rec1_2_step2 = v
	default:
		log.Fatalf("Expected type Aggregator_rec1_2, received type %T", v)
	}
	var sending_3_1_bool bool
	var sending_3_2_string string
	Aggregator_rec1_choice1, err3 := Aggregator_rec1_2_step2.Send_ProvideFlightInformation_bool_string(sending_3_1_bool, sending_3_2_string)
	if err3 != nil {
		log.Fatal(err3)
	}
	retVal, err4 := Make_Aggregator_rec1_choice1_Choices(Aggregator_rec1_choice1)
	if err4 != nil {
		log.Fatal(err4)
	}
	return retVal, nil
}

func runAggregator(wg *sync.WaitGroup, Aggregator1 *Aggregator1) (interface{}, error) {
	defer wg.Done()
	Aggregator2, err1 := Aggregator1.Receive_GreetAggregator()
	if err1 != nil {
		log.Fatal(err1)
	}
	Aggregator_rec1_1, err2 := Aggregator2.Send_GreetClient()
	if err2 != nil {
		log.Fatal(err2)
	}
	retVal, err3 := LoopAggregator_rec1(Aggregator_rec1_1)
	if err3 != nil {
		log.Fatal(err3)
	}
	return retVal, nil
}

// Functions

func runAggregator_rec1_par1(Aggregator_rec1_par1_start *Aggregator_rec1_par1_start) (interface{}, error) {
	Aggregator_rec1_par1_end, Aggregator_rec1_par1_A1, Aggregator_rec1_par1_B1, err := Aggregator_rec1_par1_start.StartPar()
	if err != nil {
		log.Fatal(err)
	}
	var newWg sync.WaitGroup
	newWg.Add(1)
	go runAggregator_rec1_par1_A(Aggregator_rec1_par1_A1, &newWg)
	newWg.Add(1)
	go runAggregator_rec1_par1_B(Aggregator_rec1_par1_B1, &newWg)
	newWg.Wait()
		Aggregator_rec12, err_end := Aggregator_rec1_par1_end.EndPar()
	if err_end != nil {
		log.Fatal(err_end)
	}
	return Aggregator_rec12, nil
}

func Make_Aggregator_rec1_choice1_Choices(Aggregator_rec1_choice1 *Aggregator_rec1_choice1) (interface{}, error) {
	var retVal interface{}
	var err error
	select {
	case received_A := <- Aggregator_rec1_choice1.Channels.TryAgainFromClientToAggregator_Empty:
		Aggregator_rec1_choice1.Channels.TryAgainFromAggregatorToAggregator_Empty <- received_A
		retVal, err = runAggregator_rec1_choice1_A(Aggregator_rec1_choice1)
		if err != nil {
			log.Fatal(err)
		}
	case received_B := <- Aggregator_rec1_choice1.Channels.RejectAndLeaveFromClientToAggregator_Empty:
		Aggregator_rec1_choice1.Channels.RejectAndLeaveFromAggregatorToAggregator_Empty <- received_B
		retVal, err = runAggregator_rec1_choice1_B(Aggregator_rec1_choice1)
		if err != nil {
			log.Fatal(err)
		}
	case received_C := <- Aggregator_rec1_choice1.Channels.AcceptFromClientToAggregator_Empty:
		Aggregator_rec1_choice1.Channels.AcceptFromAggregatorToAggregator_Empty <- received_C
		retVal, err = runAggregator_rec1_choice1_C(Aggregator_rec1_choice1)
		if err != nil {
			log.Fatal(err)
		}
	}
	return retVal, nil
}

func StartAggregator(chans *Channels) (*Aggregator1) {
	start := &Aggregator1{Channels: chans}
	return start
}

func LoopAggregator_rec1(Aggregator_rec1_1_old *Aggregator_rec1_1) (interface{}, error) {
	var retVal interface{}
	var err error
	looping := true
	for looping {
		retVal, err = runAggregator_rec1(Aggregator_rec1_1_old)
		if err != nil {
			log.Fatal(err)
		}
		switch t := retVal.(type) {
		case *Aggregator_rec1_1:
			Aggregator_rec1_1_old = t
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
	startStruct := StartAggregator(chans)
	var newWg sync.WaitGroup
	newWg.Add(1)
	go runAggregator(&newWg, startStruct)
	newWg.Wait()
	CloseNetworkConnections(chans)
}

