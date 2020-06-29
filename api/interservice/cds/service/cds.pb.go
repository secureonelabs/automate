// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/interservice/cds/service/cds.proto

package service

import (
	context "context"
	fmt "fmt"
	request "github.com/chef/automate/api/external/cds/request"
	response "github.com/chef/automate/api/external/cds/response"
	common "github.com/chef/automate/api/external/common"
	version "github.com/chef/automate/api/external/common/version"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

func init() {
	proto.RegisterFile("api/interservice/cds/service/cds.proto", fileDescriptor_efed58aeb69669da)
}

var fileDescriptor_efed58aeb69669da = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x41, 0x4b, 0xfb, 0x40,
	0x10, 0xc5, 0xff, 0x7f, 0x04, 0x0f, 0xeb, 0x45, 0xd6, 0x5b, 0x8e, 0xb5, 0x08, 0x45, 0xdd, 0xd5,
	0xea, 0x59, 0xd4, 0x56, 0xa4, 0xe0, 0xc9, 0x83, 0x07, 0x2f, 0xb2, 0x6d, 0xa6, 0x76, 0x21, 0xbb,
	0x93, 0xee, 0x4c, 0x6b, 0xc1, 0x4f, 0xeb, 0x37, 0x91, 0x34, 0x5b, 0xdd, 0xd2, 0x94, 0xe2, 0x29,
	0xc9, 0xe4, 0xf7, 0xde, 0x63, 0x86, 0x27, 0x4e, 0x4c, 0x69, 0xb5, 0xf5, 0x0c, 0x81, 0x20, 0xcc,
	0xed, 0x08, 0xf4, 0x28, 0x27, 0x9d, 0xbc, 0xab, 0x32, 0x20, 0xa3, 0xec, 0x8c, 0x26, 0x30, 0x56,
	0x66, 0xc6, 0xe8, 0x0c, 0x83, 0xca, 0xd1, 0x19, 0xeb, 0x97, 0xdf, 0xc6, 0x31, 0xbc, 0x55, 0x60,
	0x14, 0x65, 0xc7, 0x95, 0x25, 0x2c, 0x18, 0x82, 0x37, 0xc5, 0xd2, 0x2e, 0xc0, 0x74, 0x06, 0xc4,
	0x3a, 0x20, 0x72, 0xed, 0x97, 0xb5, 0x1b, 0x20, 0x2a, 0xd1, 0x13, 0xa4, 0x54, 0x6b, 0x9d, 0x42,
	0xe7, 0xd0, 0xff, 0x80, 0x91, 0xe9, 0x34, 0x31, 0x73, 0x08, 0x64, 0x7f, 0x9f, 0x35, 0xda, 0xfd,
	0xda, 0x13, 0x07, 0x77, 0x71, 0x85, 0x5e, 0x4e, 0x72, 0x2a, 0xc4, 0x23, 0xf0, 0x4b, 0xcd, 0xc8,
	0x6b, 0xb5, 0xbe, 0xa3, 0x29, 0xad, 0xaa, 0xed, 0xd4, 0xca, 0x26, 0xa2, 0x03, 0x3f, 0xc6, 0xe7,
	0x7a, 0xa7, 0xec, 0xfc, 0x4f, 0x2a, 0xe9, 0xc4, 0xe1, 0x93, 0x25, 0xee, 0xa1, 0x67, 0xf0, 0x3c,
	0x60, 0x70, 0x24, 0x4f, 0x9b, 0x2c, 0x72, 0x52, 0xf1, 0x6c, 0x2a, 0x85, 0xb3, 0xb3, 0xad, 0x70,
	0xbc, 0xca, 0x9a, 0xf5, 0xa7, 0x90, 0x03, 0x4f, 0x6c, 0x8a, 0x22, 0x19, 0xcb, 0xcb, 0x1d, 0x81,
	0x9b, 0x92, 0xac, 0xbb, 0x2b, 0xb6, 0x21, 0x86, 0xc4, 0x51, 0x1f, 0x3f, 0x7c, 0x81, 0x26, 0x4f,
	0xc7, 0xdd, 0x1d, 0xe9, 0x0d, 0x9a, 0xac, 0xbd, 0xfd, 0xca, 0x0f, 0x8b, 0x12, 0x03, 0xf7, 0x0d,
	0x9b, 0xd6, 0xbf, 0x8b, 0xff, 0xf7, 0xb7, 0xaf, 0x37, 0xef, 0x96, 0x27, 0xb3, 0x61, 0xf5, 0x5f,
	0x57, 0x2a, 0xbd, 0x52, 0xe9, 0x8d, 0xae, 0xa7, 0xfd, 0x5d, 0x95, 0x7e, 0xb8, 0xbf, 0x2c, 0xcb,
	0xd5, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa9, 0xaa, 0x4e, 0x30, 0x1b, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AutomateCdsClient is the client API for AutomateCds service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AutomateCdsClient interface {
	GetVersion(ctx context.Context, in *version.VersionInfoRequest, opts ...grpc.CallOption) (*version.VersionInfo, error)
	ListContentItems(ctx context.Context, in *request.ContentItems, opts ...grpc.CallOption) (*response.ContentItems, error)
	InstallContentItem(ctx context.Context, in *request.InstallContentItem, opts ...grpc.CallOption) (*response.InstallContentItem, error)
	DownloadContentItem(ctx context.Context, in *request.DownloadContentItem, opts ...grpc.CallOption) (AutomateCds_DownloadContentItemClient, error)
}

type automateCdsClient struct {
	cc grpc.ClientConnInterface
}

func NewAutomateCdsClient(cc grpc.ClientConnInterface) AutomateCdsClient {
	return &automateCdsClient{cc}
}

func (c *automateCdsClient) GetVersion(ctx context.Context, in *version.VersionInfoRequest, opts ...grpc.CallOption) (*version.VersionInfo, error) {
	out := new(version.VersionInfo)
	err := c.cc.Invoke(ctx, "/chef.automate.domain.autoamte_cds.service.AutomateCds/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *automateCdsClient) ListContentItems(ctx context.Context, in *request.ContentItems, opts ...grpc.CallOption) (*response.ContentItems, error) {
	out := new(response.ContentItems)
	err := c.cc.Invoke(ctx, "/chef.automate.domain.autoamte_cds.service.AutomateCds/ListContentItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *automateCdsClient) InstallContentItem(ctx context.Context, in *request.InstallContentItem, opts ...grpc.CallOption) (*response.InstallContentItem, error) {
	out := new(response.InstallContentItem)
	err := c.cc.Invoke(ctx, "/chef.automate.domain.autoamte_cds.service.AutomateCds/InstallContentItem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *automateCdsClient) DownloadContentItem(ctx context.Context, in *request.DownloadContentItem, opts ...grpc.CallOption) (AutomateCds_DownloadContentItemClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AutomateCds_serviceDesc.Streams[0], "/chef.automate.domain.autoamte_cds.service.AutomateCds/DownloadContentItem", opts...)
	if err != nil {
		return nil, err
	}
	x := &automateCdsDownloadContentItemClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AutomateCds_DownloadContentItemClient interface {
	Recv() (*common.ExportData, error)
	grpc.ClientStream
}

type automateCdsDownloadContentItemClient struct {
	grpc.ClientStream
}

func (x *automateCdsDownloadContentItemClient) Recv() (*common.ExportData, error) {
	m := new(common.ExportData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AutomateCdsServer is the server API for AutomateCds service.
type AutomateCdsServer interface {
	GetVersion(context.Context, *version.VersionInfoRequest) (*version.VersionInfo, error)
	ListContentItems(context.Context, *request.ContentItems) (*response.ContentItems, error)
	InstallContentItem(context.Context, *request.InstallContentItem) (*response.InstallContentItem, error)
	DownloadContentItem(*request.DownloadContentItem, AutomateCds_DownloadContentItemServer) error
}

// UnimplementedAutomateCdsServer can be embedded to have forward compatible implementations.
type UnimplementedAutomateCdsServer struct {
}

func (*UnimplementedAutomateCdsServer) GetVersion(ctx context.Context, req *version.VersionInfoRequest) (*version.VersionInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (*UnimplementedAutomateCdsServer) ListContentItems(ctx context.Context, req *request.ContentItems) (*response.ContentItems, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListContentItems not implemented")
}
func (*UnimplementedAutomateCdsServer) InstallContentItem(ctx context.Context, req *request.InstallContentItem) (*response.InstallContentItem, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstallContentItem not implemented")
}
func (*UnimplementedAutomateCdsServer) DownloadContentItem(req *request.DownloadContentItem, srv AutomateCds_DownloadContentItemServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadContentItem not implemented")
}

func RegisterAutomateCdsServer(s *grpc.Server, srv AutomateCdsServer) {
	s.RegisterService(&_AutomateCds_serviceDesc, srv)
}

func _AutomateCds_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(version.VersionInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AutomateCdsServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chef.automate.domain.autoamte_cds.service.AutomateCds/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AutomateCdsServer).GetVersion(ctx, req.(*version.VersionInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AutomateCds_ListContentItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.ContentItems)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AutomateCdsServer).ListContentItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chef.automate.domain.autoamte_cds.service.AutomateCds/ListContentItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AutomateCdsServer).ListContentItems(ctx, req.(*request.ContentItems))
	}
	return interceptor(ctx, in, info, handler)
}

func _AutomateCds_InstallContentItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.InstallContentItem)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AutomateCdsServer).InstallContentItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chef.automate.domain.autoamte_cds.service.AutomateCds/InstallContentItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AutomateCdsServer).InstallContentItem(ctx, req.(*request.InstallContentItem))
	}
	return interceptor(ctx, in, info, handler)
}

func _AutomateCds_DownloadContentItem_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(request.DownloadContentItem)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AutomateCdsServer).DownloadContentItem(m, &automateCdsDownloadContentItemServer{stream})
}

type AutomateCds_DownloadContentItemServer interface {
	Send(*common.ExportData) error
	grpc.ServerStream
}

type automateCdsDownloadContentItemServer struct {
	grpc.ServerStream
}

func (x *automateCdsDownloadContentItemServer) Send(m *common.ExportData) error {
	return x.ServerStream.SendMsg(m)
}

var _AutomateCds_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chef.automate.domain.autoamte_cds.service.AutomateCds",
	HandlerType: (*AutomateCdsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _AutomateCds_GetVersion_Handler,
		},
		{
			MethodName: "ListContentItems",
			Handler:    _AutomateCds_ListContentItems_Handler,
		},
		{
			MethodName: "InstallContentItem",
			Handler:    _AutomateCds_InstallContentItem_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadContentItem",
			Handler:       _AutomateCds_DownloadContentItem_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/interservice/cds/service/cds.proto",
}
