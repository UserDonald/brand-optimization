// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: content/pb/content.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ContentService_GetContentFormats_FullMethodName       = "/content.ContentService/GetContentFormats"
	ContentService_GetContentFormat_FullMethodName        = "/content.ContentService/GetContentFormat"
	ContentService_CreateContentFormat_FullMethodName     = "/content.ContentService/CreateContentFormat"
	ContentService_UpdateContentFormat_FullMethodName     = "/content.ContentService/UpdateContentFormat"
	ContentService_DeleteContentFormat_FullMethodName     = "/content.ContentService/DeleteContentFormat"
	ContentService_GetFormatPerformance_FullMethodName    = "/content.ContentService/GetFormatPerformance"
	ContentService_UpdateFormatPerformance_FullMethodName = "/content.ContentService/UpdateFormatPerformance"
	ContentService_GetScheduledPosts_FullMethodName       = "/content.ContentService/GetScheduledPosts"
	ContentService_GetScheduledPost_FullMethodName        = "/content.ContentService/GetScheduledPost"
	ContentService_SchedulePost_FullMethodName            = "/content.ContentService/SchedulePost"
	ContentService_UpdateScheduledPost_FullMethodName     = "/content.ContentService/UpdateScheduledPost"
	ContentService_DeleteScheduledPost_FullMethodName     = "/content.ContentService/DeleteScheduledPost"
	ContentService_GetPostsDue_FullMethodName             = "/content.ContentService/GetPostsDue"
)

// ContentServiceClient is the client API for ContentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// ContentService provides APIs for managing content formats, performance, and scheduled posts
type ContentServiceClient interface {
	// Content format management
	GetContentFormats(ctx context.Context, in *GetContentFormatsRequest, opts ...grpc.CallOption) (*GetContentFormatsResponse, error)
	GetContentFormat(ctx context.Context, in *GetContentFormatRequest, opts ...grpc.CallOption) (*GetContentFormatResponse, error)
	CreateContentFormat(ctx context.Context, in *CreateContentFormatRequest, opts ...grpc.CallOption) (*CreateContentFormatResponse, error)
	UpdateContentFormat(ctx context.Context, in *UpdateContentFormatRequest, opts ...grpc.CallOption) (*UpdateContentFormatResponse, error)
	DeleteContentFormat(ctx context.Context, in *DeleteContentFormatRequest, opts ...grpc.CallOption) (*DeleteContentFormatResponse, error)
	// Content format performance
	GetFormatPerformance(ctx context.Context, in *GetFormatPerformanceRequest, opts ...grpc.CallOption) (*GetFormatPerformanceResponse, error)
	UpdateFormatPerformance(ctx context.Context, in *UpdateFormatPerformanceRequest, opts ...grpc.CallOption) (*UpdateFormatPerformanceResponse, error)
	// Scheduled posts
	GetScheduledPosts(ctx context.Context, in *GetScheduledPostsRequest, opts ...grpc.CallOption) (*GetScheduledPostsResponse, error)
	GetScheduledPost(ctx context.Context, in *GetScheduledPostRequest, opts ...grpc.CallOption) (*GetScheduledPostResponse, error)
	SchedulePost(ctx context.Context, in *SchedulePostRequest, opts ...grpc.CallOption) (*SchedulePostResponse, error)
	UpdateScheduledPost(ctx context.Context, in *UpdateScheduledPostRequest, opts ...grpc.CallOption) (*UpdateScheduledPostResponse, error)
	DeleteScheduledPost(ctx context.Context, in *DeleteScheduledPostRequest, opts ...grpc.CallOption) (*DeleteScheduledPostResponse, error)
	GetPostsDue(ctx context.Context, in *GetPostsDueRequest, opts ...grpc.CallOption) (*GetPostsDueResponse, error)
}

type contentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewContentServiceClient(cc grpc.ClientConnInterface) ContentServiceClient {
	return &contentServiceClient{cc}
}

func (c *contentServiceClient) GetContentFormats(ctx context.Context, in *GetContentFormatsRequest, opts ...grpc.CallOption) (*GetContentFormatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetContentFormatsResponse)
	err := c.cc.Invoke(ctx, ContentService_GetContentFormats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) GetContentFormat(ctx context.Context, in *GetContentFormatRequest, opts ...grpc.CallOption) (*GetContentFormatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetContentFormatResponse)
	err := c.cc.Invoke(ctx, ContentService_GetContentFormat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) CreateContentFormat(ctx context.Context, in *CreateContentFormatRequest, opts ...grpc.CallOption) (*CreateContentFormatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateContentFormatResponse)
	err := c.cc.Invoke(ctx, ContentService_CreateContentFormat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) UpdateContentFormat(ctx context.Context, in *UpdateContentFormatRequest, opts ...grpc.CallOption) (*UpdateContentFormatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateContentFormatResponse)
	err := c.cc.Invoke(ctx, ContentService_UpdateContentFormat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) DeleteContentFormat(ctx context.Context, in *DeleteContentFormatRequest, opts ...grpc.CallOption) (*DeleteContentFormatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteContentFormatResponse)
	err := c.cc.Invoke(ctx, ContentService_DeleteContentFormat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) GetFormatPerformance(ctx context.Context, in *GetFormatPerformanceRequest, opts ...grpc.CallOption) (*GetFormatPerformanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFormatPerformanceResponse)
	err := c.cc.Invoke(ctx, ContentService_GetFormatPerformance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) UpdateFormatPerformance(ctx context.Context, in *UpdateFormatPerformanceRequest, opts ...grpc.CallOption) (*UpdateFormatPerformanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateFormatPerformanceResponse)
	err := c.cc.Invoke(ctx, ContentService_UpdateFormatPerformance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) GetScheduledPosts(ctx context.Context, in *GetScheduledPostsRequest, opts ...grpc.CallOption) (*GetScheduledPostsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetScheduledPostsResponse)
	err := c.cc.Invoke(ctx, ContentService_GetScheduledPosts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) GetScheduledPost(ctx context.Context, in *GetScheduledPostRequest, opts ...grpc.CallOption) (*GetScheduledPostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetScheduledPostResponse)
	err := c.cc.Invoke(ctx, ContentService_GetScheduledPost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) SchedulePost(ctx context.Context, in *SchedulePostRequest, opts ...grpc.CallOption) (*SchedulePostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SchedulePostResponse)
	err := c.cc.Invoke(ctx, ContentService_SchedulePost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) UpdateScheduledPost(ctx context.Context, in *UpdateScheduledPostRequest, opts ...grpc.CallOption) (*UpdateScheduledPostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateScheduledPostResponse)
	err := c.cc.Invoke(ctx, ContentService_UpdateScheduledPost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) DeleteScheduledPost(ctx context.Context, in *DeleteScheduledPostRequest, opts ...grpc.CallOption) (*DeleteScheduledPostResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteScheduledPostResponse)
	err := c.cc.Invoke(ctx, ContentService_DeleteScheduledPost_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentServiceClient) GetPostsDue(ctx context.Context, in *GetPostsDueRequest, opts ...grpc.CallOption) (*GetPostsDueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPostsDueResponse)
	err := c.cc.Invoke(ctx, ContentService_GetPostsDue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContentServiceServer is the server API for ContentService service.
// All implementations must embed UnimplementedContentServiceServer
// for forward compatibility.
//
// ContentService provides APIs for managing content formats, performance, and scheduled posts
type ContentServiceServer interface {
	// Content format management
	GetContentFormats(context.Context, *GetContentFormatsRequest) (*GetContentFormatsResponse, error)
	GetContentFormat(context.Context, *GetContentFormatRequest) (*GetContentFormatResponse, error)
	CreateContentFormat(context.Context, *CreateContentFormatRequest) (*CreateContentFormatResponse, error)
	UpdateContentFormat(context.Context, *UpdateContentFormatRequest) (*UpdateContentFormatResponse, error)
	DeleteContentFormat(context.Context, *DeleteContentFormatRequest) (*DeleteContentFormatResponse, error)
	// Content format performance
	GetFormatPerformance(context.Context, *GetFormatPerformanceRequest) (*GetFormatPerformanceResponse, error)
	UpdateFormatPerformance(context.Context, *UpdateFormatPerformanceRequest) (*UpdateFormatPerformanceResponse, error)
	// Scheduled posts
	GetScheduledPosts(context.Context, *GetScheduledPostsRequest) (*GetScheduledPostsResponse, error)
	GetScheduledPost(context.Context, *GetScheduledPostRequest) (*GetScheduledPostResponse, error)
	SchedulePost(context.Context, *SchedulePostRequest) (*SchedulePostResponse, error)
	UpdateScheduledPost(context.Context, *UpdateScheduledPostRequest) (*UpdateScheduledPostResponse, error)
	DeleteScheduledPost(context.Context, *DeleteScheduledPostRequest) (*DeleteScheduledPostResponse, error)
	GetPostsDue(context.Context, *GetPostsDueRequest) (*GetPostsDueResponse, error)
	mustEmbedUnimplementedContentServiceServer()
}

// UnimplementedContentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedContentServiceServer struct{}

func (UnimplementedContentServiceServer) GetContentFormats(context.Context, *GetContentFormatsRequest) (*GetContentFormatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContentFormats not implemented")
}
func (UnimplementedContentServiceServer) GetContentFormat(context.Context, *GetContentFormatRequest) (*GetContentFormatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContentFormat not implemented")
}
func (UnimplementedContentServiceServer) CreateContentFormat(context.Context, *CreateContentFormatRequest) (*CreateContentFormatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateContentFormat not implemented")
}
func (UnimplementedContentServiceServer) UpdateContentFormat(context.Context, *UpdateContentFormatRequest) (*UpdateContentFormatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateContentFormat not implemented")
}
func (UnimplementedContentServiceServer) DeleteContentFormat(context.Context, *DeleteContentFormatRequest) (*DeleteContentFormatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteContentFormat not implemented")
}
func (UnimplementedContentServiceServer) GetFormatPerformance(context.Context, *GetFormatPerformanceRequest) (*GetFormatPerformanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFormatPerformance not implemented")
}
func (UnimplementedContentServiceServer) UpdateFormatPerformance(context.Context, *UpdateFormatPerformanceRequest) (*UpdateFormatPerformanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFormatPerformance not implemented")
}
func (UnimplementedContentServiceServer) GetScheduledPosts(context.Context, *GetScheduledPostsRequest) (*GetScheduledPostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScheduledPosts not implemented")
}
func (UnimplementedContentServiceServer) GetScheduledPost(context.Context, *GetScheduledPostRequest) (*GetScheduledPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScheduledPost not implemented")
}
func (UnimplementedContentServiceServer) SchedulePost(context.Context, *SchedulePostRequest) (*SchedulePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SchedulePost not implemented")
}
func (UnimplementedContentServiceServer) UpdateScheduledPost(context.Context, *UpdateScheduledPostRequest) (*UpdateScheduledPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateScheduledPost not implemented")
}
func (UnimplementedContentServiceServer) DeleteScheduledPost(context.Context, *DeleteScheduledPostRequest) (*DeleteScheduledPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteScheduledPost not implemented")
}
func (UnimplementedContentServiceServer) GetPostsDue(context.Context, *GetPostsDueRequest) (*GetPostsDueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostsDue not implemented")
}
func (UnimplementedContentServiceServer) mustEmbedUnimplementedContentServiceServer() {}
func (UnimplementedContentServiceServer) testEmbeddedByValue()                        {}

// UnsafeContentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContentServiceServer will
// result in compilation errors.
type UnsafeContentServiceServer interface {
	mustEmbedUnimplementedContentServiceServer()
}

func RegisterContentServiceServer(s grpc.ServiceRegistrar, srv ContentServiceServer) {
	// If the following call pancis, it indicates UnimplementedContentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ContentService_ServiceDesc, srv)
}

func _ContentService_GetContentFormats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetContentFormatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetContentFormats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetContentFormats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetContentFormats(ctx, req.(*GetContentFormatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_GetContentFormat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetContentFormatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetContentFormat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetContentFormat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetContentFormat(ctx, req.(*GetContentFormatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_CreateContentFormat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateContentFormatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).CreateContentFormat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_CreateContentFormat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).CreateContentFormat(ctx, req.(*CreateContentFormatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_UpdateContentFormat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateContentFormatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).UpdateContentFormat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_UpdateContentFormat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).UpdateContentFormat(ctx, req.(*UpdateContentFormatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_DeleteContentFormat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteContentFormatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).DeleteContentFormat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_DeleteContentFormat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).DeleteContentFormat(ctx, req.(*DeleteContentFormatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_GetFormatPerformance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFormatPerformanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetFormatPerformance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetFormatPerformance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetFormatPerformance(ctx, req.(*GetFormatPerformanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_UpdateFormatPerformance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFormatPerformanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).UpdateFormatPerformance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_UpdateFormatPerformance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).UpdateFormatPerformance(ctx, req.(*UpdateFormatPerformanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_GetScheduledPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScheduledPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetScheduledPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetScheduledPosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetScheduledPosts(ctx, req.(*GetScheduledPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_GetScheduledPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScheduledPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetScheduledPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetScheduledPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetScheduledPost(ctx, req.(*GetScheduledPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_SchedulePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SchedulePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).SchedulePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_SchedulePost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).SchedulePost(ctx, req.(*SchedulePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_UpdateScheduledPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateScheduledPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).UpdateScheduledPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_UpdateScheduledPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).UpdateScheduledPost(ctx, req.(*UpdateScheduledPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_DeleteScheduledPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteScheduledPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).DeleteScheduledPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_DeleteScheduledPost_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).DeleteScheduledPost(ctx, req.(*DeleteScheduledPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentService_GetPostsDue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostsDueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServiceServer).GetPostsDue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentService_GetPostsDue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServiceServer).GetPostsDue(ctx, req.(*GetPostsDueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ContentService_ServiceDesc is the grpc.ServiceDesc for ContentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ContentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "content.ContentService",
	HandlerType: (*ContentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetContentFormats",
			Handler:    _ContentService_GetContentFormats_Handler,
		},
		{
			MethodName: "GetContentFormat",
			Handler:    _ContentService_GetContentFormat_Handler,
		},
		{
			MethodName: "CreateContentFormat",
			Handler:    _ContentService_CreateContentFormat_Handler,
		},
		{
			MethodName: "UpdateContentFormat",
			Handler:    _ContentService_UpdateContentFormat_Handler,
		},
		{
			MethodName: "DeleteContentFormat",
			Handler:    _ContentService_DeleteContentFormat_Handler,
		},
		{
			MethodName: "GetFormatPerformance",
			Handler:    _ContentService_GetFormatPerformance_Handler,
		},
		{
			MethodName: "UpdateFormatPerformance",
			Handler:    _ContentService_UpdateFormatPerformance_Handler,
		},
		{
			MethodName: "GetScheduledPosts",
			Handler:    _ContentService_GetScheduledPosts_Handler,
		},
		{
			MethodName: "GetScheduledPost",
			Handler:    _ContentService_GetScheduledPost_Handler,
		},
		{
			MethodName: "SchedulePost",
			Handler:    _ContentService_SchedulePost_Handler,
		},
		{
			MethodName: "UpdateScheduledPost",
			Handler:    _ContentService_UpdateScheduledPost_Handler,
		},
		{
			MethodName: "DeleteScheduledPost",
			Handler:    _ContentService_DeleteScheduledPost_Handler,
		},
		{
			MethodName: "GetPostsDue",
			Handler:    _ContentService_GetPostsDue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "content/pb/content.proto",
}
