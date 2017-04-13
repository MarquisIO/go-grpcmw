// Code generated by protoc-gen-go.
// source: interceptors.proto
// DO NOT EDIT!

/*
Package grpcmw is a generated protocol buffer package.

It is generated from these files:
	interceptors.proto

It has these top-level messages:
	Interceptors
*/
package grpcmw

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
	Symbols          []string `protobuf:"bytes,1,rep,name=symbols" json:"symbols,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *Interceptors) Reset()                    { *m = Interceptors{} }
func (m *Interceptors) String() string            { return proto.CompactTextString(m) }
func (*Interceptors) ProtoMessage()               {}
func (*Interceptors) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Interceptors) GetSymbols() []string {
	if m != nil {
		return m.Symbols
	}
	return nil
}

var E_PackageInterceptors = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FileOptions)(nil),
	ExtensionType: (*Interceptors)(nil),
	Field:         56780,
	Name:          "grpcmw.package_interceptors",
	Tag:           "bytes,56780,opt,name=package_interceptors,json=packageInterceptors",
	Filename:      "interceptors.proto",
}

var E_ServiceInterceptors = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.ServiceOptions)(nil),
	ExtensionType: (*Interceptors)(nil),
	Field:         56781,
	Name:          "grpcmw.service_interceptors",
	Tag:           "bytes,56781,opt,name=service_interceptors,json=serviceInterceptors",
	Filename:      "interceptors.proto",
}

var E_MethodInterceptors = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*Interceptors)(nil),
	Field:         56782,
	Name:          "grpcmw.method_interceptors",
	Tag:           "bytes,56782,opt,name=method_interceptors,json=methodInterceptors",
	Filename:      "interceptors.proto",
}

func init() {
	proto.RegisterType((*Interceptors)(nil), "grpcmw.Interceptors")
	proto.RegisterExtension(E_PackageInterceptors)
	proto.RegisterExtension(E_ServiceInterceptors)
	proto.RegisterExtension(E_MethodInterceptors)
}

func init() { proto.RegisterFile("interceptors.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 218 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0xca, 0xcc, 0x2b, 0x49,
	0x2d, 0x4a, 0x4e, 0x2d, 0x28, 0xc9, 0x2f, 0x2a, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x4b, 0x2f, 0x2a, 0x48, 0xce, 0x2d, 0x97, 0x52, 0x48, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x07,
	0x8b, 0x26, 0x95, 0xa6, 0xe9, 0xa7, 0xa4, 0x16, 0x27, 0x17, 0x65, 0x82, 0x94, 0x42, 0x54, 0x2a,
	0x69, 0x70, 0xf1, 0x78, 0x22, 0xe9, 0x17, 0x92, 0xe0, 0x62, 0x2f, 0xae, 0xcc, 0x4d, 0xca, 0xcf,
	0x29, 0x96, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x0c, 0x82, 0x71, 0xad, 0xd2, 0xb9, 0x44, 0x0a, 0x12,
	0x93, 0xb3, 0x13, 0xd3, 0x53, 0xe3, 0x91, 0x6d, 0x14, 0x92, 0xd1, 0x83, 0x58, 0xa2, 0x07, 0xb3,
	0x44, 0xcf, 0x2d, 0x33, 0x27, 0xd5, 0xbf, 0xa0, 0x24, 0x33, 0x3f, 0xaf, 0x58, 0xe2, 0xcc, 0x6e,
	0x66, 0x05, 0x46, 0x0d, 0x6e, 0x23, 0x11, 0x3d, 0x88, 0x93, 0xf4, 0x90, 0x6d, 0x0b, 0x12, 0x86,
	0x9a, 0x88, 0x2c, 0x68, 0x95, 0xc5, 0x25, 0x52, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x8c, 0x66, 0x91,
	0x3c, 0x86, 0x45, 0xc1, 0x10, 0x65, 0x30, 0xbb, 0xce, 0xe2, 0xb7, 0x0b, 0x6a, 0x28, 0x8a, 0x5d,
	0xe9, 0x5c, 0xc2, 0xb9, 0xa9, 0x25, 0x19, 0xf9, 0x29, 0xa8, 0x56, 0xc9, 0x61, 0x58, 0xe5, 0x0b,
	0x56, 0x05, 0xb3, 0xe9, 0x1c, 0x5e, 0x9b, 0x84, 0x20, 0x46, 0x22, 0x8b, 0x01, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x51, 0x2c, 0x64, 0xfa, 0xa6, 0x01, 0x00, 0x00,
}