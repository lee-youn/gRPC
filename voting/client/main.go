package main

import (
	"fmt"
	"net/rpc"
	"sync"
)

type Vehicle struct {
	Number    string
	Direction string
	Address   int
	Votes     int
}

type Request struct {
	Vehicle       Vehicle
	TotalVehicles int
}

type Response struct {
	Vehicle Vehicle
}

func sendRequest(vehicle Vehicle, targetAddress string, totalVehicles int, wg *sync.WaitGroup) {
	defer wg.Done()

	client, err := rpc.Dial("tcp", targetAddress)
	if err != nil {
		fmt.Printf("Dialing error for vehicle at %s: %v\n", targetAddress, err)
		return
	}
	defer client.Close()

	req := Request{
		Vehicle:       vehicle,
		TotalVehicles: totalVehicles,
	}
	var reply Response

	err = client.Call("VehicleRPC.SendRequest", req, &reply)
	if err != nil {
		fmt.Printf("RPC call error to vehicle at %s: %v\n", targetAddress, err)
	} else {
		fmt.Printf("Successfully sent request to vehicle at %s.\n", targetAddress)
	}
}

func main() {
	vehicles := []Vehicle{
		{Number: "0", Direction: "N", Address: 0, Votes: 0},
		{Number: "1", Direction: "S", Address: 1, Votes: 0},
		{Number: "2", Direction: "E", Address: 2, Votes: 0},
	}

	totalVehicles := len(vehicles)

	var wg sync.WaitGroup
	for _, targetVehicle := range vehicles {
		wg.Add(1)
		go sendRequest(targetVehicle, fmt.Sprintf("localhost:%d", 8000+targetVehicle.Address), totalVehicles, &wg)
	}
	wg.Wait()
	fmt.Println("All requests have been sent.")
}
