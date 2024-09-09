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
	pbtimestamp "google.golang.org/protobuf/types/known/timestamppb"
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

func (s *server) LeaderElection(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	fmt.Printf("#LeaderElection # Port %s made a request to port %s.\n", req.Port, s.Port)

	if s.Vehicle.ReceiveVotes < req.Vehicle.ReceiveVotes {
		fmt.Printf("요청 차량의 투표수가 더 많습니다.\n")
		response := &pb.Response{
			Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Address),
			Status:  "acknowledged",
			Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
		}
		return response, nil
	} else {
		fmt.Printf("요청 차량의 투표수가 더 적거나 동일합니다.")
		if s.Vehicle.ReceiveVotes == req.Vehicle.ReceiveVotes {
			reqTime := req.Vehicle.ElectionTime.AsTime().UnixNano()
			serverTime := s.Vehicle.ElectionTime.AsTime().UnixNano()

			if reqTime > serverTime {
				fmt.Printf("요청 차량의 선거 시간이 더 최신입니다.\n")
				response := &pb.Response{
					Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Address),
					Status:  "acknowledged",
					Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
				}
				return response, nil
			} else {
				fmt.Printf("요청 차량의 선거 시간이 더 예전입니다.\n")
				response := &pb.Response{
					Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Address),
					Status:  "ignored",
					Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
				}
				return response, nil
			}
		}
		fmt.Printf("요청 차량의 투표가 더 적습니다. \n")
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

func (s *server) UpdateVoteCount(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("#UpdateVoteCount# Do the server and client addresses match?")
	fmt.Printf("server : %d, client : %d\n", s.Vehicle.Address, req.Vehicle.Address)

	if s.Vehicle.Address == req.Vehicle.Address {
		// 요청에서 받은 값으로 ReceiveVotes 업데이트
		s.Vehicle.ReceiveVotes = req.Vehicle.ReceiveVotes
		s.Vehicle.ElectionTime = req.Vehicle.ElectionTime
		fmt.Printf("차량 %d의 투표 수를 %d로 업데이트했습니다.\n", s.Vehicle.Address, s.Vehicle.ReceiveVotes)

		// 성공을 나타내는 응답 준비 및 반환
		return &pb.Response{
			Message: fmt.Sprintf("차량 %d의 투표 수가 업데이트되었습니다.\n", s.Vehicle.Address),
			Status:  "success",
		}, nil
	} else {
		fmt.Printf("주소 불일치. 차량 %d의 업데이트에 실패했습니다.\n", req.Vehicle.Address)
		return &pb.Response{
			Message: fmt.Sprintf("주소 불일치. 차량 %d의 업데이트에 실패했습니다.\n", req.Vehicle.Address),
			Status:  "failed",
		}, nil
	}
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

// 차량과 관련된 차량을 VECHICLES 리스트에서 제거하는 함수
func removeVehiclesIfQuorumReached(vehicle *pb.Vehicle) bool {
	VEHICLES = RemoveValue(VEHICLES, vehicle.Address)

	// Covehicles 처리
	for _, covehicle := range vehicle.Covehicle {
		if !Contains(VEHICLES, covehicle.Address) {
			continue
		}
		VEHICLES = RemoveValue(VEHICLES, covehicle.Address)
		fmt.Printf("Address %d Vehicle passed \n", covehicle.Address)

		// Covehicle의 Covehicle 처리
		for _, subCovehicle := range covehicle.Covehicle {
			if !Contains(VEHICLES, subCovehicle.Address) {
				continue
			}
			VEHICLES = RemoveValue(VEHICLES, subCovehicle.Address)
			fmt.Printf("Address %d Vehicle passed \n", subCovehicle.Address)
		}
	}
	fmt.Printf("Address %d Vehicle passed \n", vehicle.Address)

	return len(VEHICLES) == 0
}

func main() {

	// 교차로가 많아질수록? 5차선 6차선일 때 합의가 어려워짐..ㅠㅠ
	var randomNum int32 = int32(rand.Intn(5))
	var TOTAL_VEHICLS int32 = randomNum
	var QUORUM int32 = (TOTAL_VEHICLS / 2) + 1

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
			go startServer(i, "0")
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

	// 3대 이상일 때
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
	var done = false

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

					if done { // 종료 플래그가 설정된 경우
						return
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

					if r.Status == "acknowledged" {
						dataMu.Lock()
						defer dataMu.Unlock()

						vehicle, exists := serverData[i]

						// 기존 vehicle이 있으면 가져오고 없으면 새로운 vehicle 생성
						if !exists {
							vehicle = &pb.Vehicle{
								Address:      i,
								ReceiveVotes: 0,
								ElectionVote: 0,
							}
						}

						// Covehicle 처리 (DirectionStatus가 "True"인 경우에만)
						if r.DirectionStatus == "True" {
							var covehicles []*pb.Vehicle
							covehicles = vehicle.Covehicle // 기존의 covehicles 가져오기

							// covehicles 리스트를 순회하면서 같은 방향의 차량을 찾고 추가
							for i := int32(0); i < int32(len(covehicles)); i++ {
								var linked_covehicles []*pb.Vehicle
								if direction.DirectionBoolean(covehicles[i].Direction, r.ResponseDirection) {
									linked_covehicles = covehicles[i].Covehicle
									linked_covehicles = append(linked_covehicles, r.Vehicle)
									covehicles[i].Covehicle = linked_covehicles
								}
							}

							// 새로 들어온 차량을 covehicles에 추가
							covehicles = append(covehicles, r.Vehicle)
							vehicle.Covehicle = covehicles
						}

						// 요청 시간 저장
						vehicle.ElectionTime = pbtimestamp.Now()

						// 투표 수 증가
						vehicle.ReceiveVotes++

						// 증가된 투표수를 서버에 반영하는 요청 추가
						addr := fmt.Sprintf("localhost:%d", GO_SERVER_PORT+i) // 자신의 서버 주소
						client, conn, ctx, cancel, err := rpcConnectTo(addr)
						if err != nil {
							log.Fatalf("connect error: %v", err)
						}
						defer conn.Close()
						defer cancel()

						// UpdateVoteCount 호출
						_, err = client.UpdateVoteCount(ctx, &pb.Request{
							Vehicle: vehicle,
						})
						if err != nil {
							log.Printf("failed to update vote count: %v", err)
						}

						// vehicle을 업데이트
						serverData[i] = vehicle

						// 만약 Receive Votes >= Quorum이면 자신을 제외한 차량에 LeaderElection 요청을 보내기
						if vehicle.ReceiveVotes >= QUORUM {
							fmt.Printf("차량 %d의 ReceiveVotes가 Quorum에 도달했습니다. LeaderElection 요청을 보냅니다.\n", vehicle.Address)

							var wg sync.WaitGroup
							// 고루틴으로 LeaderElection 요청 전송
							for k := int32(0); k < TOTAL_VEHICLS; k++ {
								if k == i {
									continue // 자기 자신에게는 요청을 보내지 않음
								}

								wg.Add(1) // 고루틴 수 증가
								go func(k int32) {
									defer wg.Done() // 고루틴 완료 시 호출
									if done {       // 종료 플래그가 설정된 경우
										return
									}
									addr := fmt.Sprintf("localhost:%d", GO_SERVER_PORT+k)
									client, conn, ctx, cancel, err := rpcConnectTo(addr)
									if err != nil {
										fmt.Printf("connect error to %s: %v\n", addr, err)
										return
									}
									defer conn.Close()
									defer cancel()

									// LeaderElection 요청 전송
									r, err = client.LeaderElection(ctx, &pb.Request{
										Vehicle: vehicle,
									})
									if err != nil {
										fmt.Printf("LeaderElection 요청 실패: %v\n", err)
									} else {
										fmt.Printf("차량 %d에 LeaderElection 요청을 성공적으로 보냈습니다.\n", k)
									}

									if r.Status == "acknowledged" {
										vehicle.ElectionVote++
										if vehicle.ElectionVote >= QUORUM {
											removeVehiclesIfQuorumReached(vehicle)
											done = true // 모든 고루틴 종료 플래그 설정
											return
										}
									}

								}(k) // 고루틴에 j를 전달
							}

							// 모든 고루틴이 완료될 때까지 대기
							wg.Wait()
						}
					} else {
						dataMu.Lock()
						vehicle, exists := serverData[i]
						if !exists {
							vehicle = &pb.Vehicle{
								Address:      i,
								ReceiveVotes: 1,
								ElectionVote: 0,
							}
							serverData[i] = vehicle
						}
						dataMu.Unlock()
					}

				}(j)
			}
		}(i)
	}

	// 모든 고루틴이 완료될 때까지 대기
	wg.Wait()
	dataMu.Lock()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	dataMu.Unlock()

	fmt.Printf("\n")
	fmt.Printf("합의 과정에 걸린 시간: %v\n", duration)

	fmt.Printf("남은 차량: %d \n", VEHICLES)
}
