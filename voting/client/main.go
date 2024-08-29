package main

import (
	"fmt"
	"net/rpc"
	"os"
	"strconv"
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
	Vehicles []Vehicle
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
		fmt.Printf("Received vehicles: %v\n", reply.Vehicles) // Print received vehicles
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <VehicleNumber>")
		return
	}

	vehicleNumber := os.Args[1]
	vehicleAddress, err := strconv.Atoi(vehicleNumber)
	if err != nil || vehicleAddress < 0 || vehicleAddress > 2 {
		fmt.Println("Invalid vehicle number. It should be 0, 1, or 2.")
		return
	}

	// Define only the current vehicle
	vehicle := Vehicle{
		Number:    vehicleNumber,
		Direction: map[int]string{0: "N", 1: "S", 2: "E"}[vehicleAddress],
		Address:   vehicleAddress,
		Votes:     0,
	}

	// Request other vehicles from the server
	targetAddress := fmt.Sprintf("localhost:%d", 8000+vehicleAddress)
	client, err := rpc.Dial("tcp", targetAddress)
	if err != nil {
		fmt.Printf("Dialing error for vehicle at %s: %v\n", targetAddress, err)
		return
	}
	defer client.Close()

	var reply Response
	err = client.Call("VehicleRPC.GetVehicles", struct{}{}, &reply)
	if err != nil {
		fmt.Printf("RPC call error to vehicle at %s: %v\n", targetAddress, err)
	}
	fmt.Printf("Received vehicles from server: %v\n", reply.Vehicles)

	var wg sync.WaitGroup
	for i := 0; i < len(reply.Vehicles); i++ {
		wg.Add(1)
		go sendRequest(vehicle, fmt.Sprintf("localhost:%d", 8000+i), len(reply.Vehicles), &wg)
	}

	wg.Wait()

	fmt.Println("All requests have been sent.")
}
