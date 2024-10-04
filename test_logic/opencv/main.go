package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	licenseplate "opencv/licensePlate"
	pb "opencv/server/proto" // 패키지 경로 확인

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const GO_SERVER_PORT = 5005

var VEHICLES []int32

type BndBox struct {
	Xmin int `xml:"xmin"`
	Ymin int `xml:"ymin"`
	Xmax int `xml:"xmax"`
	Ymax int `xml:"ymax"`
}

// 번호판 객체 구조체
type Object struct {
	Name   string `xml:"name"`
	BndBox BndBox `xml:"bndbox"`
}

// XML 파일의 주 구조체
type Annotation struct {
	Filename string   `xml:"filename"`
	Path     string   `xml:"path"`
	Objects  []Object `xml:"object"`
}

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
func startServer(address int32, license int32) {
	// 서버 인스턴스 생성 및 초기화
	s := &server{
		Port:    fmt.Sprintf("%d", GO_SERVER_PORT+int(address)),       // 포트 번호를 문자열로 변환
		Vehicle: &pb.Vehicle{Address: address, LicensePlate: license}, // 기본 차량 정보로 초기화
	}

	// TCP 리스너 생성
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GO_SERVER_PORT+int(address)))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Server is listening on port %d\n", GO_SERVER_PORT+int(address))

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
	vehicle1License := fmt.Sprintf("%d", s.Vehicle.LicensePlate)
	vehicle2License := fmt.Sprintf("%d", req.Vehicle.LicensePlate)

	fmt.Printf("Port %s license : %s\n", s.Port, vehicle1License)
	fmt.Printf("Port %s license : %s\n", req.Port, vehicle2License)

	fmt.Printf("Vehicle %d received request from port %s: %v\n", s.Vehicle.Address, req.Port, req.Vehicle)

	// 방향의 유효성을 검사하고 응답 메시지 작성
	PassStatus := "False"
	if licenseplate.LicensePlateComparison(vehicle2License, vehicle1License) {
		PassStatus = "True"

		response := &pb.Response{
			Message:         fmt.Sprintf("Vote registered from port %s to port %s", req.Port, s.Port),
			Status:          "passed",
			DirectionStatus: PassStatus,
			Vehicle:         s.Vehicle, // 현재 서버의 차량 정보를 포함
		}

		return response, nil
	} else {
		response := &pb.Response{
			Message:         fmt.Sprintf("Vote registered from port %s to port %s", req.Port, s.Port),
			Status:          "ignored",
			DirectionStatus: PassStatus,
			Vehicle:         s.Vehicle, // 현재 서버의 차량 정보를 포함
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
	var TOTAL_VEHICLS int32 = 3

	for i := int32(0); i < TOTAL_VEHICLS; i++ {
		VEHICLES = append(VEHICLES, i)
	}

	var vehicle_license_list []int32 // 방향 저장 리스트

	// 서버를 각각의 고루틴에서 실행
	for i := int32(0); i < TOTAL_VEHICLS; i++ {
		vehicle_license_list = append(vehicle_license_list, int32(rand.Intn(249)+1))
		go startServer(i, vehicle_license_list[i]) // 랜덤 방향 같이 전달
	}

	// 잠시 대기하여 서버가 시작될 시간 확보
	// time.Sleep(time.Second * 2)

	// WaitGroup을 사용하여 모든 고루틴이 끝날 때까지 대기
	var wg sync.WaitGroup
	wg.Add(int(TOTAL_VEHICLS))

	// 모든 서버의 데이터를 저장할 맵
	serverData := make(map[int32]*pb.Vehicle)
	var dataMu sync.Mutex // 맵 접근을 위한 뮤텍스

	startTime := time.Now()
	var count = 0

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
								Address:      i,
								LicensePlate: vehicle_license_list[i],
							},
							Port:          fmt.Sprintf("%d", GO_SERVER_PORT+i),
							TotalVehicles: TOTAL_VEHICLS,
						},
					)
					count++
					if err != nil {
						log.Fatalf("grpc call error: %v", err)
					}

					// 투표가 "acknowledged"인 경우에만 ReceiveVotes를 증가
					if r.Status == "passed" {
						dataMu.Lock()
						vehicle, exists := serverData[i]
						if !exists {
							vehicle = &pb.Vehicle{
								Address:      i,
								ReceiveVotes: 0,
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
								ReceiveVotes: 0,
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
	for len(VEHICLES) > 0 {
		for addr, vehicle := range serverData {
			if !Contains(VEHICLES, vehicle.Address) {
				continue
			}

			fmt.Printf("\n")
			fmt.Printf("Address: %d, Vehicle: %+v\n", addr, vehicle)
			if vehicle.ReceiveVotes == TOTAL_VEHICLS-1 {
				VEHICLES = RemoveValue(VEHICLES, vehicle.Address)
				TOTAL_VEHICLS -= 1
				fmt.Printf("Address %d Vehicle passed \n", vehicle.Address)
				if len(VEHICLES) == 0 {
					break outerLoop
				}
			}
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	dataMu.Unlock()

	fmt.Printf("\n")
	fmt.Printf("합의 과정에 걸린 시간: %v\n", duration)

	fmt.Printf("남은 차량: %d \n", VEHICLES)
	fmt.Printf("번호판 비교 횟수: %d\n", count)
}
