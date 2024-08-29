package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
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

type ConcurrentVehicle struct {
	Vehicle Vehicle
	Next    *ConcurrentVehicle
}

type VehicleRPC struct {
	Address               int
	Vehicles              map[int]*Vehicle
	TotalVehicles         int
	VoteCount             int
	ConcurrentVehicleList *ConcurrentVehicle
	mu                    sync.Mutex
}

func (v *VehicleRPC) GetVehicles(req struct{}, reply *Response) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Convert the map of vehicles to a slice
	vehicles := make([]Vehicle, 0, len(v.Vehicles))
	for _, vehicle := range v.Vehicles {
		if vehicle != nil { // Ensure we don't dereference a nil pointer
			vehicles = append(vehicles, *vehicle) // Dereference the pointer to get the Vehicle value
		}
	}
	reply.Vehicles = vehicles

	return nil
}

func (v *VehicleRPC) SendRequest(req Request, reply *Response) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	fmt.Printf("Received request for vehicle: %v\n", req.Vehicle)

	// Increment votes for the vehicle if it exists
	if vehicle, exists := v.Vehicles[req.Vehicle.Address]; exists {
		vehicle.Votes++
		fmt.Printf("Updated vehicle %d votes to %d\n", vehicle.Address, vehicle.Votes)
	} else {
		fmt.Printf("Vehicle %d not found in the system.\n", req.Vehicle.Address)
	}

	// Prepare to send RPC requests to other vehicles
	var wg sync.WaitGroup
	var requestVehicles []Vehicle

	for i := 0; i < req.TotalVehicles; i++ {
		if i != req.Vehicle.Address {
			address := fmt.Sprintf("localhost:%d", 8000+i)
			fmt.Printf("Vehicle %d attempting to connect to vehicle %d at %s\n", req.Vehicle.Address, i, address)

			// Check if the vehicle exists before attempting to append
			if vehicle, exists := v.Vehicles[i]; exists {
				requestVehicles = append(requestVehicles, *vehicle) // Dereference the pointer to get Vehicle value
			}

			wg.Add(1)
			go func(addr string, vehicleAddr int) {
				defer wg.Done()

				client, err := rpc.Dial("tcp", addr)
				if err != nil {
					fmt.Printf("Failed to connect to vehicle %d at %s: %v\n", vehicleAddr, addr, err)
					return
				}
				defer client.Close()

				fmt.Printf("Attempting RPC call to vehicle %d at %s\n", vehicleAddr, addr)

				go func() {
					var response Response
					req := Request{
						Vehicle:       *v.Vehicles[req.Vehicle.Address],
						TotalVehicles: req.TotalVehicles,
					}

					err = client.Call("VehicleRPC.SendResponse", req, &response)
					if err != nil {
						fmt.Printf("RPC call error to vehicle %d: %v\n", vehicleAddr, err)
					} else {
						fmt.Printf("Successfully sent response to vehicle %d.\n", vehicleAddr)
					}
				}()
			}(address, i)
		}
	}

	wg.Wait()
	reply.Vehicles = requestVehicles

	return nil
}

func (v *VehicleRPC) SendResponse(req Request, reply *Response) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	fmt.Printf("Vehicle %d received response request: %v\n", v.Address, req.Vehicle)
	reply.Vehicles = []Vehicle{req.Vehicle}

	return nil
}

func (v *VehicleRPC) ReceiveResponse(req Request, reply *Response) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.VoteCount++
	fmt.Printf("Vehicle %d received response from vehicle: %v\n", v.Address, req.Vehicle)

	if v.ConcurrentVehicleList == nil {
		v.ConcurrentVehicleList = &ConcurrentVehicle{Vehicle: req.Vehicle}
	} else {
		current := v.ConcurrentVehicleList
		for current.Next != nil {
			if canPassTogether(current.Vehicle.Direction, req.Vehicle.Direction) {
				current = current.Next
			} else {
				newConcurrentVehicle := &ConcurrentVehicle{Vehicle: req.Vehicle, Next: current.Next}
				current.Next = newConcurrentVehicle
				break
			}
		}
	}

	if v.VoteCount >= v.TotalVehicles {
		fmt.Printf("Vehicle %d is now the leader.\n", v.Address)
		current := v.ConcurrentVehicleList
		for current != nil {
			fmt.Printf("Vehicle %d is passing.\n", current.Vehicle.Address)
			current = current.Next
		}
		fmt.Printf("Vehicle %d is passing.\n", v.Address)
	}

	return nil
}

func canPassTogether(direction1, direction2 string) bool {
	return direction1 == direction2
}

func startServer(address int, vehicles map[int]*Vehicle, totalVehicles int) {
	vehicleRPC := &VehicleRPC{
		Address:       address,
		Vehicles:      vehicles,
		TotalVehicles: totalVehicles,
		VoteCount:     0,
	}

	rpc.Register(vehicleRPC)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000+address))
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer listener.Close()

	fmt.Printf("Vehicle %d listening on port %d\n", address, 8000+address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func main() {
	totalVehicles := 3

	// Initialize vehicles for each server
	for i := 0; i < totalVehicles; i++ {
		vehicles := make(map[int]*Vehicle)
		for j := 0; j < totalVehicles; j++ {
			vehicles[j] = &Vehicle{Number: fmt.Sprintf("%d", j+1), Direction: "N", Address: j, Votes: 0}
		}
		go startServer(i, vehicles, totalVehicles)
	}

	time.Sleep(80 * time.Second) // Wait enough time for all servers to start and run
}
