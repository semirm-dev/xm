package gateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"xm/companies"
	"xm/companies/datastore"
	"xm/companies/notifiers"
	"xm/gateway"
	"xm/internal/web"
	"xm/proto/gen"
)

const (
	bufSize = 1024 * 1024
	addr    = "8001"
)

var lis *bufconn.Listener

func startGrpcServer() {
	lis = bufconn.Listen(bufSize)
	srv := grpc.NewServer()

	ds := datastore.NewInMemoryStore()
	ds.Companies = []*companies.Company{
		{
			ID:           "123",
			Name:         "Company n",
			Description:  "Company n description",
			EmployeesNum: 1,
			Registered:   true,
			CompanyType:  companies.NonProfitType,
		},
	}
	svc := companies.NewService(ds, notifiers.NewNoopNotifier())
	usrSrv := companies.NewGrpcServer(addr, svc)

	gen.RegisterCompaniesServer(srv, usrSrv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			logrus.Fatalf("grpc server failed: %v", err)
		}
	}()
}

func grpcConn(addr string) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		logrus.Fatal(err)
	}

	return conn
}

func grpcClient() gen.CompaniesClient {
	conn := grpcConn(addr)
	return gen.NewCompaniesClient(conn)
}

func TestAddCompanyHandler(t *testing.T) {
	startGrpcServer()

	router := web.NewRouter()

	router.POST("/companies", gateway.AddCompanyHandler(grpcClient()))

	payload := `{
		"name": "Company 1",
		"description": "Company 1 description",
		"employees_num": 1,
		"registered": true,
		"company_type": 2
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/companies", bytes.NewBuffer([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	var res gen.CompanyResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Id)
	assert.Equal(t, "Company 1", res.Name)
	assert.Equal(t, "Company 1 description", res.Description)
	assert.Equal(t, uint32(1), res.EmployeesNum)
	assert.True(t, res.Registered)
	assert.Equal(t, gen.CompanyType_Cooperative, res.CompanyType)
}

func TestModifyCompanyHandler(t *testing.T) {
	startGrpcServer()

	router := web.NewRouter()

	router.PUT("/companies/:id", gateway.ModifyCompanyHandler(grpcClient()))

	payload := `{
		"name": "Company 1 changed",
		"description": "Company 1 description changed",
		"employees_num": 2,
		"registered": false,
		"company_type": 3
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/companies/123", bytes.NewBuffer([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	var res gen.CompanyResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Id)
	assert.Equal(t, "Company n", res.Name)
	assert.Equal(t, "Company 1 description changed", res.Description)
	assert.Equal(t, uint32(2), res.EmployeesNum)
	assert.False(t, res.Registered)
	assert.Equal(t, gen.CompanyType_SoleProprietorship, res.CompanyType)
}

func TestModifyCompanyHandler_NotExists(t *testing.T) {
	startGrpcServer()

	router := web.NewRouter()

	router.PUT("/companies/:id", gateway.ModifyCompanyHandler(grpcClient()))

	payload := `{
		"name": "Company 1 changed",
		"description": "Company 1 description changed",
		"employees_num": 2,
		"registered": false,
		"company_type": 3
	}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/companies/111", bytes.NewBuffer([]byte(payload)))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	var res gen.CompanyResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.Error(t, err)
}

func TestDeleteCompanyHandler(t *testing.T) {
	startGrpcServer()

	router := web.NewRouter()

	router.DELETE("/companies/:id", gateway.DeleteCompanyHandler(grpcClient()))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/companies/11", nil)
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	var res string
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
}

func TestFindCompanyByIDHandler(t *testing.T) {
	startGrpcServer()

	router := web.NewRouter()

	router.GET("/companies/:id", gateway.FindCompanyByIDHandler(grpcClient()))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/companies/123", nil)
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	var res gen.CompanyResponse
	err := json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.NotEmpty(t, res.Id)
	assert.Equal(t, "Company n", res.Name)
	assert.Equal(t, "Company n description", res.Description)
	assert.Equal(t, uint32(1), res.EmployeesNum)
	assert.True(t, res.Registered)
	assert.Equal(t, gen.CompanyType_NonProfit, res.CompanyType)
}
