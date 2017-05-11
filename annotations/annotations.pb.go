// Code generated by protoc-gen-go.
// source: annotations.proto
// DO NOT EDIT!

/*
Package annotations is a generated protocol buffer package.

It is generated from these files:
	annotations.proto

It has these top-level messages:
	Interceptors
*/
package annotations

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Interceptors struct {
	Indexes          []string `protobuf:"bytes,1,rep,name=indexes" json:"indexes,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Interceptors) Reset()                    { *m = Interceptors{} }
func (m *Interceptors) String() string            { return proto.CompactTextString(m) }
func (*Interceptors) ProtoMessage()               {}
func (*Interceptors) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Interceptors) GetIndexes() []string {
	if m != nil {
		return m.Indexes
	}
	return nil
}

var E_PackageInterceptors = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FileOptions)(nil),
	ExtensionType: (*Interceptors)(nil),
	Field:         56780,
	Name:          "grpcmw.package_interceptors",
	Tag:           "bytes,56780,opt,name=package_interceptors,json=packageInterceptors",
	Filename:      "annotations.proto",
}

var E_ServiceInterceptors = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.ServiceOptions)(nil),
	ExtensionType: (*Interceptors)(nil),
	Field:         56781,
	Name:          "grpcmw.service_interceptors",
	Tag:           "bytes,56781,opt,name=service_interceptors,json=serviceInterceptors",
	Filename:      "annotations.proto",
}

var E_MethodInterceptors = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*Interceptors)(nil),
	Field:         56782,
	Name:          "grpcmw.method_interceptors",
	Tag:           "bytes,56782,opt,name=method_interceptors,json=methodInterceptors",
	Filename:      "annotations.proto",
}

func init() {
	proto.RegisterType((*Interceptors)(nil), "grpcmw.Interceptors")
	proto.RegisterExtension(E_PackageInterceptors)
	proto.RegisterExtension(E_ServiceInterceptors)
	proto.RegisterExtension(E_MethodInterceptors)
}

func init() { proto.RegisterFile("annotations.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 223 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xcc, 0xcb, 0xcb,
	0x2f, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4b,
	0x2f, 0x2a, 0x48, 0xce, 0x2d, 0x97, 0x52, 0x48, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x07, 0x8b,
	0x26, 0x95, 0xa6, 0xe9, 0xa7, 0xa4, 0x16, 0x27, 0x17, 0x65, 0x16, 0x94, 0xe4, 0x17, 0x41, 0x54,
	0x2a, 0x69, 0x70, 0xf1, 0x78, 0xe6, 0x95, 0xa4, 0x16, 0x25, 0xa7, 0x82, 0x04, 0x8b, 0x85, 0x24,
	0xb8, 0xd8, 0x33, 0xf3, 0x52, 0x52, 0x2b, 0x52, 0x8b, 0x25, 0x18, 0x15, 0x98, 0x35, 0x38, 0x83,
	0x60, 0x5c, 0xab, 0x74, 0x2e, 0x91, 0x82, 0xc4, 0xe4, 0xec, 0xc4, 0xf4, 0xd4, 0xf8, 0x4c, 0x64,
	0x1d, 0x32, 0x7a, 0x10, 0x4b, 0xf4, 0x60, 0x96, 0xe8, 0xb9, 0x65, 0xe6, 0xa4, 0xfa, 0x17, 0x80,
	0xdd, 0x23, 0x71, 0x66, 0x37, 0xb3, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x88, 0x1e, 0xc4, 0x49, 0x7a,
	0xc8, 0xb6, 0x05, 0x09, 0x43, 0x4d, 0x44, 0x16, 0xb4, 0xca, 0xe2, 0x12, 0x29, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x46, 0xb3, 0x48, 0x1e, 0xc3, 0xa2, 0x60, 0x88, 0x32, 0x98, 0x5d, 0x67, 0xf1, 0xdb,
	0x05, 0x35, 0x14, 0xc5, 0xae, 0x74, 0x2e, 0xe1, 0xdc, 0xd4, 0x92, 0x8c, 0xfc, 0x14, 0x54, 0xab,
	0xe4, 0x30, 0xac, 0xf2, 0x05, 0xab, 0x82, 0xd9, 0x74, 0x0e, 0xaf, 0x4d, 0x42, 0x10, 0x23, 0x91,
	0xc5, 0x00, 0x01, 0x00, 0x00, 0xff, 0xff, 0x68, 0x2e, 0x4d, 0x7d, 0xa5, 0x01, 0x00, 0x00,
}