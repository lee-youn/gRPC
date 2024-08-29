// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: goServer.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: vehicle.proto

package vehicleServer

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Vehicle 메시지 정의
type Vehicle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number       string               `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`                                  // 차량 번호
	Direction    string               `protobuf:"bytes,2,opt,name=direction,proto3" json:"direction,omitempty"`                            // 방향
	Address      int32                `protobuf:"varint,3,opt,name=address,proto3" json:"address,omitempty"`                               // 주소
	SendVotes    int32                `protobuf:"varint,4,opt,name=send_votes,json=sendVotes,proto3" json:"send_votes,omitempty"`          // 보내는 투표 수
	ReceiveVotes int32                `protobuf:"varint,5,opt,name=receive_votes,json=receiveVotes,proto3" json:"receive_votes,omitempty"` // 받은 투표 수
	Covehicle    []*ConcurrentVehicle `protobuf:"bytes,6,rep,name=covehicle,proto3" json:"covehicle,omitempty"`                            // 같은 방향 차량
	RandomNumber int32                `protobuf:"varint,7,opt,name=random_number,json=randomNumber,proto3" json:"random_number,omitempty"`
}

func (x *Vehicle) Reset() {
	*x = Vehicle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vehicle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vehicle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vehicle) ProtoMessage() {}

func (x *Vehicle) ProtoReflect() protoreflect.Message {
	mi := &file_vehicle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vehicle.ProtoReflect.Descriptor instead.
func (*Vehicle) Descriptor() ([]byte, []int) {
	return file_vehicle_proto_rawDescGZIP(), []int{0}
}

func (x *Vehicle) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *Vehicle) GetDirection() string {
	if x != nil {
		return x.Direction
	}
	return ""
}

func (x *Vehicle) GetAddress() int32 {
	if x != nil {
		return x.Address
	}
	return 0
}

func (x *Vehicle) GetSendVotes() int32 {
	if x != nil {
		return x.SendVotes
	}
	return 0
}

func (x *Vehicle) GetReceiveVotes() int32 {
	if x != nil {
		return x.ReceiveVotes
	}
	return 0
}

func (x *Vehicle) GetCovehicle() []*ConcurrentVehicle {
	if x != nil {
		return x.Covehicle
	}
	return nil
}

func (x *Vehicle) GetRandomNumber() int32 {
	if x != nil {
		return x.RandomNumber
	}
	return 0
}

// Request 메시지 정의
type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vehicle       *Vehicle `protobuf:"bytes,1,opt,name=vehicle,proto3" json:"vehicle,omitempty"` // 차량 정보
	Port          string   `protobuf:"bytes,2,opt,name=Port,proto3" json:"Port,omitempty"`
	TotalVehicles int32    `protobuf:"varint,3,opt,name=total_vehicles,json=totalVehicles,proto3" json:"total_vehicles,omitempty"` // 총 차량 수
	RandomNumber  int32    `protobuf:"varint,4,opt,name=RandomNumber,proto3" json:"RandomNumber,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vehicle_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_vehicle_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_vehicle_proto_rawDescGZIP(), []int{1}
}

func (x *Request) GetVehicle() *Vehicle {
	if x != nil {
		return x.Vehicle
	}
	return nil
}

func (x *Request) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *Request) GetTotalVehicles() int32 {
	if x != nil {
		return x.TotalVehicles
	}
	return 0
}

func (x *Request) GetRandomNumber() int32 {
	if x != nil {
		return x.RandomNumber
	}
	return 0
}

// Response 메시지 정의
type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message         string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Status          string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	DirectionStatus string   `protobuf:"bytes,3,opt,name=direction_status,json=directionStatus,proto3" json:"direction_status,omitempty"`
	Vehicle         *Vehicle `protobuf:"bytes,4,opt,name=vehicle,proto3" json:"vehicle,omitempty"` // Vehicle 정보를 포함
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vehicle_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_vehicle_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_vehicle_proto_rawDescGZIP(), []int{2}
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Response) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Response) GetDirectionStatus() string {
	if x != nil {
		return x.DirectionStatus
	}
	return ""
}

func (x *Response) GetVehicle() *Vehicle {
	if x != nil {
		return x.Vehicle
	}
	return nil
}

// ConcurrentVehicle 메시지 정의
type ConcurrentVehicle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vehicle *Vehicle           `protobuf:"bytes,1,opt,name=vehicle,proto3" json:"vehicle,omitempty"` // 차량 정보
	Next    *ConcurrentVehicle `protobuf:"bytes,2,opt,name=next,proto3" json:"next,omitempty"`       // 다음 ConcurrentVehicle 포인터
}

func (x *ConcurrentVehicle) Reset() {
	*x = ConcurrentVehicle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vehicle_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConcurrentVehicle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConcurrentVehicle) ProtoMessage() {}

func (x *ConcurrentVehicle) ProtoReflect() protoreflect.Message {
	mi := &file_vehicle_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConcurrentVehicle.ProtoReflect.Descriptor instead.
func (*ConcurrentVehicle) Descriptor() ([]byte, []int) {
	return file_vehicle_proto_rawDescGZIP(), []int{3}
}

func (x *ConcurrentVehicle) GetVehicle() *Vehicle {
	if x != nil {
		return x.Vehicle
	}
	return nil
}

func (x *ConcurrentVehicle) GetNext() *ConcurrentVehicle {
	if x != nil {
		return x.Next
	}
	return nil
}

// VehicleRPC 메시지 정의
type VehicleRPC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address               int32              `protobuf:"varint,1,opt,name=address,proto3" json:"address,omitempty"`                                                                                           // 주소
	Vehicles              map[int32]*Vehicle `protobuf:"bytes,2,rep,name=vehicles,proto3" json:"vehicles,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // 차량 맵
	TotalVehicles         int32              `protobuf:"varint,3,opt,name=total_vehicles,json=totalVehicles,proto3" json:"total_vehicles,omitempty"`                                                          // 총 차량 수
	VoteCount             int32              `protobuf:"varint,4,opt,name=vote_count,json=voteCount,proto3" json:"vote_count,omitempty"`                                                                      // 투표 수
	ConcurrentVehicleList *ConcurrentVehicle `protobuf:"bytes,5,opt,name=concurrent_vehicle_list,json=concurrentVehicleList,proto3" json:"concurrent_vehicle_list,omitempty"`                                 // ConcurrentVehicle 리스트
}

func (x *VehicleRPC) Reset() {
	*x = VehicleRPC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vehicle_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VehicleRPC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VehicleRPC) ProtoMessage() {}

func (x *VehicleRPC) ProtoReflect() protoreflect.Message {
	mi := &file_vehicle_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VehicleRPC.ProtoReflect.Descriptor instead.
func (*VehicleRPC) Descriptor() ([]byte, []int) {
	return file_vehicle_proto_rawDescGZIP(), []int{4}
}

func (x *VehicleRPC) GetAddress() int32 {
	if x != nil {
		return x.Address
	}
	return 0
}

func (x *VehicleRPC) GetVehicles() map[int32]*Vehicle {
	if x != nil {
		return x.Vehicles
	}
	return nil
}

func (x *VehicleRPC) GetTotalVehicles() int32 {
	if x != nil {
		return x.TotalVehicles
	}
	return 0
}

func (x *VehicleRPC) GetVoteCount() int32 {
	if x != nil {
		return x.VoteCount
	}
	return 0
}

func (x *VehicleRPC) GetConcurrentVehicleList() *ConcurrentVehicle {
	if x != nil {
		return x.ConcurrentVehicleList
	}
	return nil
}

var File_vehicle_proto protoreflect.FileDescriptor

var file_vehicle_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0d, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x22, 0x82,
	0x02, 0x0a, 0x07, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65,
	0x6e, 0x64, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x73, 0x65, 0x6e, 0x64, 0x56, 0x6f, 0x74, 0x65, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x56, 0x6f, 0x74, 0x65, 0x73, 0x12, 0x3e,
	0x0a, 0x09, 0x63, 0x6f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x20, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x52, 0x09, 0x63, 0x6f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x23,
	0x0a, 0x0d, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x22, 0x9a, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x30, 0x0a, 0x07, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x07, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x76,
	0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x12, 0x22, 0x0a, 0x0c,
	0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0c, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x22, 0x99, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x29, 0x0a, 0x10, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x30, 0x0a, 0x07, 0x76, 0x65,
	0x68, 0x69, 0x63, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x76, 0x65,
	0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x56, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x52, 0x07, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x22, 0x7b, 0x0a, 0x11,
	0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c,
	0x65, 0x12, 0x30, 0x0a, 0x07, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x07, 0x76, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x12, 0x34, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x20, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x22, 0xe0, 0x02, 0x0a, 0x0a, 0x56, 0x65,
	0x68, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x50, 0x43, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x43, 0x0a, 0x08, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x50, 0x43, 0x2e,
	0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x76,
	0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x5f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x12, 0x1d,
	0x0a, 0x0a, 0x76, 0x6f, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x76, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x58, 0x0a,
	0x17, 0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43,
	0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65,
	0x52, 0x15, 0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x53, 0x0a, 0x0d, 0x56, 0x65, 0x68, 0x69, 0x63,
	0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x76, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c,
	0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x97, 0x01, 0x0a,
	0x0e, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x41, 0x0a, 0x0e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x76, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x42, 0x0a, 0x0f, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x41, 0x67, 0x72, 0x65,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x2e, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x3b, 0x76, 0x65, 0x68, 0x69,
	0x63, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_vehicle_proto_rawDescOnce sync.Once
	file_vehicle_proto_rawDescData = file_vehicle_proto_rawDesc
)

func file_vehicle_proto_rawDescGZIP() []byte {
	file_vehicle_proto_rawDescOnce.Do(func() {
		file_vehicle_proto_rawDescData = protoimpl.X.CompressGZIP(file_vehicle_proto_rawDescData)
	})
	return file_vehicle_proto_rawDescData
}

var file_vehicle_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_vehicle_proto_goTypes = []any{
	(*Vehicle)(nil),           // 0: vehicleServer.Vehicle
	(*Request)(nil),           // 1: vehicleServer.Request
	(*Response)(nil),          // 2: vehicleServer.Response
	(*ConcurrentVehicle)(nil), // 3: vehicleServer.ConcurrentVehicle
	(*VehicleRPC)(nil),        // 4: vehicleServer.VehicleRPC
	nil,                       // 5: vehicleServer.VehicleRPC.VehiclesEntry
}
var file_vehicle_proto_depIdxs = []int32{
	3,  // 0: vehicleServer.Vehicle.covehicle:type_name -> vehicleServer.ConcurrentVehicle
	0,  // 1: vehicleServer.Request.vehicle:type_name -> vehicleServer.Vehicle
	0,  // 2: vehicleServer.Response.vehicle:type_name -> vehicleServer.Vehicle
	0,  // 3: vehicleServer.ConcurrentVehicle.vehicle:type_name -> vehicleServer.Vehicle
	3,  // 4: vehicleServer.ConcurrentVehicle.next:type_name -> vehicleServer.ConcurrentVehicle
	5,  // 5: vehicleServer.VehicleRPC.vehicles:type_name -> vehicleServer.VehicleRPC.VehiclesEntry
	3,  // 6: vehicleServer.VehicleRPC.concurrent_vehicle_list:type_name -> vehicleServer.ConcurrentVehicle
	0,  // 7: vehicleServer.VehicleRPC.VehiclesEntry.value:type_name -> vehicleServer.Vehicle
	1,  // 8: vehicleServer.VehicleService.ReceiveRequest:input_type -> vehicleServer.Request
	1,  // 9: vehicleServer.VehicleService.RandomAgreement:input_type -> vehicleServer.Request
	2,  // 10: vehicleServer.VehicleService.ReceiveRequest:output_type -> vehicleServer.Response
	2,  // 11: vehicleServer.VehicleService.RandomAgreement:output_type -> vehicleServer.Response
	10, // [10:12] is the sub-list for method output_type
	8,  // [8:10] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_vehicle_proto_init() }
func file_vehicle_proto_init() {
	if File_vehicle_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vehicle_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Vehicle); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vehicle_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vehicle_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vehicle_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ConcurrentVehicle); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vehicle_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*VehicleRPC); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vehicle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vehicle_proto_goTypes,
		DependencyIndexes: file_vehicle_proto_depIdxs,
		MessageInfos:      file_vehicle_proto_msgTypes,
	}.Build()
	File_vehicle_proto = out.File
	file_vehicle_proto_rawDesc = nil
	file_vehicle_proto_goTypes = nil
	file_vehicle_proto_depIdxs = nil
}
