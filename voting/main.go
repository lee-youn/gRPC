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
	Vehicle Vehicle
}

type ConcurrentVehicle struct {
	Vehicle Vehicle
	Next    *ConcurrentVehicle
}

type VehicleRPC struct {
	Address               int
	Vehicles              []Vehicle
	TotalVehicles         int
	VoteCount             int
	ConcurrentVehicleList *ConcurrentVehicle // 동시 통행 가능한 차량들의 리스트를 나타내는 포인터
	mu                    sync.Mutex
}

func (v *VehicleRPC) SendRequest(req Request, reply *Response) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	fmt.Printf("Vehicle : %v\n", req.Vehicle)
	req.Vehicle.Votes++

	// 모든 차량에게 요청을 전송
	var wg sync.WaitGroup

	for i := 0; i < req.TotalVehicles; i++ {
		if i != req.Vehicle.Address {
			address := fmt.Sprintf("localhost:%d", 8000+i)
			fmt.Printf("Vehicle %d attempting to connect to vehicle %d at %s\n", req.Vehicle.Address, i, address)

			wg.Add(1) // wg에 작업 추가
			go func(addr string, vehicleAddr int) {
				defer wg.Done() // 작업 완료 후 wg에서 제거

				// 슬라이스 인덱스 검증
				if vehicleAddr >= len(v.Vehicles) {
					fmt.Printf("Vehicle %d has an invalid index %d\n", req.Vehicle.Address, vehicleAddr)
					return
				}

				client, err := rpc.Dial("tcp", addr)
				if err != nil {
					fmt.Printf("Failed to connect to vehicle %d at %s: %v\n", vehicleAddr, addr, err)
					return
				}
				defer client.Close()

				fmt.Printf("Attempting RPC call to vehicle %d at %s\n", vehicleAddr, addr)

				var response Response

				// 요청 데이터 수정
				req := Request{
					Vehicle:       v.Vehicles[vehicleAddr], // 서버의 차량 목록에서 올바른 차량 정보 가져오기
					TotalVehicles: req.TotalVehicles,
				}

				err = client.Call("VehicleRPC.SendResponse", req, &response)
				if err != nil {
					fmt.Printf("RPC call error to vehicle %d: %v\n", vehicleAddr, err)
				} else {
					fmt.Printf("Vehicle %d successfully sent request to vehicle %d.\n", req.Vehicle.Address, vehicleAddr)
				}
			}(address, i)
		}
	}

	wg.Wait()

	return nil
}

func (v *VehicleRPC) SendResponse(req Request, reply *Response) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Print received request
	fmt.Printf("Vehicle %d received response request: %v\n", v.Address, req.Vehicle)

	// Example logic: just echo the vehicle data back as the response
	reply.Vehicle = req.Vehicle

	// Print response being sent
	fmt.Printf("Vehicle %d sending response: %v\n", v.Address, reply.Vehicle)

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

func startServer(address int, vehicles []Vehicle, totalVehicles int) {
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
	api := new(VehicleRPC)
	err := rpc.Register(api)
	if err != nil {
		log.Fatal("error registering API", err)
	}

	rpc.HandleHTTP()

	vehicles := []Vehicle{
		{Number: "1", Direction: "N", Address: 0, Votes: 0},
		{Number: "2", Direction: "S", Address: 1, Votes: 0},
		{Number: "3", Direction: "E", Address: 2, Votes: 0},
	}

	totalVehicles := len(vehicles)

	for i := 0; i < totalVehicles; i++ {
		go startServer(i, vehicles, totalVehicles)
	}

	time.Sleep(80 * time.Second) // 충분한 시간 대기
}
