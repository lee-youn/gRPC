package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	direction "direction/directionBoolean"
	pb "direction/server/proto" // 패키지 경로 확인

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pbtimestamp "google.golang.org/protobuf/types/known/timestamppb"
)

var GO_SERVER_PORT = 50000

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
func startServer(address int32, direction string, number int32) (*grpc.Server, string) {
	// 서버 인스턴스 생성 및 초기화
	s := &server{
		Port:    fmt.Sprintf("%d", GO_SERVER_PORT+int(address)),                      // 포트 번호를 문자열로 변환
		Vehicle: &pb.Vehicle{Number: number, Address: address, Direction: direction}, // 기본 차량 정보로 초기화
	}

	// TCP 리스너 생성
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GO_SERVER_PORT+int(address)))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return nil, s.Port // 오류 발생 시 nil 반환
	}
	log.Printf("Server is listening on port %d\n", GO_SERVER_PORT+int(address))

	// gRPC 서버 생성
	grpcServer := grpc.NewServer(
		grpc.MaxSendMsgSize(1024*1024*10), // 10MB
		grpc.MaxRecvMsgSize(1024*1024*10), // 10MB
	)

	pb.RegisterVehicleServiceServer(grpcServer, s)

	// 서버 시작 (차단 함수)를 고루틴에서 실행
	go func() {
		defer lis.Close() // 리스너를 고루틴 내에서 닫도록 수정
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
		log.Printf("grpcServer %d 실행 중\n", address)
	}()

	log.Printf("grpcServer 함수 테스트 \n")

	return grpcServer, s.Port
}

func (s *server) ReceiveRequest(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	// log.Printf("Port %s made a request to port %s.\n", req.Port, s.Port)

	// 동시에 접근 못하게 함(즉, 동시에 request가 도착하는 일 없음)
	s.mu.Lock()
	defer s.mu.Unlock()

	if req == nil {
		return nil, fmt.Errorf("received nil request")
	}

	log.Printf("Port %s direction : %s\n", s.Port, s.Vehicle.Direction)
	log.Printf("Port %s direction : %s\n", req.Port, req.Vehicle.Direction)

	if s.Vehicle.SendVotes == 0 {
		// 투표 응답을 처리하는 경우
		log.Printf("Vehicle %d received request from port %s: %v\n", s.Vehicle.Number, req.Port, req.Vehicle)
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
		log.Printf("Vehicle %d has already responded to a request.\n", s.Vehicle.Number)

		// 응답 메시지 작성
		response := &pb.Response{
			Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Number),
			Status:  "ignored",
			Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
		}

		return response, nil
	}
}

func (s *server) LeaderElection(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	log.Printf("\n#LeaderElection#\n Port (%s) made a request to port (%+v).\n\n", req.Vehicle, s.Vehicle)
	log.Printf("vehicle test : %v\n", s.Vehicle.ReceiveVotes)

	if s.Vehicle.ReceiveVotes < req.Vehicle.ReceiveVotes {
		log.Printf("요청 차량의 투표수가 더 많습니다.\n")
		response := &pb.Response{
			Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Number),
			Status:  "acknowledged",
			Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
		}
		return response, nil
	} else {
		log.Printf("요청 차량의 투표수가 더 적거나 동일합니다.")
		if s.Vehicle.ReceiveVotes == req.Vehicle.ReceiveVotes {
			reqTime := req.Vehicle.ElectionTime.AsTime().UnixNano()
			serverTime := s.Vehicle.ElectionTime.AsTime().UnixNano()

			if reqTime > serverTime {
				log.Printf("요청 차량의 선거 시간이 더 최신입니다.\n")
				response := &pb.Response{
					Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Number),
					Status:  "acknowledged",
					Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
				}
				return response, nil
			} else {
				log.Printf("요청 차량의 선거 시간이 더 예전입니다.\n")
				response := &pb.Response{
					Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Number),
					Status:  "ignored",
					Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
				}
				return response, nil
			}
		}
		log.Printf("요청 차량의 투표가 더 적습니다. \n")
		response := &pb.Response{
			Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Number),
			Status:  "ignored",
			Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
		}
		return response, nil
	}

}

func (s *server) UpdateVoteCount(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	log.Printf("#UpdateVoteCount# Do the server and client addresses match?")
	log.Printf("server : %d, client : %d\n", s.Vehicle.Number, req.Vehicle.Number)

	if s.Vehicle.Number == req.Vehicle.Number {
		// 요청에서 받은 값으로 ReceiveVotes 업데이트
		s.Vehicle.ReceiveVotes = req.Vehicle.ReceiveVotes
		s.Vehicle.ElectionTime = req.Vehicle.ElectionTime
		log.Printf("잘업데이트됐나테스트 : %+v\n", s.Vehicle)
		log.Printf("차량 %d의 투표 수를 %d로 업데이트했습니다.\n", s.Vehicle.Number, s.Vehicle.ReceiveVotes)

		// 성공을 나타내는 응답 준비 및 반환
		return &pb.Response{
			Message: fmt.Sprintf("차량 %d의 투표 수가 업데이트되었습니다.\n", s.Vehicle.Number),
			Status:  "success",
		}, nil
	} else {
		log.Printf("주소 불일치. 차량 %d의 업데이트에 실패했습니다.\n", req.Vehicle.Number)
		return &pb.Response{
			Message: fmt.Sprintf("주소 불일치. 차량 %d의 업데이트에 실패했습니다.\n", req.Vehicle.Number),
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
	if !Contains(slice, valueToRemove) {
		log.Printf("차량 %v는 이미 통과된 차량입니다. \n", valueToRemove)
		return slice
	} else {
		result := []int32{}
		for _, value := range slice {
			if value != valueToRemove {
				result = append(result, value)
			}
		}
		return result
	}
}

var TOTAL_VEHICLES int32

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

	VEHICLES = RemoveValue(VEHICLES, vehicle.Number)
	PASS_COUNT++
	log.Printf("제거할 차량 : %v, 제거됨? : %v\n", vehicle, VEHICLES)

	if len(vehicle.Covehicle) > 0 {
		firstCovehicle := vehicle.Covehicle[0]
		if Contains(VEHICLES, firstCovehicle.Number) {
			log.Printf("통과할 차량 번호: %d \n", firstCovehicle.Number)
			log.Printf("통과할 covehicle 차량정보 : %v \n", firstCovehicle)
			VEHICLES = RemoveValue(VEHICLES, firstCovehicle.Number)
			log.Printf("Address %d Vehicle passed \n", firstCovehicle.Number)
			PASS_COUNT++
		}
		if firstCovehicle.Covehicle != nil && len(firstCovehicle.Covehicle) > 0 {
			for _, subCovehicle := range vehicle.Covehicle[0].Covehicle {
				if !Contains(VEHICLES, subCovehicle.Number) {
					continue
				}
				VEHICLES = RemoveValue(VEHICLES, subCovehicle.Number)
				log.Printf("Address %d Vehicle passed \n", subCovehicle.Number)
				PASS_COUNT++
			}
		}
	}

	log.Printf("Address %d Vehicle passed \n", vehicle.Number)

	return len(VEHICLES) == 0
}

// 주어진 차량 배열에서 n개의 차량을 무작위로 선택하는 함수
func selectRandomVehicles(vehicles []int, n int) []int {
	selected := make([]int, 0, n)
	vehicleCopy := append([]int{}, vehicles...) // 원본 배열을 복사
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(vehicleCopy))
		selected = append(selected, vehicleCopy[idx])
		vehicleCopy = append(vehicleCopy[:idx], vehicleCopy[idx+1:]...)
	}
	return selected
}

// 두 배열 간의 차집합을 구하는 함수
func difference(a, b []int) []int {
	m := make(map[int]bool)
	for _, item := range b {
		m[item] = true
	}
	diff := []int{}
	for _, item := range a {
		if !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

var PASS_COUNT int

func main() {
	// var TOTOA_TIME float64
	// const testNumber = 0
	filePath := "/Users/leeyounjeong/Documents/LAB/실험로그/50_04.txt"

	// 로그 파일 생성
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("로그 파일을 열 수 없습니다: %v", err)
	}
	defer file.Close()

	// 로그 출력 대상을 파일로 설정
	log.SetOutput(file)
	totalStartTime := time.Now()

	const NUMBER_OF_TOTAL_VEHICLES = 50
	hvRatio := 0.4

	// 전체 차량 배열 생성
	totalVehicles := make([]int, NUMBER_OF_TOTAL_VEHICLES)
	for i := 0; i < NUMBER_OF_TOTAL_VEHICLES; i++ {
		totalVehicles[i] = i + 1 // 차량 번호 1부터 360까지 설정
	}

	// 비잔틴 장애 차량 수 계산
	numHV := int(float64(NUMBER_OF_TOTAL_VEHICLES) * hvRatio)
	// 비잔틴 장애 차량 배열
	hvVehicles := selectRandomVehicles(totalVehicles, numHV)
	// 나머지 정상 차량 배열 (av 차량)
	avVehicles := difference(totalVehicles, hvVehicles)
	// 결과 출력
	log.Printf("비잔틴 장애 차량 (hv): %v\n", hvVehicles)
	log.Printf("정상 차량 (av): %v\n", avVehicles)

	for len(totalVehicles) > 0 {
		TIMEOUT := time.Now()

		var randomNum int32 = int32(rand.Intn(4) + 1) // 1 또는 2
		if randomNum > int32(len(totalVehicles)) {
			randomNum = int32(len(totalVehicles)) // 차량 수보다 많지 않도록 조정
		}
		var TOTAL_VEHICLES int32 = randomNum
		selectedVehicles := selectRandomVehicles(totalVehicles, int(TOTAL_VEHICLES))
		log.Printf("선택된 차량: %v\n", selectedVehicles)

		var RandomByzantine []int32

		for i := 0; i < len(selectedVehicles); i++ {
			VEHICLES = append(VEHICLES, int32(selectedVehicles[i])) // i 값을 VEHICLES에 추가
			if contains(hvVehicles, selectedVehicles[i]) {
				RandomByzantine = append(RandomByzantine, int32(selectedVehicles[i]))
			}
		}

		var STOP_VEHICLES_PASS_TIME int

		for len(VEHICLES) > 0 {
			END_TIMEOUT := time.Now()
			duration := END_TIMEOUT.Sub(TIMEOUT)
			log.Printf("TIME : %v\n", duration)
			log.Printf("TIME2 : %v\n", STOP_VEHICLES_PASS_TIME)
			log.Printf("비교시간 : %s\n", time.Duration(500)*time.Millisecond)

			if len(VEHICLES) > 1 && duration-(time.Duration(STOP_VEHICLES_PASS_TIME)*time.Millisecond) >= time.Duration(500)*time.Millisecond {
				log.Printf("TIMEOUT! 비전시스템으로 합의 진행\n")
				log.Printf("TIME3 : %v\n", duration-(time.Duration(STOP_VEHICLES_PASS_TIME)*time.Millisecond))
				time.Sleep(time.Duration(500) * time.Millisecond)

				for _, i := range VEHICLES {
					log.Printf("%v \n", VEHICLES)

					// 리스트가 비어있지 않다면 첫 번째 차량에 접근
					if len(VEHICLES) > 0 && !Contains(RandomByzantine, i) {
						log.Printf("Address %d Vehicle passed \n", i)
						VEHICLES = RemoveValue(VEHICLES, i)
						PASS_COUNT++
					} else {
						log.Printf("Address %d Vehicle은 무응답 차량이므로 통과하지 못합니다. \n", i)
					}
				}
				break
			}

			var TOTAL_VEHICLES int32 = int32(len(VEHICLES))
			// var QUORUM int32 = (int32(TOTAL_VEHICLES)+1)/2 + 1 // 리더는 자기 자신한테 투표를 하니깐..1을 더 더해줘야함
			var QUORUM int32 = TOTAL_VEHICLES

			if TOTAL_VEHICLES == 0 {
				log.Printf("There are no vehicles entering at the same time.\n")
				log.Printf("\n")
				return
			}

			if TOTAL_VEHICLES == 1 {
				log.Printf("Address %d Vehicle passed \n", 0)
				PASS_COUNT++
				log.Printf("\n")
				VEHICLES = RemoveValue(VEHICLES, 0)
				break
			}

			if TOTAL_VEHICLES == 2 {
				log.Printf("TOTAL_VEHICLES : %v\n", VEHICLES)

				var VISION = 0

				for _, i := range VEHICLES {
					if Contains(RandomByzantine, i) {
						break
					} else {
						VISION++
					}
				}

				if VISION != 2 {
					var RANDOM_PASS_TIME = 3000
					STOP_VEHICLES_PASS_TIME += RANDOM_PASS_TIME
					// var TIME_WEIGHT = float64(2) * 0.2 * float64(RANDOM_PASS_TIME) //충돌?
					// time.Sleep(time.Duration(RANDOM_PASS_TIME+int(TIME_WEIGHT)) * time.Millisecond)
					time.Sleep(time.Duration(RANDOM_PASS_TIME) * time.Millisecond)
					log.Printf("2대 합의과정인데, hv가 한 대라도 있을 경우 : %v\n", time.Duration(RANDOM_PASS_TIME)*time.Millisecond)
				} else {
					// 비전시스템
					time.Sleep(500 * time.Millisecond)
				}

				if VEHICLES[0] == VEHICLES[1] {
					VEHICLES = RemoveValue(VEHICLES, VEHICLES[1])
					PASS_COUNT++
				} else {
					VEHICLES = RemoveValue(VEHICLES, VEHICLES[1])
					PASS_COUNT++
					VEHICLES = RemoveValue(VEHICLES, VEHICLES[0])
					PASS_COUNT++
				}
				break
			}

			// 3대 이상일 때
			if TOTAL_VEHICLES >= 3 {
				PASS_COUNT = 0
				log.Printf("테스트1: %d \n", VEHICLES)
				var directions = [12]string{"Rs", "Rl", "Rr", "Ls", "Ll", "Lr", "Ds", "Dl", "Dr", "Us", "Ul", "Ur"}

				var DirectionMap map[int32]string
				var grpcServers []*grpc.Server
				var dataMu sync.Mutex
				var wg sync.WaitGroup

				wg.Add(len(VEHICLES))
				DirectionMap = make(map[int32]string)

				// 서버를 각각의 고루틴에서 실행
				for _, i := range VEHICLES {
					// vehicle_direction_list = append(vehicle_direction_list, directions[rand.Intn(len(directions))])
					DirectionMap[i] = directions[rand.Intn(len(directions))]
					log.Printf("directionMap test : %v\n", DirectionMap)

					// 고루틴 내에서 실행
					go func(index int32, direction string, number int32) {
						defer wg.Done()

						grpcServer, port := startServer(index, direction, number) // 랜덤 방향 같이 전달

						if grpcServer != nil {
							// defer grpcServer.GracefulStop()            // 서버 종료 처리
							log.Printf("서버 %d의 포트: %s\n", index, port) // 포트 정보 출력
						} else {
							log.Printf("gRPC 서버 생성 실패\n")
						}

						// 뮤텍스를 사용하여 슬라이스에 접근
						dataMu.Lock()
						grpcServers = append(grpcServers, grpcServer)
						dataMu.Unlock()
					}(i, DirectionMap[i], i) // 인덱스와 방향을 고루틴에 전달
				}

				wg.Wait()

				// grpcServers 상태 출력
				log.Printf("grpcServers : %v\n", grpcServers)

				// WaitGroup을 사용하여 모든 고루틴이 끝날 때까지 대기
				wg.Add(int(TOTAL_VEHICLES))

				// 모든 서버의 데이터를 저장할 맵
				serverData := make(map[int32]*pb.Vehicle)

				var done = false

				// 각 서버에 요청 보내기
				for _, i := range VEHICLES {
					if Contains(RandomByzantine, i) { // i가 RandomByzantine에 포함되어 있는지 확인
						log.Printf("비잔틴 장애 차량 number = %d\n", i)
					}

					go func(i int32) {
						defer wg.Done()
						for _, j := range VEHICLES {
							if i == j {
								continue // 자기 자신에게는 요청을 보내지 않음
							}
							if Contains(RandomByzantine, j) {
								continue // 비잔틴 장애로부터 요청 못받음을 그냥 continue로 표현
							}
							if Contains(RandomByzantine, i) {
								continue // 비잔틴 장애로부터 요청 못받음을 그냥 continue로 표현
							}

							wg.Add(1) // 새로운 고루틴을 시작하므로 WaitGroup 카운터를 증가
							go func(j int32) {
								defer wg.Done()

								if done { // 종료 플래그가 설정된 경우
									return
								}

								addr := fmt.Sprintf("localhost:%d", GO_SERVER_PORT+int(j))
								client, conn, ctx, cancel, err := rpcConnectTo(addr)
								if err != nil {
									log.Printf("connect error: %v", err)
								}
								defer conn.Close()
								defer cancel()

								log.Printf("VEHICLES : %v\n", VEHICLES)

								r, err := client.ReceiveRequest(
									ctx,
									&pb.Request{
										Vehicle: &pb.Vehicle{
											Number:    i,
											Address:   i,
											Direction: DirectionMap[i],
										},
										Port:          fmt.Sprintf("%d", GO_SERVER_PORT+int(i)),
										TotalVehicles: TOTAL_VEHICLES,
									},
								)
								if err != nil {
									log.Printf("grpc call error: %v", err)
								}

								if r == nil {
									log.Printf("received nil response from server at %s", addr)
									return
								}

								// 종료 플래그가 설정된 경우
								if done {
									return
								}

								if r.Status == "acknowledged" {
									dataMu.Lock()
									defer dataMu.Unlock()
									if done { // 종료 플래그가 설정된 경우
										return
									}

									vehicle, exists := serverData[i]

									// 기존 vehicle이 있으면 가져오고 없으면 새로운 vehicle 생성
									if !exists {
										vehicle = &pb.Vehicle{
											Number:       i,
											Address:      i,
											ReceiveVotes: 0,
											ElectionVote: 0,
										}
									}

									// Covehicle 처리 (DirectionStatus가 "True"인 경우에만)
									if r.DirectionStatus == "True" {
										var covehicles []*pb.Vehicle
										covehicles = vehicle.Covehicle // 기존의 covehicles 가져오기
										var covehicleCheck = false

										// covehicles 리스트를 순회하면서 같은 방향의 차량을 찾고 추가
										for i := int32(0); i < int32(len(covehicles)); i++ {
											if direction.DirectionBoolean(covehicles[i].Direction, r.ResponseDirection) {
												covehicleCheck = true

												if covehicles[i].Covehicle == nil {
													covehicles[i].Covehicle = []*pb.Vehicle{}
												}
												covehicles[i].Covehicle = append(covehicles[i].Covehicle, r.Vehicle)
												log.Printf("\n")
												log.Printf("linked_covehicles test : %v\n", covehicles[i].Covehicle)
												log.Printf("r.Vehicle test : %v\n", r.Vehicle)
												log.Printf("\n")
											}
										}

										if covehicleCheck == false {
											covehicles = append(covehicles, r.Vehicle)
											vehicle.Covehicle = covehicles
										}

										log.Printf("\n")
										log.Printf("covehicles updated: %v\n", vehicle.Covehicle)
										log.Printf("\n")
									}

									// 요청 시간 저장
									vehicle.ElectionTime = pbtimestamp.Now()

									// 비동기적으로 UpdateVoteCount 호출
									go func(vehicle *pb.Vehicle) {
										if done { // 종료 플래그가 설정된 경우
											return
										}

										// 투표 수 증가
										vehicle.ReceiveVotes++
										addr := fmt.Sprintf("localhost:%d", GO_SERVER_PORT+int(vehicle.Address)) // 자신의 서버 주소
										client, conn, ctx, cancel, err := rpcConnectTo(addr)
										if err != nil {
											log.Printf("connect error: %v", err)
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
									}(vehicle) // 고루틴에 vehicle 전달

									// vehicle을 업데이트
									serverData[i] = vehicle

									// 만약 Receive Votes >= Quorum이면 자신을 제외한 차량에 LeaderElection 요청을 보내기
									if vehicle.ReceiveVotes >= QUORUM-1 {
										log.Printf("차량 %d의 ReceiveVotes가 Quorum에 도달했습니다. LeaderElection 요청을 보냅니다.\n", vehicle.Address)

										var wg sync.WaitGroup
										// 고루틴으로 LeaderElection 요청 전송
										for _, k := range VEHICLES {
											if k == i {
												continue // 자기 자신에게는 요청을 보내지 않음
											}
											if Contains(RandomByzantine, k) {
												continue // 비잔틴 장애로부터 요청 못받음을 그냥 continue로 표현
											}

											wg.Add(1) // 고루틴 수 증가
											go func(k int32) {
												defer wg.Done() // 고루틴 완료 시 호출
												if done {       // 종료 플래그가 설정된 경우
													return
												}
												addr := fmt.Sprintf("localhost:%d", GO_SERVER_PORT+int(k))
												client, conn, ctx, cancel, err := rpcConnectTo(addr)
												if err != nil {
													log.Printf("connect error to %s: %v\n", addr, err)
													return
												}
												defer conn.Close()
												defer cancel()

												// LeaderElection 요청 전송
												r, err = client.LeaderElection(ctx, &pb.Request{
													Vehicle: vehicle,
												})
												if err != nil {
													log.Printf("LeaderElection 요청 실패: %v\n", err)
												} else {
													log.Printf("차량 %d에 LeaderElection 요청을 성공적으로 보냈습니다.\n", k)
												}

												if r == nil {
													log.Printf("received nil response from server at %s", addr)
													return
												}

												// 종료 플래그가 설정된 경우
												if done {
													return
												}

												if r.Status == "acknowledged" {
													vehicle.ElectionVote++
													log.Printf("현재 동의 수 : %d\n", vehicle.ElectionVote)
													if vehicle.ElectionVote >= QUORUM-1 {
														log.Printf("vehicle이 비어있는지 확인 : %v\n", vehicle)
														removeVehiclesIfQuorumReached(vehicle)
														log.Printf("test2\n")
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
									log.Printf("test number: %v, test address: %d\n", VEHICLES, i)
									log.Printf("pass count: %d\n", PASS_COUNT)
									if !exists {
										vehicle = &pb.Vehicle{
											Number:       i,
											Address:      i,
											ReceiveVotes: 0,
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

				// endTime := time.Now()
				// duration := endTime.Sub(startTime)
				// TOTOA_TIME += float64(duration)

				// 여기서 grpcServers를 사용하여 모든 서버 종료 가능
				for _, server := range grpcServers {
					if server != nil {
						server.GracefulStop()
					} else {
						log.Println("A grpcServer is nil, skipping graceful stop")
					}
				}

				dataMu.Unlock()

				log.Printf("\n")
				// log.Printf("합의 과정에 걸린 시간: %v\n", duration)

				log.Printf("남은 차량: %d \n", VEHICLES)

				// 랜덤 무응답 차량 통과
				var NUMBER_OF_PASS_STOP_VEHICLES = rand.Intn(len(RandomByzantine) + 1)
				log.Printf("RandomByzantine 길이 : %d \n", len(RandomByzantine))

				for i := 1; i <= NUMBER_OF_PASS_STOP_VEHICLES; i++ {
					var PASS_STOP_VEHICLES = RandomByzantine[rand.Intn(len(RandomByzantine))]
					log.Printf("무응답차량 %d가 통과하였습니다.\n", PASS_STOP_VEHICLES)
					VEHICLES = RemoveValue(VEHICLES, int32(PASS_STOP_VEHICLES))
					PASS_COUNT++
					RandomByzantine = RemoveValue(RandomByzantine, int32(PASS_STOP_VEHICLES))
				}

				if NUMBER_OF_PASS_STOP_VEHICLES >= 1 {
					// var RANDOM_PASS_TIME = rand.Intn(3000-100) + 100
					var RANDOM_PASS_TIME = 3000
					STOP_VEHICLES_PASS_TIME += RANDOM_PASS_TIME
					// var TIME_WEIGHT = float64(NUMBER_OF_PASS_STOP_VEHICLES) * 0.2 * float64(RANDOM_PASS_TIME)
					// log.Printf("time weight test : %v\n", TIME_WEIGHT)
					log.Printf("time weight test2 : %v\n", time.Duration(RANDOM_PASS_TIME)*time.Millisecond)
					time.Sleep(time.Duration(RANDOM_PASS_TIME) * time.Millisecond)
				}

				log.Printf("무응답차량 통과하고 남은 차량: %d \n", VEHICLES)

				// GO_SERVER_PORT += int(TOTAL_VEHICLES)
				log.Printf("GO_SERVER_PORT 변경 : %v\n", GO_SERVER_PORT)
				// log.Printf("합의 걸리는 최종시간 : %v\n", TOTOA_TIME)
				log.Printf("\n")
			}
		}
		totalVehicles = difference(totalVehicles, selectedVehicles)
		log.Printf("결과 차량: %v\n", totalVehicles)

	}
	totalEndTime := time.Now()
	duration := totalEndTime.Sub(totalStartTime)
	log.Printf("합의 걸리는 최종시간 : %v\n", duration)
}
