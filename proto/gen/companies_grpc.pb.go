// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: companies.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Companies_AddCompany_FullMethodName      = "/companies.Companies/AddCompany"
	Companies_ModifyCompany_FullMethodName   = "/companies.Companies/ModifyCompany"
	Companies_DeleteCompany_FullMethodName   = "/companies.Companies/DeleteCompany"
	Companies_FindCompanyByID_FullMethodName = "/companies.Companies/FindCompanyByID"
)

// CompaniesClient is the client API for Companies service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CompaniesClient interface {
	AddCompany(ctx context.Context, in *AddCompanyRequest, opts ...grpc.CallOption) (*CompanyResponse, error)
	ModifyCompany(ctx context.Context, in *ModifyCompanyRequest, opts ...grpc.CallOption) (*CompanyResponse, error)
	DeleteCompany(ctx context.Context, in *DeleteCompanyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	FindCompanyByID(ctx context.Context, in *FindCompanyByIDRequest, opts ...grpc.CallOption) (*CompanyResponse, error)
}

type companiesClient struct {
	cc grpc.ClientConnInterface
}

func NewCompaniesClient(cc grpc.ClientConnInterface) CompaniesClient {
	return &companiesClient{cc}
}

func (c *companiesClient) AddCompany(ctx context.Context, in *AddCompanyRequest, opts ...grpc.CallOption) (*CompanyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CompanyResponse)
	err := c.cc.Invoke(ctx, Companies_AddCompany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *companiesClient) ModifyCompany(ctx context.Context, in *ModifyCompanyRequest, opts ...grpc.CallOption) (*CompanyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CompanyResponse)
	err := c.cc.Invoke(ctx, Companies_ModifyCompany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *companiesClient) DeleteCompany(ctx context.Context, in *DeleteCompanyRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Companies_DeleteCompany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *companiesClient) FindCompanyByID(ctx context.Context, in *FindCompanyByIDRequest, opts ...grpc.CallOption) (*CompanyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CompanyResponse)
	err := c.cc.Invoke(ctx, Companies_FindCompanyByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompaniesServer is the server API for Companies service.
// All implementations must embed UnimplementedCompaniesServer
// for forward compatibility.
type CompaniesServer interface {
	AddCompany(context.Context, *AddCompanyRequest) (*CompanyResponse, error)
	ModifyCompany(context.Context, *ModifyCompanyRequest) (*CompanyResponse, error)
	DeleteCompany(context.Context, *DeleteCompanyRequest) (*emptypb.Empty, error)
	FindCompanyByID(context.Context, *FindCompanyByIDRequest) (*CompanyResponse, error)
	mustEmbedUnimplementedCompaniesServer()
}

// UnimplementedCompaniesServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCompaniesServer struct{}

func (UnimplementedCompaniesServer) AddCompany(context.Context, *AddCompanyRequest) (*CompanyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCompany not implemented")
}
func (UnimplementedCompaniesServer) ModifyCompany(context.Context, *ModifyCompanyRequest) (*CompanyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyCompany not implemented")
}
func (UnimplementedCompaniesServer) DeleteCompany(context.Context, *DeleteCompanyRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCompany not implemented")
}
func (UnimplementedCompaniesServer) FindCompanyByID(context.Context, *FindCompanyByIDRequest) (*CompanyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindCompanyByID not implemented")
}
func (UnimplementedCompaniesServer) mustEmbedUnimplementedCompaniesServer() {}
func (UnimplementedCompaniesServer) testEmbeddedByValue()                   {}

// UnsafeCompaniesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CompaniesServer will
// result in compilation errors.
type UnsafeCompaniesServer interface {
	mustEmbedUnimplementedCompaniesServer()
}

func RegisterCompaniesServer(s grpc.ServiceRegistrar, srv CompaniesServer) {
	// If the following call pancis, it indicates UnimplementedCompaniesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Companies_ServiceDesc, srv)
}

func _Companies_AddCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).AddCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_AddCompany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).AddCompany(ctx, req.(*AddCompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Companies_ModifyCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyCompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).ModifyCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_ModifyCompany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).ModifyCompany(ctx, req.(*ModifyCompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Companies_DeleteCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).DeleteCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_DeleteCompany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).DeleteCompany(ctx, req.(*DeleteCompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Companies_FindCompanyByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindCompanyByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompaniesServer).FindCompanyByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Companies_FindCompanyByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompaniesServer).FindCompanyByID(ctx, req.(*FindCompanyByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Companies_ServiceDesc is the grpc.ServiceDesc for Companies service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Companies_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "companies.Companies",
	HandlerType: (*CompaniesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddCompany",
			Handler:    _Companies_AddCompany_Handler,
		},
		{
			MethodName: "ModifyCompany",
			Handler:    _Companies_ModifyCompany_Handler,
		},
		{
			MethodName: "DeleteCompany",
			Handler:    _Companies_DeleteCompany_Handler,
		},
		{
			MethodName: "FindCompanyByID",
			Handler:    _Companies_FindCompanyByID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "companies.proto",
}
