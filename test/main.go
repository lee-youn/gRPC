package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "test/server/proto" // 패키지 경로 확인

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
}

// 서버 인스턴스 생성 및 초기화 후 실행
func startServer(address int32) {
	// 서버 인스턴스 생성 및 초기화
	s := &server{
		Port:    fmt.Sprintf("%d", GO_SERVER_PORT+address),      // 포트 번호를 문자열로 변환
		Vehicle: &pb.Vehicle{Address: address, ReceiveVotes: 0}, // 기본 차량 정보로 초기화
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

	if s.Vehicle.SendVotes == 0 {
		// 투표 응답 response
		fmt.Printf("Vehicle %d received response request: %v\n", s.Vehicle.Address, req.Vehicle)
		s.Vehicle.SendVotes = 1 // 투표 완료로 설정

		// 응답 메시지 작성
		response := &pb.Response{
			Message: fmt.Sprintf("Vote registered from port %s to port %s", req.Port, s.Port),
			Status:  "acknowledged",
			Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
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

func main() {

	// 교차로가 많아질수록? 5차선 6차선일 때 합의가 어려워짐..ㅠㅠ
	var randomNum int32 = int32(rand.Intn(5))
	var TOTAL_VEHICLS int32 = randomNum

	for i := int32(0); i < TOTAL_VEHICLS; i++ {
		VEHICLES = append(VEHICLES, i)
	}

	if TOTAL_VEHICLS == 0 {
		fmt.Printf("There are no vehicles entering at the same time.\n")
		return
	}

	if TOTAL_VEHICLS == 1 {
		fmt.Printf("Address %d Vehicle passed \n", 0)
		VEHICLES = RemoveValue(VEHICLES, 0)
		return
	}

	if TOTAL_VEHICLS == 2 {

		// 서버를 각각의 고루틴에서 실행
		for i := int32(0); i < TOTAL_VEHICLS; i++ {
			go startServer(i)
		}

		// 잠시 대기하여 서버가 시작될 시간 확보
		time.Sleep(time.Second * 2)

		var vehicle_number1 int32 = int32(rand.Intn(2))
		var vehicle_number2 int32 = int32(rand.Intn(2))

		fmt.Printf("동시진입차량: %d \n", VEHICLES)

		vehicles := []*pb.Vehicle{
			{
				Address:      0,
				RandomNumber: vehicle_number1,
			},
			{
				Address:      1,
				RandomNumber: vehicle_number2,
			},
		}

		fmt.Printf("랜덤 숫자: %d, %d\n", vehicle_number1, vehicle_number2)

		// 서버 연결을 미리 설정
		connections := make([]struct {
			client pb.VehicleServiceClient
			conn   *grpc.ClientConn
			ctx    context.Context
			cancel context.CancelFunc
		}, TOTAL_VEHICLS)

		for i := int32(0); i < TOTAL_VEHICLS; i++ {
			addr := fmt.Sprintf("localhost:%d", GO_SERVER_PORT+i)
			client, conn, ctx, cancel, err := rpcConnectTo(addr)
			if err != nil {
				log.Fatalf("connect error: %v", err)
			}
			connections[i] = struct {
				client pb.VehicleServiceClient
				conn   *grpc.ClientConn
				ctx    context.Context
				cancel context.CancelFunc
			}{client, conn, ctx, cancel}
		}

		startTime := time.Now()

		for len(VEHICLES) > 0 {
			var wg sync.WaitGroup
			wg.Add(len(VEHICLES))

			// 합의 상태 추적
			var allVehiclesPassed bool

			connData := connections[0]

			r, err := connData.client.RandomAgreement(connData.ctx, &pb.Request{
				Vehicle:       vehicles[0],
				Port:          fmt.Sprintf("%d", GO_SERVER_PORT+0),
				TotalVehicles: TOTAL_VEHICLS,
				RandomNumber:  vehicles[1].RandomNumber,
			})
			if err != nil {
				log.Fatalf("grpc call error: %v", err)
			}

			if r.Status == "pass" {
				VEHICLES = RemoveValue(VEHICLES, 0)
				fmt.Printf("Address %d Vehicle passed \n", 0)
				remainingVehicle := VEHICLES[0]
				VEHICLES = RemoveValue(VEHICLES, remainingVehicle)
				fmt.Printf("Address %d Vehicle passed \n", remainingVehicle)
				break

			} else if r.Status == "wait" {
				fmt.Printf("Address %d Vehicle has to wait \n", 0)
				VEHICLES = RemoveValue(VEHICLES, 1)
				fmt.Printf("Address %d Vehicle passed \n", 1)
				remainingVehicle := VEHICLES[0]
				VEHICLES = RemoveValue(VEHICLES, remainingVehicle)
				fmt.Printf("Address %d Vehicle passed \n", remainingVehicle)
				break

			} else if r.Status == "draw" {
				vehicles[0].RandomNumber = int32(rand.Intn(2))
				vehicles[1].RandomNumber = int32(rand.Intn(2))
				fmt.Printf("랜덤 숫자: %d, %d\n", vehicles[0].RandomNumber, vehicles[1].RandomNumber)
				continue
			}

			// 모든 차량이 통과했는지 확인
			if len(VEHICLES) == 0 {
				allVehiclesPassed = true
			}
			// 모든 고루틴이 완료될 때까지 대기
			wg.Wait()

			// 모든 차량이 통과했으면 종료
			if allVehiclesPassed {
				break
			}
		}
		endTime := time.Now()
		duration := endTime.Sub(startTime)

		fmt.Printf("합의 과정에 걸린 시간: %v\n", duration)
		return
	}

	// 서버를 각각의 고루틴에서 실행
	for i := int32(0); i < TOTAL_VEHICLS; i++ {
		go startServer(i)
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
							Direction: "North",
							Address:   i,
						},

						Port:          fmt.Sprintf("%d", GO_SERVER_PORT+i),
						TotalVehicles: 3,
					},
				)
				if err != nil {
					log.Fatalf("grpc call error: %v", err)
				}

				// 투표가 "acknowledged"인 경우에만 ReceiveVotes를 증가
				if r.Status == "acknowledged" {
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
				}

				log.Printf("Response from port %d: %v", GO_SERVER_PORT+j, r.Vehicle)
			}
		}(i)
	}

	// 모든 고루틴이 완료될 때까지 대기
	wg.Wait()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// 최종 서버 데이터를 출력
	fmt.Println("Final server data:")
	dataMu.Lock()
	for addr, vehicle := range serverData {
		fmt.Printf("Address: %d, Vehicle: %+v\n", addr, vehicle)
		if vehicle.ReceiveVotes == TOTAL_VEHICLS {
			VEHICLES = RemoveValue(VEHICLES, vehicle.Address)
			fmt.Printf("Address %d Vehicle passed \n", vehicle.Address)
		}
	}
	dataMu.Unlock()

	fmt.Printf("합의 과정에 걸린 시간: %v\n", duration)

	fmt.Printf("남은 차량: %d \n", VEHICLES)
}
