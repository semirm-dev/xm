package companies

import (
	"context"
	"github.com/sirupsen/logrus"
	grpcLib "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"xm/internal/grpc"
	"xm/proto/gen"
)

// GrpcServer exposes companies Service over grpc.
type GrpcServer struct {
	gen.UnimplementedCompaniesServer
	addr   string
	cmpSvc *Service
}

func NewGrpcServer(addr string, cmpSvc *Service) *GrpcServer {
	return &GrpcServer{
		addr:   addr,
		cmpSvc: cmpSvc,
	}
}

func (srv *GrpcServer) StartListening(ctx context.Context) {
	grpc.Start(ctx, srv.addr, "companies", func(s *grpcLib.Server) {
		gen.RegisterCompaniesServer(s, srv)
	})
}

func (srv *GrpcServer) AddCompany(ctx context.Context, req *gen.AddCompanyRequest) (*gen.CompanyResponse, error) {
	cmp, err := srv.cmpSvc.AddCompany(ctx, Company{
		Name:         req.Name,
		Description:  req.Description,
		EmployeesNum: req.EmployeesNum,
		Registered:   req.Registered,
		CompanyType:  CompanyType(req.CompanyType),
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("new company saved: %v", cmp)

	return mapCompanyToProto(cmp), nil
}

func (srv *GrpcServer) ModifyCompany(ctx context.Context, req *gen.ModifyCompanyRequest) (*gen.CompanyResponse, error) {
	cmp, err := srv.cmpSvc.ModifyCompany(ctx, Company{
		ID:           req.Id,
		Description:  req.Description,
		EmployeesNum: req.EmployeesNum,
		Registered:   req.Registered,
		CompanyType:  CompanyType(req.CompanyType),
	})
	if err != nil {
		return nil, err
	}

	logrus.Infof("company modified: %v", cmp)

	return mapCompanyToProto(cmp), nil
}

func (srv *GrpcServer) DeleteCompany(ctx context.Context, req *gen.DeleteCompanyRequest) (*emptypb.Empty, error) {
	if err := srv.cmpSvc.DeleteCompany(ctx, req.Id); err != nil {
		return nil, err
	}

	logrus.Infof("company deleted: %s", req.Id)

	return &emptypb.Empty{}, nil
}

func (srv *GrpcServer) FindCompanyByID(ctx context.Context, req *gen.FindCompanyByIDRequest) (*gen.CompanyResponse, error) {
	cmp, err := srv.cmpSvc.FindCompanyByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return mapCompanyToProto(cmp), nil
}

func mapCompanyToProto(cmp Company) *gen.CompanyResponse {
	return &gen.CompanyResponse{
		Id:           cmp.ID,
		Name:         cmp.Name,
		Description:  cmp.Description,
		EmployeesNum: cmp.EmployeesNum,
		Registered:   cmp.Registered,
		CompanyType:  gen.CompanyType(cmp.CompanyType),
	}
}
