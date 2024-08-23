package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	direction "direction/directionBoolean"
	pb "direction/server/proto" // 패키지 경로 확인

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const GO_SERVER_PORT = 5005

var VEHICLES []int32

////////////////////////////////////////
/////////////////SERVER/////////////////
////////////////////////////////////////

// 서버 구조체 정의
type server struct {
	pb.UnimplementedVehicleServiceServer
	Port    string
	Vehicle *pb.Vehicle // 포트당 하나의 차량 정보를 저장
	mu      sync.Mutex
}

// 서버 인스턴스 생성 및 초기화 후 실행
func startServer(address int32, direction string) {
	// 서버 인스턴스 생성 및 초기화
	s := &server{
		Port:    fmt.Sprintf("%d", GO_SERVER_PORT+address),           // 포트 번호를 문자열로 변환
		Vehicle: &pb.Vehicle{Address: address, Direction: direction}, // 기본 차량 정보로 초기화
	}

	// TCP 리스너 생성
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GO_SERVER_PORT+address))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Server is listening on port %d\n", GO_SERVER_PORT+address)

	// gRPC 서버 생성
	grpcServer := grpc.NewServer()
	pb.RegisterVehicleServiceServer(grpcServer, s)

	// 서버 시작 (차단 함수)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (s *server) ReceiveRequest(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	fmt.Printf("Port %s made a request to port %s.\n", req.Port, s.Port)

	// 동시에 접근 못하게 함(즉, 동시에 request가 도착하는 일 없음)
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("Port %s direction : %s\n", s.Port, s.Vehicle.Direction)
	fmt.Printf("Port %s direction : %s\n", req.Port, req.Vehicle.Direction)

	if s.Vehicle.SendVotes == 0 {
		// 투표 응답을 처리하는 경우
		fmt.Printf("Vehicle %d received request from port %s: %v\n", s.Vehicle.Address, req.Port, req.Vehicle)
		s.Vehicle.SendVotes = 1 // 투표 완료로 설정

		// 방향의 유효성을 검사하고 응답 메시지 작성
		directionStatus := "False"
		if direction.DirectionBoolean(req.Vehicle.Direction, s.Vehicle.Direction) {
			directionStatus = "True"
		}

		response := &pb.Response{
			Message:           fmt.Sprintf("Vote registered from port %s to port %s", req.Port, s.Port),
			Status:            "acknowledged",
			DirectionStatus:   directionStatus,
			ResponseDirection: s.Vehicle.Direction,
			Vehicle:           s.Vehicle, // 현재 서버의 차량 정보를 포함
		}

		return response, nil
	} else {
		// 이미 투표한 경우
		fmt.Printf("Vehicle %d has already responded to a request.\n", s.Vehicle.Address)

		// 응답 메시지 작성
		response := &pb.Response{
			Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Address),
			Status:  "ignored",
			Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
		}

		return response, nil
	}
}

func (s *server) RandomAgreement(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	fmt.Printf("Port %s made a request to port %s.\n", req.Port, s.Port)

	s.Vehicle.RandomNumber = req.RandomNumber

	var response *pb.Response

	if req.Vehicle.RandomNumber == s.Vehicle.RandomNumber {
		fmt.Printf("reqeust randomnumber: %d\n", req.Vehicle.RandomNumber)
		fmt.Printf("server randomnumber: %d\n", s.Vehicle.RandomNumber)
		fmt.Printf("Cannot reach an agreement.\n")
		response = &pb.Response{
			Message: fmt.Sprintf("Cannot reach an agreement"),
			Status:  "draw",
		}
	} else if req.Vehicle.RandomNumber > s.Vehicle.RandomNumber {
		fmt.Printf("reqeust randomnumber: %d\n", req.Vehicle.RandomNumber)
		fmt.Printf("server randomnumber: %d\n", s.Vehicle.RandomNumber)
		fmt.Printf("Vehicle %d may pass first.\n", req.Vehicle.Address)
		response = &pb.Response{
			Message: fmt.Sprintf("Vehicle %d may pass first.", req.Vehicle.Address),
			Status:  "pass",
		}
	} else {
		fmt.Printf("reqeust randomnumber: %d\n", req.Vehicle.RandomNumber)
		fmt.Printf("server randomnumber: %d\n", s.Vehicle.RandomNumber)
		fmt.Printf("Vehicle %d has to wait.\n", req.Vehicle.Address)
		response = &pb.Response{
			Message: fmt.Sprintf("Vehicle %d has to wait.", req.Vehicle.Address),
			Status:  "wait",
		}
	}

	return response, nil
}

////////////////////////////////////////
/////////////////CLIENT/////////////////
////////////////////////////////////////

func rpcConnectTo(ip string) (pb.VehicleServiceClient, *grpc.ClientConn, context.Context, context.CancelFunc, error) {
	conn, err := grpc.Dial(
		ip,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("did not connect: %v", err)
	}

	c := pb.NewVehicleServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	return c, conn, ctx, cancel, nil
}

func RemoveValue(slice []int32, valueToRemove int32) []int32 {
	result := []int32{}
	for _, value := range slice {
		if value != valueToRemove {
			result = append(result, value)
		}
	}
	return result
}

var TOTAL_VEHICLS int32

func Contains(slice []int32, item int32) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func main() {

	// 교차로가 많아질수록? 5차선 6차선일 때 합의가 어려워짐..ㅠㅠ
	// var randomNum int32 = int32(rand.Intn(5))
	var TOTAL_VEHICLS int32 = 4

	for i := int32(0); i < TOTAL_VEHICLS; i++ {
		VEHICLES = append(VEHICLES, i)
	}

	var directions = [12]string{"Rs", "Rl", "Rr", "Ls", "Ll", "Lr", "Ds", "Dl", "Dr", "Us", "Ul", "Ur"}

	var vehicle_direction_list []string // 방향 저장 리스트

	// 서버를 각각의 고루틴에서 실행
	for i := int32(0); i < TOTAL_VEHICLS; i++ {
		vehicle_direction_list = append(vehicle_direction_list, directions[rand.Intn(len(directions))])
		go startServer(i, vehicle_direction_list[i]) // 랜덤 방향 같이 전달
	}

	// 잠시 대기하여 서버가 시작될 시간 확보
	time.Sleep(time.Second * 2)

	// WaitGroup을 사용하여 모든 고루틴이 끝날 때까지 대기
	var wg sync.WaitGroup
	wg.Add(int(TOTAL_VEHICLS))

	// 모든 서버의 데이터를 저장할 맵
	serverData := make(map[int32]*pb.Vehicle)
	var dataMu sync.Mutex // 맵 접근을 위한 뮤텍스

	startTime := time.Now()

	// 각 서버에 요청 보내기
	for i := int32(0); i < TOTAL_VEHICLS; i++ {
		go func(i int32) {
			defer wg.Done()
			for j := int32(0); j < TOTAL_VEHICLS; j++ {
				if i == j {
					continue // 자기 자신에게는 요청을 보내지 않음
				}

				wg.Add(1) // 새로운 고루틴을 시작하므로 WaitGroup 카운터를 증가
				go func(j int32) {
					defer wg.Done()

					addr := fmt.Sprintf("localhost:%d", GO_SERVER_PORT+j)
					client, conn, ctx, cancel, err := rpcConnectTo(addr)
					if err != nil {
						log.Fatalf("connect error: %v", err)
					}
					defer conn.Close()
					defer cancel()

					r, err := client.ReceiveRequest(
						ctx,
						&pb.Request{
							Vehicle: &pb.Vehicle{
								Address:   i,
								Direction: vehicle_direction_list[i],
							},
							Port:          fmt.Sprintf("%d", GO_SERVER_PORT+i),
							TotalVehicles: TOTAL_VEHICLS,
						},
					)
					if err != nil {
						log.Fatalf("grpc call error: %v", err)
					}

					// 투표가 "acknowledged"인 경우에만 ReceiveVotes를 증가
					if r.Status == "acknowledged" && r.DirectionStatus == "True" {
						dataMu.Lock()
						vehicle, exists := serverData[i]

						var covehicles []*pb.Vehicle
						if exists {
							covehicles = vehicle.Covehicle // 기존의 covehicles 가져오기

							for i := int32(0); i < int32(len(covehicles)); i++ {
								var linked_covehicles []*pb.Vehicle
								if direction.DirectionBoolean(covehicles[i].Direction, r.ResponseDirection) {
									linked_covehicles = covehicles[i].Covehicle
									linked_covehicles = append(linked_covehicles, r.Vehicle)
									covehicles[i].Covehicle = linked_covehicles
									// if linked_covehicles != nil {
									// 	linked_covehicles = append(linked_covehicles, r.Vehicle)
									// }
								}
							}

							covehicles = append(covehicles, r.Vehicle)
							vehicle.Covehicle = covehicles
						}

						if !exists {
							covehicles = append(covehicles, r.Vehicle)

							vehicle = &pb.Vehicle{
								Address:      i,
								Covehicle:    covehicles,
								ReceiveVotes: 1,
							}
							serverData[i] = vehicle
						}
						vehicle.ReceiveVotes++
						dataMu.Unlock()

					} else if r.Status == "acknowledged" {
						dataMu.Lock()
						vehicle, exists := serverData[i]
						if !exists {
							vehicle = &pb.Vehicle{
								Address:      i,
								ReceiveVotes: 1,
							}
							serverData[i] = vehicle
						}
						vehicle.ReceiveVotes++
						dataMu.Unlock()
					} else {
						dataMu.Lock()
						vehicle, exists := serverData[i]
						if !exists {
							vehicle = &pb.Vehicle{
								Address:      i,
								ReceiveVotes: 1,
							}
							serverData[i] = vehicle
						}
						dataMu.Unlock()
					}

					// log.Printf("Response from port %d: %v", GO_SERVER_PORT+j, r.Vehicle)
				}(j)
			}
		}(i)
	}

	// 모든 고루틴이 완료될 때까지 대기
	wg.Wait()

	// 최종 서버 데이터를 출력
	fmt.Println("Final server data:")
	dataMu.Lock()

outerLoop:
	for addr, vehicle := range serverData {

		if !Contains(VEHICLES, vehicle.Address) {
			continue
		}

		fmt.Printf("\n")
		fmt.Printf("Address: %d, Vehicle: %+v\n", addr, vehicle)
		fmt.Printf("Covehicle: %s\n", vehicle.Covehicle)
		for i := int32(0); i < int32(len(vehicle.Covehicle)); i++ {
			if !Contains(VEHICLES, vehicle.Covehicle[i].Address) {
				continue
			}
			fmt.Printf("Covehicle's: %s\n", vehicle.Covehicle[i])
			fmt.Printf("Covehicle's covehicle: %s\n", vehicle.Covehicle[i].Covehicle)
		}
		if vehicle.ReceiveVotes == TOTAL_VEHICLS {
			VEHICLES = RemoveValue(VEHICLES, vehicle.Address)
			for i := int32(0); i < int32(len(vehicle.Covehicle)); i++ {
				if !Contains(VEHICLES, vehicle.Covehicle[i].Address) {
					continue
				}

				VEHICLES = RemoveValue(VEHICLES, vehicle.Covehicle[i].Address)
				fmt.Printf("Address %d Vehicle passed \n", vehicle.Covehicle[i].Address)
				for j := int32(0); j < int32(len(vehicle.Covehicle[i].Covehicle)); j++ {
					if !Contains(VEHICLES, vehicle.Covehicle[i].Covehicle[j].Address) {
						continue
					}
					VEHICLES = RemoveValue(VEHICLES, vehicle.Covehicle[i].Covehicle[j].Address)
					fmt.Printf("Address %d Vehicle passed \n", vehicle.Covehicle[i].Covehicle[j].Address)
				}
			}
			fmt.Printf("Address %d Vehicle passed \n", vehicle.Address)
			if len(VEHICLES) == 0 {
				break outerLoop
			}
		}
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	dataMu.Unlock()

	fmt.Printf("\n")
	fmt.Printf("합의 과정에 걸린 시간: %v\n", duration)

	fmt.Printf("남은 차량: %d \n", VEHICLES)
}
