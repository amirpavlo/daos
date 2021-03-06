// Code generated by protoc-gen-go. DO NOT EDIT.
// source: acl.proto

package acl

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RequestStatus int32

const (
	RequestStatus_SUCCESS               RequestStatus = 0
	RequestStatus_ERR_UNKNOWN           RequestStatus = -1
	RequestStatus_ERR_PERM_DENIED       RequestStatus = -2
	RequestStatus_ERR_INVALID_PERMS     RequestStatus = -3
	RequestStatus_ERR_INVALID_PRINCIPAL RequestStatus = -4
	RequestStatus_ERR_INVALID_UUID      RequestStatus = -5
)

var RequestStatus_name = map[int32]string{
	0:  "SUCCESS",
	-1: "ERR_UNKNOWN",
	-2: "ERR_PERM_DENIED",
	-3: "ERR_INVALID_PERMS",
	-4: "ERR_INVALID_PRINCIPAL",
	-5: "ERR_INVALID_UUID",
}
var RequestStatus_value = map[string]int32{
	"SUCCESS":               0,
	"ERR_UNKNOWN":           -1,
	"ERR_PERM_DENIED":       -2,
	"ERR_INVALID_PERMS":     -3,
	"ERR_INVALID_PRINCIPAL": -4,
	"ERR_INVALID_UUID":      -5,
}

func (x RequestStatus) String() string {
	return proto.EnumName(RequestStatus_name, int32(x))
}
func (RequestStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_acl_a67c5b5daa5f79f1, []int{0}
}

// Bits representing access permissions
type Permissions int32

const (
	Permissions_NO_ACCESS Permissions = 0
	Permissions_READ      Permissions = 1
	Permissions_WRITE     Permissions = 2
)

var Permissions_name = map[int32]string{
	0: "NO_ACCESS",
	1: "READ",
	2: "WRITE",
}
var Permissions_value = map[string]int32{
	"NO_ACCESS": 0,
	"READ":      1,
	"WRITE":     2,
}

func (x Permissions) String() string {
	return proto.EnumName(Permissions_name, int32(x))
}
func (Permissions) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_acl_a67c5b5daa5f79f1, []int{1}
}

// A given user/group may have multiple different types of entries
type EntryType int32

const (
	EntryType_ALLOW EntryType = 0
	EntryType_AUDIT EntryType = 1
	EntryType_ALARM EntryType = 2
)

var EntryType_name = map[int32]string{
	0: "ALLOW",
	1: "AUDIT",
	2: "ALARM",
}
var EntryType_value = map[string]int32{
	"ALLOW": 0,
	"AUDIT": 1,
	"ALARM": 2,
}

func (x EntryType) String() string {
	return proto.EnumName(EntryType_name, int32(x))
}
func (EntryType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_acl_a67c5b5daa5f79f1, []int{2}
}

// Bits representing flags on a given ACL entry
type Flags int32

const (
	Flags_NO_FLAGS       Flags = 0
	Flags_GROUP          Flags = 1
	Flags_ACCESS_SUCCESS Flags = 2
	Flags_ACCESS_FAILURE Flags = 4
	Flags_POOL_INHERIT   Flags = 8
)

var Flags_name = map[int32]string{
	0: "NO_FLAGS",
	1: "GROUP",
	2: "ACCESS_SUCCESS",
	4: "ACCESS_FAILURE",
	8: "POOL_INHERIT",
}
var Flags_value = map[string]int32{
	"NO_FLAGS":       0,
	"GROUP":          1,
	"ACCESS_SUCCESS": 2,
	"ACCESS_FAILURE": 4,
	"POOL_INHERIT":   8,
}

func (x Flags) String() string {
	return proto.EnumName(Flags_name, int32(x))
}
func (Flags) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_acl_a67c5b5daa5f79f1, []int{3}
}

type Response struct {
	Status               RequestStatus     `protobuf:"varint,1,opt,name=status,proto3,enum=acl.RequestStatus" json:"status,omitempty"`
	Permissions          *EntryPermissions `protobuf:"bytes,2,opt,name=permissions,proto3" json:"permissions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_acl_a67c5b5daa5f79f1, []int{0}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetStatus() RequestStatus {
	if m != nil {
		return m.Status
	}
	return RequestStatus_SUCCESS
}

func (m *Response) GetPermissions() *EntryPermissions {
	if m != nil {
		return m.Permissions
	}
	return nil
}

// Identifier for a specific Access Control Entry
type Entry struct {
	Type                 EntryType `protobuf:"varint,1,opt,name=type,proto3,enum=acl.EntryType" json:"type,omitempty"`
	Flags                uint32    `protobuf:"varint,2,opt,name=flags,proto3" json:"flags,omitempty"`
	Entity               string    `protobuf:"bytes,3,opt,name=entity,proto3" json:"entity,omitempty"`
	Identity             string    `protobuf:"bytes,4,opt,name=identity,proto3" json:"identity,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_acl_a67c5b5daa5f79f1, []int{1}
}
func (m *Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entry.Unmarshal(m, b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
}
func (dst *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(dst, src)
}
func (m *Entry) XXX_Size() int {
	return xxx_messageInfo_Entry.Size(m)
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetType() EntryType {
	if m != nil {
		return m.Type
	}
	return EntryType_ALLOW
}

func (m *Entry) GetFlags() uint32 {
	if m != nil {
		return m.Flags
	}
	return 0
}

func (m *Entry) GetEntity() string {
	if m != nil {
		return m.Entity
	}
	return ""
}

func (m *Entry) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

// Permissions for the given entry
type EntryPermissions struct {
	Entry                *Entry   `protobuf:"bytes,1,opt,name=entry,proto3" json:"entry,omitempty"`
	PermissionBits       uint64   `protobuf:"varint,2,opt,name=permission_bits,json=permissionBits,proto3" json:"permission_bits,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EntryPermissions) Reset()         { *m = EntryPermissions{} }
func (m *EntryPermissions) String() string { return proto.CompactTextString(m) }
func (*EntryPermissions) ProtoMessage()    {}
func (*EntryPermissions) Descriptor() ([]byte, []int) {
	return fileDescriptor_acl_a67c5b5daa5f79f1, []int{2}
}
func (m *EntryPermissions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntryPermissions.Unmarshal(m, b)
}
func (m *EntryPermissions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntryPermissions.Marshal(b, m, deterministic)
}
func (dst *EntryPermissions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntryPermissions.Merge(dst, src)
}
func (m *EntryPermissions) XXX_Size() int {
	return xxx_messageInfo_EntryPermissions.Size(m)
}
func (m *EntryPermissions) XXX_DiscardUnknown() {
	xxx_messageInfo_EntryPermissions.DiscardUnknown(m)
}

var xxx_messageInfo_EntryPermissions proto.InternalMessageInfo

func (m *EntryPermissions) GetEntry() *Entry {
	if m != nil {
		return m.Entry
	}
	return nil
}

func (m *EntryPermissions) GetPermissionBits() uint64 {
	if m != nil {
		return m.PermissionBits
	}
	return 0
}

func init() {
	proto.RegisterType((*Response)(nil), "acl.Response")
	proto.RegisterType((*Entry)(nil), "acl.Entry")
	proto.RegisterType((*EntryPermissions)(nil), "acl.EntryPermissions")
	proto.RegisterEnum("acl.RequestStatus", RequestStatus_name, RequestStatus_value)
	proto.RegisterEnum("acl.Permissions", Permissions_name, Permissions_value)
	proto.RegisterEnum("acl.EntryType", EntryType_name, EntryType_value)
	proto.RegisterEnum("acl.Flags", Flags_name, Flags_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AccessControlClient is the client API for AccessControl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AccessControlClient interface {
	// Set the permissions on a given ACE or create it if it doesn't exist
	SetPermissions(ctx context.Context, in *EntryPermissions, opts ...grpc.CallOption) (*Response, error)
	// Fetch the permissions on a given ACE
	GetPermissions(ctx context.Context, in *Entry, opts ...grpc.CallOption) (*Response, error)
	// Remove the given ACE completely from the ACL
	DestroyAclEntry(ctx context.Context, in *Entry, opts ...grpc.CallOption) (*Response, error)
}

type accessControlClient struct {
	cc *grpc.ClientConn
}

func NewAccessControlClient(cc *grpc.ClientConn) AccessControlClient {
	return &accessControlClient{cc}
}

func (c *accessControlClient) SetPermissions(ctx context.Context, in *EntryPermissions, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/acl.AccessControl/SetPermissions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessControlClient) GetPermissions(ctx context.Context, in *Entry, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/acl.AccessControl/GetPermissions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessControlClient) DestroyAclEntry(ctx context.Context, in *Entry, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/acl.AccessControl/DestroyAclEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessControlServer is the server API for AccessControl service.
type AccessControlServer interface {
	// Set the permissions on a given ACE or create it if it doesn't exist
	SetPermissions(context.Context, *EntryPermissions) (*Response, error)
	// Fetch the permissions on a given ACE
	GetPermissions(context.Context, *Entry) (*Response, error)
	// Remove the given ACE completely from the ACL
	DestroyAclEntry(context.Context, *Entry) (*Response, error)
}

func RegisterAccessControlServer(s *grpc.Server, srv AccessControlServer) {
	s.RegisterService(&_AccessControl_serviceDesc, srv)
}

func _AccessControl_SetPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntryPermissions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessControlServer).SetPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/acl.AccessControl/SetPermissions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessControlServer).SetPermissions(ctx, req.(*EntryPermissions))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessControl_GetPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Entry)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessControlServer).GetPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/acl.AccessControl/GetPermissions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessControlServer).GetPermissions(ctx, req.(*Entry))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessControl_DestroyAclEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Entry)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessControlServer).DestroyAclEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/acl.AccessControl/DestroyAclEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessControlServer).DestroyAclEntry(ctx, req.(*Entry))
	}
	return interceptor(ctx, in, info, handler)
}

var _AccessControl_serviceDesc = grpc.ServiceDesc{
	ServiceName: "acl.AccessControl",
	HandlerType: (*AccessControlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetPermissions",
			Handler:    _AccessControl_SetPermissions_Handler,
		},
		{
			MethodName: "GetPermissions",
			Handler:    _AccessControl_GetPermissions_Handler,
		},
		{
			MethodName: "DestroyAclEntry",
			Handler:    _AccessControl_DestroyAclEntry_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "acl.proto",
}

func init() { proto.RegisterFile("acl.proto", fileDescriptor_acl_a67c5b5daa5f79f1) }

var fileDescriptor_acl_a67c5b5daa5f79f1 = []byte{
	// 509 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x93, 0x5f, 0xaf, 0xd2, 0x3e,
	0x18, 0xc7, 0x19, 0x07, 0xf8, 0xc1, 0xb3, 0xc3, 0xe8, 0xaf, 0xf1, 0x18, 0x72, 0xa2, 0x86, 0xec,
	0x46, 0x42, 0x94, 0x44, 0xbc, 0xd0, 0xdb, 0xca, 0x0a, 0x36, 0xee, 0x6c, 0xa4, 0x63, 0x92, 0x98,
	0x98, 0x85, 0x83, 0xd5, 0x2c, 0x41, 0x86, 0xb4, 0x5c, 0xf0, 0x8e, 0x7c, 0x07, 0xbe, 0x38, 0xff,
	0x66, 0xdd, 0x60, 0x43, 0x8d, 0x5c, 0xd1, 0xe7, 0xf3, 0xe9, 0xfa, 0xdd, 0xd3, 0x67, 0xd0, 0x5a,
	0xae, 0xd6, 0xc3, 0xed, 0x2e, 0x51, 0x09, 0xbe, 0x58, 0xae, 0xd6, 0x76, 0x02, 0x4d, 0x2e, 0xe4,
	0x36, 0xd9, 0x48, 0x81, 0x07, 0xd0, 0x90, 0x6a, 0xa9, 0xf6, 0xb2, 0x6b, 0xf4, 0x8c, 0xbe, 0x35,
	0xc2, 0xc3, 0x54, 0xe6, 0xe2, 0xd3, 0x5e, 0x48, 0x15, 0x68, 0xc2, 0x73, 0x03, 0x3f, 0x03, 0x73,
	0x2b, 0x76, 0x1f, 0x63, 0x29, 0xe3, 0x64, 0x23, 0xbb, 0xd5, 0x9e, 0xd1, 0x37, 0x47, 0x57, 0x7a,
	0x03, 0xdd, 0xa8, 0xdd, 0x61, 0x56, 0x40, 0x5e, 0x36, 0xed, 0x3d, 0xd4, 0xb5, 0x80, 0x6d, 0xa8,
	0xa9, 0xc3, 0x56, 0xe4, 0x67, 0x59, 0xc5, 0xd6, 0xf9, 0x61, 0x2b, 0xb8, 0x66, 0xf8, 0x0e, 0xd4,
	0xdf, 0xaf, 0x97, 0x1f, 0xb2, 0xe7, 0xb7, 0x79, 0xb6, 0xc0, 0x77, 0xa1, 0x21, 0x36, 0x2a, 0x56,
	0x87, 0xee, 0x45, 0xcf, 0xe8, 0xb7, 0x78, 0xbe, 0xc2, 0xd7, 0xd0, 0x8c, 0xdf, 0xe5, 0xa4, 0xa6,
	0xc9, 0x69, 0x6d, 0xbf, 0x05, 0xf4, 0x7b, 0x2e, 0xdc, 0x83, 0xba, 0x48, 0x6b, 0x3a, 0x82, 0x39,
	0x82, 0x22, 0x02, 0xcf, 0x00, 0x7e, 0x08, 0x9d, 0x22, 0x7b, 0x74, 0x1b, 0xab, 0x2c, 0x49, 0x8d,
	0x5b, 0x45, 0xf9, 0x45, 0xac, 0xe4, 0xe0, 0x8b, 0x01, 0xed, 0xb3, 0x46, 0x61, 0x13, 0xfe, 0x0b,
	0xc2, 0xf1, 0x98, 0x06, 0x01, 0xaa, 0xe0, 0x2e, 0x98, 0x94, 0xf3, 0x28, 0xf4, 0x5e, 0x79, 0xfe,
	0xc2, 0x43, 0x3f, 0x8f, 0x3f, 0x03, 0xdf, 0x83, 0x4e, 0x4a, 0x66, 0x94, 0xdf, 0x44, 0x0e, 0xf5,
	0x18, 0x75, 0xd0, 0x8f, 0x82, 0x3e, 0x80, 0xff, 0x53, 0xca, 0xbc, 0xd7, 0xc4, 0x65, 0x8e, 0xb6,
	0x02, 0xf4, 0xbd, 0xe0, 0x36, 0x5c, 0x9d, 0x71, 0xce, 0xbc, 0x31, 0x9b, 0x11, 0x17, 0x7d, 0x2b,
	0x9c, 0xfb, 0x80, 0xca, 0x4e, 0x18, 0x32, 0x07, 0x7d, 0x3d, 0xe1, 0xc1, 0x13, 0x30, 0xcb, 0x3d,
	0x69, 0x43, 0xcb, 0xf3, 0x23, 0x72, 0x0c, 0xde, 0x84, 0x1a, 0xa7, 0xc4, 0x41, 0x06, 0x6e, 0x41,
	0x7d, 0xc1, 0xd9, 0x9c, 0xa2, 0xea, 0xe0, 0x11, 0xb4, 0x4e, 0x17, 0x95, 0xd6, 0x89, 0xeb, 0xfa,
	0x0b, 0x54, 0xd1, 0x7f, 0x43, 0x87, 0xcd, 0x33, 0x9b, 0xb8, 0x84, 0xdf, 0xa0, 0xea, 0xe0, 0x0d,
	0xd4, 0x27, 0xfa, 0xda, 0x2e, 0xa1, 0xe9, 0xf9, 0xd1, 0xc4, 0x25, 0xd3, 0x20, 0x93, 0xa7, 0xdc,
	0x0f, 0x67, 0xc8, 0xc0, 0x18, 0xac, 0xec, 0xc0, 0xe8, 0xd8, 0xb1, 0x6a, 0xa9, 0x36, 0x21, 0xcc,
	0x0d, 0x39, 0x45, 0x35, 0x8c, 0xe0, 0x72, 0xe6, 0xfb, 0x6e, 0xc4, 0xbc, 0x97, 0x94, 0xb3, 0x39,
	0x6a, 0x8e, 0x3e, 0x1b, 0xd0, 0x26, 0xab, 0x95, 0x90, 0x72, 0x9c, 0x6c, 0xd4, 0x2e, 0x59, 0xe3,
	0xe7, 0x60, 0x05, 0x42, 0x95, 0xdf, 0xe8, 0xef, 0x43, 0x79, 0xdd, 0xce, 0x87, 0x3b, 0x9b, 0x7d,
	0xbb, 0x82, 0x1f, 0x83, 0x35, 0x3d, 0xdf, 0x59, 0x1a, 0x88, 0x3f, 0xf5, 0x21, 0x74, 0x1c, 0x21,
	0xd5, 0x2e, 0x39, 0x90, 0xd5, 0x3a, 0x9b, 0xe8, 0x7f, 0xf9, 0xb7, 0x0d, 0xfd, 0xd1, 0x3d, 0xfd,
	0x15, 0x00, 0x00, 0xff, 0xff, 0x7a, 0x2e, 0x6c, 0xf0, 0x81, 0x03, 0x00, 0x00,
}
