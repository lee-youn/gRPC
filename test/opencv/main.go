package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/otiai10/gosseract/v2"

	"gocv.io/x/gocv"
	// 패키지 경로 확인
)

const GO_SERVER_PORT = 5005

var VEHICLES []int32

////////////////////////////////////////
/////////////////SERVER/////////////////
////////////////////////////////////////

// // 서버 구조체 정의
// type server struct {
// 	pb.UnimplementedVehicleServiceServer
// 	Port    string
// 	Vehicle *pb.Vehicle // 포트당 하나의 차량 정보를 저장
// 	mu      sync.Mutex
// }

// // 서버 인스턴스 생성 및 초기화 후 실행
// func startServer(address int32, direction string) {
// 	// 서버 인스턴스 생성 및 초기화
// 	s := &server{
// 		Port:    fmt.Sprintf("%d", GO_SERVER_PORT+address),           // 포트 번호를 문자열로 변환
// 		Vehicle: &pb.Vehicle{Address: address, Direction: direction}, // 기본 차량 정보로 초기화
// 	}

// 	// TCP 리스너 생성
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GO_SERVER_PORT+address))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	fmt.Printf("Server is listening on port %d\n", GO_SERVER_PORT+address)

// 	// gRPC 서버 생성
// 	grpcServer := grpc.NewServer()
// 	pb.RegisterVehicleServiceServer(grpcServer, s)

// 	// 서버 시작 (차단 함수)
// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}

// }

// func (s *server) ReceiveRequest(ctx context.Context, req *pb.Request) (*pb.Response, error) {
// 	fmt.Printf("Port %s made a request to port %s.\n", req.Port, s.Port)

// 	// 동시에 접근 못하게 함(즉, 동시에 request가 도착하는 일 없음)
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	fmt.Printf("Port %s direction : %s\n", s.Port, s.Vehicle.Direction)
// 	fmt.Printf("Port %s direction : %s\n", req.Port, req.Vehicle.Direction)

// 	if s.Vehicle.SendVotes == 0 {
// 		// 투표 응답을 처리하는 경우
// 		fmt.Printf("Vehicle %d received request from port %s: %v\n", s.Vehicle.Address, req.Port, req.Vehicle)
// 		s.Vehicle.SendVotes = 1 // 투표 완료로 설정

// 		response := &pb.Response{
// 			Message:           fmt.Sprintf("Vote registered from port %s to port %s", req.Port, s.Port),
// 			Status:            "acknowledged",
// 			ResponseDirection: s.Vehicle.Direction,
// 			Vehicle:           s.Vehicle, // 현재 서버의 차량 정보를 포함
// 		}

// 		return response, nil
// 	} else {
// 		// 이미 투표한 경우
// 		fmt.Printf("Vehicle %d has already responded to a request.\n", s.Vehicle.Address)

// 		// 응답 메시지 작성
// 		response := &pb.Response{
// 			Message: fmt.Sprintf("Vehicle %d has already voted", s.Vehicle.Address),
// 			Status:  "ignored",
// 			Vehicle: s.Vehicle, // 현재 서버의 차량 정보를 포함
// 		}

// 		return response, nil
// 	}
// }

// func (s *server) RandomAgreement(ctx context.Context, req *pb.Request) (*pb.Response, error) {
// 	fmt.Printf("Port %s made a request to port %s.\n", req.Port, s.Port)

// 	s.Vehicle.RandomNumber = req.RandomNumber

// 	var response *pb.Response

// 	if req.Vehicle.RandomNumber == s.Vehicle.RandomNumber {
// 		fmt.Printf("reqeust randomnumber: %d\n", req.Vehicle.RandomNumber)
// 		fmt.Printf("server randomnumber: %d\n", s.Vehicle.RandomNumber)
// 		fmt.Printf("Cannot reach an agreement.\n")
// 		response = &pb.Response{
// 			Message: fmt.Sprintf("Cannot reach an agreement"),
// 			Status:  "draw",
// 		}
// 	} else if req.Vehicle.RandomNumber > s.Vehicle.RandomNumber {
// 		fmt.Printf("reqeust randomnumber: %d\n", req.Vehicle.RandomNumber)
// 		fmt.Printf("server randomnumber: %d\n", s.Vehicle.RandomNumber)
// 		fmt.Printf("Vehicle %d may pass first.\n", req.Vehicle.Address)
// 		response = &pb.Response{
// 			Message: fmt.Sprintf("Vehicle %d may pass first.", req.Vehicle.Address),
// 			Status:  "pass",
// 		}
// 	} else {
// 		fmt.Printf("reqeust randomnumber: %d\n", req.Vehicle.RandomNumber)
// 		fmt.Printf("server randomnumber: %d\n", s.Vehicle.RandomNumber)
// 		fmt.Printf("Vehicle %d has to wait.\n", req.Vehicle.Address)
// 		response = &pb.Response{
// 			Message: fmt.Sprintf("Vehicle %d has to wait.", req.Vehicle.Address),
// 			Status:  "wait",
// 		}
// 	}

// 	return response, nil
// }

// ////////////////////////////////////////
// /////////////////CLIENT/////////////////
// ////////////////////////////////////////

// func rpcConnectTo(ip string) (pb.VehicleServiceClient, *grpc.ClientConn, context.Context, context.CancelFunc, error) {
// 	conn, err := grpc.Dial(
// 		ip,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		return nil, nil, nil, nil, fmt.Errorf("did not connect: %v", err)
// 	}

// 	c := pb.NewVehicleServiceClient(conn)
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
// 	return c, conn, ctx, cancel, nil
// }

// func RemoveValue(slice []int32, valueToRemove int32) []int32 {
// 	result := []int32{}
// 	for _, value := range slice {
// 		if value != valueToRemove {
// 			result = append(result, value)
// 		}
// 	}
// 	return result
// }

// var TOTAL_VEHICLS int32

//	func Contains(slice []int32, item int32) bool {
//		for _, v := range slice {
//			if v == item {
//				return true
//			}
//		}
//		return false
//	}
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

func main() {
	startTime := time.Now()

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	model := "models/yolov4.weights"
	config := "models/yolov4.cfg"

	net := gocv.ReadNet(model, config)
	if net.Empty() {
		log.Fatalf("Error reading network model from: %v %v", model, config)
	}

	imagePath := "images/N2.jpeg"
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		log.Fatalf("이미지를 읽지 못했습니다: %s", imagePath)
	}
	defer img.Close()

	fmt.Printf("이미지 로드 성공 \n")

	// Load XML file with license plate coordinates
	xmlFile, err := os.Open("images/N2.xml")
	if err != nil {
		log.Fatalf("Failed to open XML file: %v\n", err)
	}
	defer xmlFile.Close()

	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatalf("Failed to read XML file: %v\n", err)
	}

	var annotation Annotation
	err = xml.Unmarshal(xmlData, &annotation)
	if err != nil {
		log.Fatalf("XML 파싱 실패: %v\n", err)
	}

	// Tesseract 클라이언트 초기화
	client := gosseract.NewClient()
	defer client.Close()

	// Iterate through the plates found in the XML
	// XML에서 발견된 객체를 순회
	for _, obj := range annotation.Objects {
		if obj.Name == "num_plate" {
			left := obj.BndBox.Xmin
			top := obj.BndBox.Ymin
			right := obj.BndBox.Xmax
			bottom := obj.BndBox.Ymax

			// 바운딩 박스 그리기
			gocv.Rectangle(&img, image.Rect(left, top, right, bottom), color.RGBA{0, 255, 0, 0}, 2)

			// 이미지에서 번호판 잘라내기
			croppedPlate := img.Region(image.Rect(left, top, right, bottom))
			if croppedPlate.Empty() {
				fmt.Println("잘린 번호판 이미지가 비어 있습니다.")
			} else {
				fmt.Printf("잘린 번호판 이미지 존재: left=%d, top=%d, right=%d, bottom=%d\n", left, top, right, bottom)

				// 잘린 번호판 이미지를 바이트로 인코딩
				plateBytes, err := gocv.IMEncode(gocv.JPEGFileExt, croppedPlate)
				if err != nil {
					log.Fatalf("번호판 이미지 인코딩 실패: %v\n", err)
				}

				// Tesseract OCR로 텍스트 인식
				client.SetImageFromBytes(plateBytes.GetBytes())
				text, err := client.Text()
				if err != nil {
					log.Fatalf("텍스트 인식 실패: %v\n", err)
				}

				// 인식된 텍스트 출력
				fmt.Printf("인식된 번호판: %s\n", text)

				// 이미지에 인식된 텍스트 표시
				gocv.PutText(&img, text, image.Pt(left, top-10), gocv.FontHersheyPlain, 2.0, color.RGBA{255, 0, 0, 0}, 2)
			}
		}
	}

	// 결과 이미지 표시
	window := gocv.NewWindow("Detected License Plate")
	defer window.Close()

	window.IMShow(img)

	elapsedTime := time.Since(startTime)
	fmt.Printf("프로그램 실행 시간: %s\n", elapsedTime)
	return
}
