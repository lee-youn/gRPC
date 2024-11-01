package licenseplate

import (
	"encoding/xml"
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

func LicensePlate(value string) string {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// model := "../models/yolov4.weights"
	// config := "../models/yolov4.cfg"

	// net := gocv.ReadNet(model, config)
	// if net.Empty() {
	// 	log.Printf("Error reading network model from: %v %v", model, config)
	// }
	var CarNumber string

	imagePath := "images/N" + value + ".jpeg"
	log.Printfln("Image Path:", imagePath)
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	// img2 := "../images/N65.jpeg"
	if img.Empty() {
		log.Printf("이미지를 읽지 못했습니다: %s", imagePath)
	}
	defer img.Close()

	log.Printf("이미지 로드 성공 \n")

	// Load XML file with license plate coordinates
	xmlPath := "images/N" + value + ".xml"
	xmlFile, err := os.Open(xmlPath)
	if err != nil {
		log.Printf("Failed to open XML file: %v\n", err)
	}
	defer xmlFile.Close()

	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Printf("Failed to read XML file: %v\n", err)
	}

	var annotation Annotation
	err = xml.Unmarshal(xmlData, &annotation)
	log.Printf("annotation 성공\n")
	if err != nil {
		log.Printf("XML 파싱 실패: %v\n", err)
	}

	// Tesseract 클라이언트 초기화
	client := gosseract.NewClient()
	defer client.Close()

	// Iterate through the plates found in the XML
	// XML에서 발견된 객체를 순회
	for _, obj := range annotation.Objects {
		log.Printf("Object Name: %s\n", obj.Name)
		if obj.Name == "number_plate" {
			left := obj.BndBox.Xmin
			top := obj.BndBox.Ymin
			right := obj.BndBox.Xmax
			bottom := obj.BndBox.Ymax

			// 바운딩 박스 그리기
			gocv.Rectangle(&img, image.Rect(left, top, right, bottom), color.RGBA{0, 255, 0, 0}, 2)

			// 이미지에서 번호판 잘라내기
			croppedPlate := img.Region(image.Rect(left, top, right, bottom))
			if croppedPlate.Empty() {
				log.Printfln("잘린 번호판 이미지가 비어 있습니다.")
			} else {
				log.Printf("잘린 번호판 이미지 존재: left=%d, top=%d, right=%d, bottom=%d\n", left, top, right, bottom)

				// 잘린 번호판 이미지를 바이트로 인코딩
				plateBytes, err := gocv.IMEncode(gocv.JPEGFileExt, croppedPlate)
				if err != nil {
					log.Printf("번호판 이미지 인코딩 실패: %v\n", err)
				}

				// log.Printfln(plateBytes)

				// Tesseract OCR로 텍스트 인식
				client.SetImageFromBytes(plateBytes.GetBytes())
				text, err := client.Text()
				if err != nil {
					log.Printf("텍스트 인식 실패: %v\n", err)
				}

				// 인식된 텍스트 출력
				log.Printf("인식된 번호판: %s\n", text)
				CarNumber = text

				// 이미지에 인식된 텍스트 표시
				gocv.PutText(&img, text, image.Pt(left, top-10), gocv.FontHersheyPlain, 2.0, color.RGBA{255, 0, 0, 0}, 2)
			}
		}
	}

	// // 결과 이미지 표시
	// window := gocv.NewWindow("Detected License Plate")
	// defer window.Close()

	// window.IMShow(img)
	log.Printf("car number : %s\n", CarNumber)

	return CarNumber
}

func LicensePlateComparison(vehicle1 string, vehicle2 string) bool {
	comparisonStartTime := time.Now()
	vehicle1_number := LicensePlate(vehicle1)
	vehicle2_number := LicensePlate(vehicle2)
	comparisonEndTime := time.Now()
	duration := comparisonEndTime.Sub(comparisonStartTime)
	log.Printf("vehicle1_number: %s\n", vehicle1_number)
	log.Printf("vehicle2_number: %s\n", vehicle2_number)
	log.Printf("duration time : %v\n", duration)

	return vehicle1_number >= vehicle2_number
}
