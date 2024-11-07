package companies_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"xm/companies"
	"xm/companies/datastore"
	"xm/companies/notifiers"
)

type NotifierMock struct {
	Event    string
	Message  any
	Finished chan bool
}

func (n *NotifierMock) Notify(ctx context.Context, event string, message any) error {
	n.Event = event
	n.Message = message
	n.Finished <- true
	return nil
}

func TestService_AddCompany(t *testing.T) {
	ds := datastore.NewInMemoryStore()
	svc := companies.NewService(ds, notifiers.NewNoopNotifier())

	saved, err := svc.AddCompany(context.Background(), companies.Company{
		Name:         "Company 1",
		Description:  "Company 1 description",
		EmployeesNum: 1,
		Registered:   true,
		CompanyType:  companies.NonProfitType,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, saved.ID)
	assert.Equal(t, "Company 1", saved.Name)
	assert.Equal(t, "Company 1 description", saved.Description)
	assert.Equal(t, uint32(1), saved.EmployeesNum)
	assert.True(t, saved.Registered)
	assert.Equal(t, companies.NonProfitType, saved.CompanyType)

	cmpInDb := ds.Companies[0]
	assert.Equal(t, saved, *cmpInDb)
}

func TestService_AddCompany_AlreadyExists(t *testing.T) {
	ds := datastore.NewInMemoryStore()
	ds.Companies = []*companies.Company{
		{
			ID:   "123",
			Name: "Company 1",
		},
	}
	svc := companies.NewService(ds, notifiers.NewNoopNotifier())

	_, err := svc.AddCompany(context.Background(), companies.Company{
		Name:         "Company 1",
		Description:  "Company 1 description",
		EmployeesNum: 1,
		Registered:   true,
		CompanyType:  companies.NonProfitType,
	})
	assert.Error(t, err)
	assert.Equal(t, errors.New("company Company 1 already exists"), err)
}

func TestService_ModifyCompany(t *testing.T) {
	ds := datastore.NewInMemoryStore()
	ds.Companies = []*companies.Company{
		{
			ID:           "123",
			Name:         "Company 1",
			Description:  "Company 1 description",
			EmployeesNum: 1,
			Registered:   true,
			CompanyType:  companies.NonProfitType,
		},
	}
	svc := companies.NewService(ds, notifiers.NewNoopNotifier())

	modified, err := svc.ModifyCompany(context.Background(), companies.Company{
		ID:           "123",
		Name:         "Company 1 changed",
		Description:  "Company 1 description changed",
		EmployeesNum: 2,
		Registered:   false,
		CompanyType:  companies.CooperativeType,
	})
	assert.NoError(t, err)

	expected := companies.Company{
		ID:           "123",
		Name:         "Company 1",
		Description:  "Company 1 description changed",
		EmployeesNum: 2,
		Registered:   false,
		CompanyType:  companies.CooperativeType,
	}

	assert.Equal(t, expected, modified)
}

func TestService_ModifyCompany_NotExists(t *testing.T) {
	ds := datastore.NewInMemoryStore()
	ds.Companies = []*companies.Company{
		{
			ID:           "123",
			Name:         "Company 1",
			Description:  "Company 1 description",
			EmployeesNum: 1,
			Registered:   true,
			CompanyType:  companies.NonProfitType,
		},
	}
	svc := companies.NewService(ds, notifiers.NewNoopNotifier())

	_, err := svc.ModifyCompany(context.Background(), companies.Company{
		ID:           "456",
		Name:         "Company 1 changed",
		Description:  "Company 1 description changed",
		EmployeesNum: 2,
		Registered:   false,
		CompanyType:  companies.CooperativeType,
	})
	assert.Error(t, err)
	assert.Equal(t, errors.New("company 456 not found"), err)
}

func TestService_DeleteCompany(t *testing.T) {
	ds := datastore.NewInMemoryStore()
	ds.Companies = []*companies.Company{
		{
			ID:           "123",
			Name:         "Company 1",
			Description:  "Company 1 description",
			EmployeesNum: 1,
			Registered:   true,
			CompanyType:  companies.NonProfitType,
		},
	}
	svc := companies.NewService(ds, notifiers.NewNoopNotifier())

	err := svc.DeleteCompany(context.Background(), "123")
	assert.NoError(t, err)
	assert.Len(t, ds.Companies, 0)
}

func TestService_FindCompanyByID(t *testing.T) {
	ds := datastore.NewInMemoryStore()
	ds.Companies = []*companies.Company{
		{
			ID:           "123",
			Name:         "Company 1",
			Description:  "Company 1 description",
			EmployeesNum: 1,
			Registered:   true,
			CompanyType:  companies.NonProfitType,
		},
	}
	svc := companies.NewService(ds, notifiers.NewNoopNotifier())

	found, err := svc.FindCompanyByID(context.Background(), "123")
	assert.NoError(t, err)
	assert.Equal(t, *ds.Companies[0], found)
}

func TestService_AddCompany_TriggerEvent(t *testing.T) {
	ds := datastore.NewInMemoryStore()
	ns := &NotifierMock{Finished: make(chan bool)}
	svc := companies.NewService(ds, ns)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go svc.ListenForNotifications(ctx)

	_, err := svc.AddCompany(context.Background(), companies.Company{
		Name:         "Company 1",
		Description:  "Company 1 description",
		EmployeesNum: 1,
		Registered:   true,
		CompanyType:  companies.NonProfitType,
	})

	<-ns.Finished

	assert.NoError(t, err)
	assert.Equal(t, "company.created", ns.Event)
	assert.NotEmpty(t, ns.Message)
}
