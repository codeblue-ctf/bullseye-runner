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
	X11Info              *X11Info `protobuf:"bytes,8,opt,name=x11info,proto3" json:"x11info,omitempty"`
	PullImage            bool     `protobuf:"varint,9,opt,name=pull_image,json=pullImage,proto3" json:"pull_image,omitempty"`
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

func (m *RunnerRequest) GetX11Info() *X11Info {
	if m != nil {
		return m.X11Info
	}
	return nil
}

func (m *RunnerRequest) GetPullImage() bool {
	if m != nil {
		return m.PullImage
	}
	return false
}

type X11Info struct {
	Width                uint64   `protobuf:"varint,1,opt,name=width,proto3" json:"width,omitempty"`
	Height               uint64   `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Depth                uint64   `protobuf:"varint,3,opt,name=depth,proto3" json:"depth,omitempty"`
	CapExt               string   `protobuf:"bytes,4,opt,name=cap_ext,json=capExt,proto3" json:"cap_ext,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *X11Info) Reset()         { *m = X11Info{} }
func (m *X11Info) String() string { return proto.CompactTextString(m) }
func (*X11Info) ProtoMessage()    {}
func (*X11Info) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{1}
}

func (m *X11Info) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_X11Info.Unmarshal(m, b)
}
func (m *X11Info) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_X11Info.Marshal(b, m, deterministic)
}
func (m *X11Info) XXX_Merge(src proto.Message) {
	xxx_messageInfo_X11Info.Merge(m, src)
}
func (m *X11Info) XXX_Size() int {
	return xxx_messageInfo_X11Info.Size(m)
}
func (m *X11Info) XXX_DiscardUnknown() {
	xxx_messageInfo_X11Info.DiscardUnknown(m)
}

var xxx_messageInfo_X11Info proto.InternalMessageInfo

func (m *X11Info) GetWidth() uint64 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *X11Info) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *X11Info) GetDepth() uint64 {
	if m != nil {
		return m.Depth
	}
	return 0
}

func (m *X11Info) GetCapExt() string {
	if m != nil {
		return m.CapExt
	}
	return ""
}

type RunnerResponse struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Succeeded            bool     `protobuf:"varint,2,opt,name=succeeded,proto3" json:"succeeded,omitempty"`
	Output               string   `protobuf:"bytes,3,opt,name=output,proto3" json:"output,omitempty"`
	X11Cap               []byte   `protobuf:"bytes,4,opt,name=x11cap,proto3" json:"x11cap,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunnerResponse) Reset()         { *m = RunnerResponse{} }
func (m *RunnerResponse) String() string { return proto.CompactTextString(m) }
func (*RunnerResponse) ProtoMessage()    {}
func (*RunnerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{2}
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

func (m *RunnerResponse) GetX11Cap() []byte {
	if m != nil {
		return m.X11Cap
	}
	return nil
}

type InfoRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoRequest) Reset()         { *m = InfoRequest{} }
func (m *InfoRequest) String() string { return proto.CompactTextString(m) }
func (*InfoRequest) ProtoMessage()    {}
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{3}
}

func (m *InfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoRequest.Unmarshal(m, b)
}
func (m *InfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoRequest.Marshal(b, m, deterministic)
}
func (m *InfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoRequest.Merge(m, src)
}
func (m *InfoRequest) XXX_Size() int {
	return xxx_messageInfo_InfoRequest.Size(m)
}
func (m *InfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InfoRequest proto.InternalMessageInfo

type InfoResponse struct {
	Cpus                 uint64   `protobuf:"varint,1,opt,name=cpus,proto3" json:"cpus,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoResponse) Reset()         { *m = InfoResponse{} }
func (m *InfoResponse) String() string { return proto.CompactTextString(m) }
func (*InfoResponse) ProtoMessage()    {}
func (*InfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{4}
}

func (m *InfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoResponse.Unmarshal(m, b)
}
func (m *InfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoResponse.Marshal(b, m, deterministic)
}
func (m *InfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoResponse.Merge(m, src)
}
func (m *InfoResponse) XXX_Size() int {
	return xxx_messageInfo_InfoResponse.Size(m)
}
func (m *InfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InfoResponse proto.InternalMessageInfo

func (m *InfoResponse) GetCpus() uint64 {
	if m != nil {
		return m.Cpus
	}
	return 0
}

func init() {
	proto.RegisterType((*RunnerRequest)(nil), "proto.RunnerRequest")
	proto.RegisterType((*X11Info)(nil), "proto.X11Info")
	proto.RegisterType((*RunnerResponse)(nil), "proto.RunnerResponse")
	proto.RegisterType((*InfoRequest)(nil), "proto.InfoRequest")
	proto.RegisterType((*InfoResponse)(nil), "proto.InfoResponse")
}

func init() { proto.RegisterFile("runner.proto", fileDescriptor_48eceea7e2abc593) }

var fileDescriptor_48eceea7e2abc593 = []byte{
	// 411 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xdf, 0x8a, 0xd3, 0x40,
	0x14, 0xc6, 0xc9, 0x36, 0x4d, 0xda, 0xb3, 0xed, 0xb2, 0x8e, 0xab, 0x0e, 0x8b, 0x42, 0x89, 0x37,
	0x01, 0x61, 0x25, 0xf5, 0x19, 0x04, 0xf7, 0x4e, 0x06, 0x05, 0xef, 0xca, 0x98, 0x9c, 0x36, 0x81,
	0x24, 0x33, 0xce, 0x1f, 0xb6, 0xfb, 0x6e, 0x3e, 0x9c, 0xcc, 0x9f, 0xb8, 0x46, 0xf6, 0xaa, 0xf3,
	0xfd, 0xce, 0x37, 0x3d, 0x27, 0xdf, 0x19, 0xd8, 0x28, 0x3b, 0x8e, 0xa8, 0xee, 0xa4, 0x12, 0x46,
	0x90, 0xa5, 0xff, 0x29, 0x7e, 0x5f, 0xc0, 0x96, 0x79, 0xce, 0xf0, 0x97, 0x45, 0x6d, 0x08, 0x81,
	0xd4, 0xda, 0xae, 0xa1, 0xc9, 0x2e, 0x29, 0xd7, 0xcc, 0x9f, 0x09, 0x85, 0xdc, 0x74, 0x03, 0x0a,
	0x6b, 0xe8, 0xc5, 0x2e, 0x29, 0x53, 0x36, 0x49, 0x72, 0x0d, 0x8b, 0xc7, 0xa1, 0xa7, 0x0b, 0x6f,
	0x76, 0x47, 0xf2, 0x1e, 0xb6, 0x0a, 0x4f, 0x9d, 0x36, 0xea, 0xf1, 0xd0, 0x0a, 0x6d, 0x68, 0xea,
	0x6b, 0x9b, 0x09, 0x7e, 0x11, 0xda, 0x90, 0x0f, 0xf0, 0xe2, 0xaf, 0xc9, 0x6a, 0x54, 0x23, 0x1f,
	0x90, 0x2e, 0xbd, 0xf1, 0x7a, 0x2a, 0x7c, 0x8f, 0x7c, 0x66, 0x96, 0x5c, 0xeb, 0x07, 0xa1, 0x1a,
	0x9a, 0xcd, 0xcd, 0x5f, 0x23, 0x77, 0xed, 0x8f, 0x3d, 0x3f, 0x1d, 0x0c, 0x0e, 0xb2, 0xe7, 0x06,
	0x69, 0x1e, 0xda, 0x3b, 0xf8, 0x2d, 0x32, 0x52, 0x42, 0x7e, 0xae, 0xaa, 0x6e, 0x3c, 0x0a, 0xba,
	0xda, 0x25, 0xe5, 0xe5, 0xfe, 0x2a, 0xa4, 0x72, 0xf7, 0xa3, 0xaa, 0xee, 0xc7, 0xa3, 0x60, 0x53,
	0x99, 0xbc, 0x03, 0x90, 0xb6, 0xef, 0x0f, 0xdd, 0xc0, 0x4f, 0x48, 0xd7, 0xbb, 0xa4, 0x5c, 0xb1,
	0xb5, 0x23, 0xf7, 0x0e, 0x14, 0x47, 0xc8, 0xe3, 0x15, 0x72, 0x03, 0xcb, 0x87, 0xae, 0x31, 0xad,
	0x0f, 0x2e, 0x65, 0x41, 0x90, 0xd7, 0x90, 0xb5, 0xd8, 0x9d, 0xda, 0x29, 0xb8, 0xa8, 0x9c, 0xbb,
	0x41, 0x69, 0x5a, 0x9f, 0x5c, 0xca, 0x82, 0x20, 0x6f, 0x20, 0xaf, 0xb9, 0x3c, 0xe0, 0x79, 0x4a,
	0x2d, 0xab, 0xb9, 0xfc, 0x7c, 0x36, 0x85, 0x82, 0xab, 0x69, 0x4b, 0x5a, 0x8a, 0x51, 0xe3, 0xb3,
	0x6b, 0x7a, 0x0b, 0x6b, 0x6d, 0xeb, 0x1a, 0xb1, 0xc1, 0xc6, 0xf7, 0x5b, 0xb1, 0x27, 0xe0, 0x46,
	0x11, 0xd6, 0x48, 0x6b, 0xe2, 0xb6, 0xa2, 0x72, 0xfc, 0x5c, 0x55, 0x35, 0x97, 0xbe, 0xe7, 0x86,
	0x45, 0x55, 0x6c, 0xe1, 0xd2, 0x67, 0x11, 0xde, 0x45, 0x51, 0xc0, 0x26, 0xc8, 0xa7, 0x01, 0x6a,
	0x69, 0x75, 0xfc, 0x5c, 0x7f, 0xde, 0x0f, 0x90, 0x85, 0x31, 0xc9, 0x1e, 0x16, 0xcc, 0x8e, 0xe4,
	0x26, 0xe6, 0x3a, 0x7b, 0x62, 0xb7, 0xaf, 0xfe, 0xa3, 0xf1, 0x1f, 0x3f, 0x42, 0xea, 0x93, 0x24,
	0xb1, 0xfc, 0x4f, 0xf7, 0xdb, 0x97, 0x33, 0x16, 0x2e, 0xfc, 0xcc, 0x3c, 0xfb, 0xf4, 0x27, 0x00,
	0x00, 0xff, 0xff, 0x5e, 0xbd, 0x14, 0x97, 0xda, 0x02, 0x00, 0x00,
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
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
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

func (c *runnerClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, "/proto.Runner/Info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RunnerServer is the server API for Runner service.
type RunnerServer interface {
	Run(context.Context, *RunnerRequest) (*RunnerResponse, error)
	Info(context.Context, *InfoRequest) (*InfoResponse, error)
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

func _Runner_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Runner/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).Info(ctx, req.(*InfoRequest))
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
		{
			MethodName: "Info",
			Handler:    _Runner_Info_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "runner.proto",
}
