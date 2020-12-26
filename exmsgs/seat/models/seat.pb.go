// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/seat/models/seat.proto

package seat_models // import "simple-restful/exmsgs/seat/models"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Seat struct {
	Row                  int64    `protobuf:"varint,1,opt,name=row" json:"row,omitempty"`
	Col                  int64    `protobuf:"varint,2,opt,name=col" json:"col,omitempty"`
	Status               string   `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Seat) Reset()         { *m = Seat{} }
func (m *Seat) String() string { return proto.CompactTextString(m) }
func (*Seat) ProtoMessage()    {}
func (*Seat) Descriptor() ([]byte, []int) {
	return fileDescriptor_seat_e5e600e58dce6789, []int{0}
}
func (m *Seat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Seat.Unmarshal(m, b)
}
func (m *Seat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Seat.Marshal(b, m, deterministic)
}
func (dst *Seat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Seat.Merge(dst, src)
}
func (m *Seat) XXX_Size() int {
	return xxx_messageInfo_Seat.Size(m)
}
func (m *Seat) XXX_DiscardUnknown() {
	xxx_messageInfo_Seat.DiscardUnknown(m)
}

var xxx_messageInfo_Seat proto.InternalMessageInfo

func (m *Seat) GetRow() int64 {
	if m != nil {
		return m.Row
	}
	return 0
}

func (m *Seat) GetCol() int64 {
	if m != nil {
		return m.Col
	}
	return 0
}

func (m *Seat) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*Seat)(nil), "protos.seat.models.Seat")
}

func init() {
	proto.RegisterFile("protos/seat/models/seat.proto", fileDescriptor_seat_e5e600e58dce6789)
}

var fileDescriptor_seat_e5e600e58dce6789 = []byte{
	// 148 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0x2f, 0xd6, 0x2f, 0x4e, 0x4d, 0x2c, 0xd1, 0xcf, 0xcd, 0x4f, 0x49, 0xcd, 0x81, 0xb0, 0xf5,
	0xc0, 0xe2, 0x42, 0x42, 0x10, 0x69, 0x3d, 0xb0, 0x10, 0x44, 0x5a, 0xc9, 0x89, 0x8b, 0x25, 0x38,
	0x35, 0xb1, 0x44, 0x48, 0x80, 0x8b, 0xb9, 0x28, 0xbf, 0x5c, 0x82, 0x51, 0x81, 0x51, 0x83, 0x39,
	0x08, 0xc4, 0x04, 0x89, 0x24, 0xe7, 0xe7, 0x48, 0x30, 0x41, 0x44, 0x92, 0xf3, 0x73, 0x84, 0xc4,
	0xb8, 0xd8, 0x8a, 0x4b, 0x12, 0x4b, 0x4a, 0x8b, 0x25, 0x98, 0x15, 0x18, 0x35, 0x38, 0x83, 0xa0,
	0x3c, 0x27, 0xfd, 0x28, 0xdd, 0xe2, 0xcc, 0xdc, 0x82, 0x9c, 0x54, 0xdd, 0xa2, 0xd4, 0xe2, 0x92,
	0xb4, 0xd2, 0x1c, 0xfd, 0xd4, 0x8a, 0xdc, 0xe2, 0x74, 0x14, 0x77, 0x58, 0x83, 0xd8, 0xf1, 0x10,
	0x76, 0x12, 0x1b, 0xd8, 0x21, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb0, 0xc9, 0x28, 0xdd,
	0xb0, 0x00, 0x00, 0x00,
}