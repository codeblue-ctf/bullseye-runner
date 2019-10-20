// Code generated by protoc-gen-go. DO NOT EDIT.
// source: runner.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RunnerRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Timeout              uint64   `protobuf:"varint,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Yml                  string   `protobuf:"bytes,3,opt,name=yml,proto3" json:"yml,omitempty"`
	RegistryHost         string   `protobuf:"bytes,4,opt,name=registry_host,json=registryHost,proto3" json:"registry_host,omitempty"`
	RegistryUsername     string   `protobuf:"bytes,5,opt,name=registry_username,json=registryUsername,proto3" json:"registry_username,omitempty"`
	RegistryPassword     string   `protobuf:"bytes,6,opt,name=registry_password,json=registryPassword,proto3" json:"registry_password,omitempty"`
	FlagTemplate         string   `protobuf:"bytes,7,opt,name=flag_template,json=flagTemplate,proto3" json:"flag_template,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunnerRequest) Reset()         { *m = RunnerRequest{} }
func (m *RunnerRequest) String() string { return proto.CompactTextString(m) }
func (*RunnerRequest) ProtoMessage()    {}
func (*RunnerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{0}
}

func (m *RunnerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunnerRequest.Unmarshal(m, b)
}
func (m *RunnerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunnerRequest.Marshal(b, m, deterministic)
}
func (m *RunnerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunnerRequest.Merge(m, src)
}
func (m *RunnerRequest) XXX_Size() int {
	return xxx_messageInfo_RunnerRequest.Size(m)
}
func (m *RunnerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RunnerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RunnerRequest proto.InternalMessageInfo

func (m *RunnerRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *RunnerRequest) GetTimeout() uint64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *RunnerRequest) GetYml() string {
	if m != nil {
		return m.Yml
	}
	return ""
}

func (m *RunnerRequest) GetRegistryHost() string {
	if m != nil {
		return m.RegistryHost
	}
	return ""
}

func (m *RunnerRequest) GetRegistryUsername() string {
	if m != nil {
		return m.RegistryUsername
	}
	return ""
}

func (m *RunnerRequest) GetRegistryPassword() string {
	if m != nil {
		return m.RegistryPassword
	}
	return ""
}

func (m *RunnerRequest) GetFlagTemplate() string {
	if m != nil {
		return m.FlagTemplate
	}
	return ""
}

type RunnerResponse struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Succeeded            bool     `protobuf:"varint,2,opt,name=succeeded,proto3" json:"succeeded,omitempty"`
	Output               string   `protobuf:"bytes,3,opt,name=output,proto3" json:"output,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunnerResponse) Reset()         { *m = RunnerResponse{} }
func (m *RunnerResponse) String() string { return proto.CompactTextString(m) }
func (*RunnerResponse) ProtoMessage()    {}
func (*RunnerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{1}
}

func (m *RunnerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunnerResponse.Unmarshal(m, b)
}
func (m *RunnerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunnerResponse.Marshal(b, m, deterministic)
}
func (m *RunnerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunnerResponse.Merge(m, src)
}
func (m *RunnerResponse) XXX_Size() int {
	return xxx_messageInfo_RunnerResponse.Size(m)
}
func (m *RunnerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RunnerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RunnerResponse proto.InternalMessageInfo

func (m *RunnerResponse) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *RunnerResponse) GetSucceeded() bool {
	if m != nil {
		return m.Succeeded
	}
	return false
}

func (m *RunnerResponse) GetOutput() string {
	if m != nil {
		return m.Output
	}
	return ""
}

func init() {
	proto.RegisterType((*RunnerRequest)(nil), "proto.RunnerRequest")
	proto.RegisterType((*RunnerResponse)(nil), "proto.RunnerResponse")
}

func init() { proto.RegisterFile("runner.proto", fileDescriptor_48eceea7e2abc593) }

var fileDescriptor_48eceea7e2abc593 = []byte{
	// 264 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x89, 0x49, 0x53, 0x3b, 0xb4, 0x52, 0x17, 0x95, 0x45, 0x3c, 0x94, 0x7a, 0x29, 0x08,
	0x3d, 0xd4, 0xab, 0x0f, 0xe0, 0x51, 0x16, 0xbd, 0x78, 0x29, 0xb1, 0x19, 0x6b, 0x20, 0xc9, 0xc6,
	0x9d, 0x19, 0xa4, 0xcf, 0xed, 0x0b, 0x48, 0x36, 0x1b, 0x25, 0xd2, 0x53, 0x66, 0xbe, 0xf9, 0xc8,
	0xec, 0x3f, 0x30, 0x75, 0x52, 0xd7, 0xe8, 0xd6, 0x8d, 0xb3, 0x6c, 0xd5, 0xc8, 0x7f, 0x96, 0xdf,
	0x11, 0xcc, 0x8c, 0xe7, 0x06, 0x3f, 0x05, 0x89, 0x95, 0x82, 0x44, 0xa4, 0xc8, 0x75, 0xb4, 0x88,
	0x56, 0x13, 0xe3, 0x6b, 0xa5, 0x61, 0xcc, 0x45, 0x85, 0x56, 0x58, 0x9f, 0x2c, 0xa2, 0x55, 0x62,
	0xfa, 0x56, 0xcd, 0x21, 0x3e, 0x54, 0xa5, 0x8e, 0xbd, 0xdc, 0x96, 0xea, 0x16, 0x66, 0x0e, 0xf7,
	0x05, 0xb1, 0x3b, 0x6c, 0x3f, 0x2c, 0xb1, 0x4e, 0xfc, 0x6c, 0xda, 0xc3, 0x47, 0x4b, 0xac, 0xee,
	0xe0, 0xfc, 0x57, 0x12, 0x42, 0x57, 0x67, 0x15, 0xea, 0x91, 0x17, 0xe7, 0xfd, 0xe0, 0x25, 0xf0,
	0x81, 0xdc, 0x64, 0x44, 0x5f, 0xd6, 0xe5, 0x3a, 0x1d, 0xca, 0x4f, 0x81, 0xb7, 0xeb, 0xdf, 0xcb,
	0x6c, 0xbf, 0x65, 0xac, 0x9a, 0x32, 0x63, 0xd4, 0xe3, 0x6e, 0x7d, 0x0b, 0x9f, 0x03, 0x5b, 0xbe,
	0xc2, 0x59, 0x1f, 0x9a, 0x1a, 0x5b, 0x13, 0x1e, 0x4d, 0x7d, 0x03, 0x13, 0x92, 0xdd, 0x0e, 0x31,
	0xc7, 0xdc, 0xe7, 0x3e, 0x35, 0x7f, 0x40, 0x5d, 0x41, 0x6a, 0x85, 0x1b, 0xe1, 0x10, 0x3e, 0x74,
	0x9b, 0x07, 0x48, 0xbb, 0x7f, 0xab, 0x0d, 0xc4, 0x46, 0x6a, 0x75, 0xd1, 0x5d, 0x7c, 0x3d, 0x38,
	0xf3, 0xf5, 0xe5, 0x3f, 0xda, 0xbd, 0xe3, 0x2d, 0xf5, 0xf4, 0xfe, 0x27, 0x00, 0x00, 0xff, 0xff,
	0x14, 0xe5, 0xfd, 0xaf, 0xad, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RunnerClient is the client API for Runner service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RunnerClient interface {
	Run(ctx context.Context, in *RunnerRequest, opts ...grpc.CallOption) (*RunnerResponse, error)
}

type runnerClient struct {
	cc *grpc.ClientConn
}

func NewRunnerClient(cc *grpc.ClientConn) RunnerClient {
	return &runnerClient{cc}
}

func (c *runnerClient) Run(ctx context.Context, in *RunnerRequest, opts ...grpc.CallOption) (*RunnerResponse, error) {
	out := new(RunnerResponse)
	err := c.cc.Invoke(ctx, "/proto.Runner/Run", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RunnerServer is the server API for Runner service.
type RunnerServer interface {
	Run(context.Context, *RunnerRequest) (*RunnerResponse, error)
}

func RegisterRunnerServer(s *grpc.Server, srv RunnerServer) {
	s.RegisterService(&_Runner_serviceDesc, srv)
}

func _Runner_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunnerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Runner/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).Run(ctx, req.(*RunnerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Runner_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Runner",
	HandlerType: (*RunnerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Run",
			Handler:    _Runner_Run_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "runner.proto",
}
