package companies

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	companyCreatedEvent  = "company.created"
	companyModifiedEvent = "company.modified"
	companyDeletedEvent  = "company.deleted"
)

const (
	CorporationsType       CompanyType = 0
	NonProfitType          CompanyType = 1
	CooperativeType        CompanyType = 2
	SoleProprietorShipType CompanyType = 3
)

type CompanyType int

// Service for companies management.
// For now the service is exposed and accessed over gRPC only.
type Service struct {
	ds            DataStore
	notifier      Notifier
	notifications chan notification
}

// Company main model.
type Company struct {
	ID           string
	Name         string
	Description  string
	EmployeesNum uint32
	Registered   bool
	CompanyType  CompanyType
}

type notification struct {
	Event   string
	Message any
}

// DataStore is main data storage for companies.
type DataStore interface {
	Save(ctx context.Context, company Company) (Company, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (Company, error)
	FindByName(ctx context.Context, name string) (Company, error)
}

// Notifier is responsible for notifying other systems about companies actions.
type Notifier interface {
	Notify(ctx context.Context, event string, message any) error
}

func NewService(ds DataStore, notifier Notifier) *Service {
	return &Service{
		ds:            ds,
		notifier:      notifier,
		notifications: make(chan notification, 10),
	}
}

// AddCompany will add new company in the system.
func (svc *Service) AddCompany(ctx context.Context, company Company) (Company, error) {
	if err := svc.companyCreateIsValid(ctx, company); err != nil {
		return Company{}, err
	}

	company.ID = uuid.New().String()

	saved, err := svc.ds.Save(ctx, company)
	if err != nil {
		return Company{}, err
	}

	go svc.notify(companyCreatedEvent, saved)

	return saved, nil
}

// ModifyCompany will modify an existing company.
// Only small subset of fields is allowed to be updated.
func (svc *Service) ModifyCompany(ctx context.Context, company Company) (Company, error) {
	found, err := svc.ds.FindByID(ctx, company.ID)
	if err != nil {
		return Company{}, err
	}

	if found.ID == "" {
		return Company{}, fmt.Errorf("company %s not found", company.ID)
	}

	found.Description = company.Description
	found.EmployeesNum = company.EmployeesNum
	found.Registered = company.Registered
	found.CompanyType = company.CompanyType

	if err = validateCompanyFields(found); err != nil {
		return Company{}, err
	}

	saved, err := svc.ds.Save(ctx, found)
	if err != nil {
		return Company{}, err
	}

	go svc.notify(companyModifiedEvent, saved)

	return saved, nil
}

// DeleteCompany will delete an existing company based on company id.
func (svc *Service) DeleteCompany(ctx context.Context, id string) error {
	if err := svc.ds.Delete(ctx, id); err != nil {
		return err
	}

	go svc.notify(companyDeletedEvent, id)

	return nil
}

// FindCompanyByID will find a company by a given id.
func (svc *Service) FindCompanyByID(ctx context.Context, id string) (Company, error) {
	return svc.ds.FindByID(ctx, id)
}

// ListenForNotifications will listen for notifications ready to be sent to whoever is listening.
// Listener is called in its own goroutine so the rest of the system is not blocked.
func (svc *Service) ListenForNotifications(ctx context.Context) {
	defer func() {
		logrus.Warn("Notifications listener stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case n, ok := <-svc.notifications:
			if !ok {
				return
			}

			go func(ctx context.Context, n notification) {
				if err := svc.notifier.Notify(ctx, n.Event, n.Message); err != nil {
					logrus.Error(err)
				}
			}(ctx, n)
		}
	}
}

// notify is used internally to trigger a notification about executed actions, such as company created, modified, deleted.
// Thread-safe notification sender.
func (svc *Service) notify(event string, message any) {
	svc.notifications <- notification{
		Event:   event,
		Message: message,
	}
}

func (svc *Service) companyCreateIsValid(ctx context.Context, company Company) error {
	if err := validateCompanyFields(company); err != nil {
		return err
	}

	found, err := svc.ds.FindByName(ctx, company.Name)
	if err != nil {
		return err
	}

	if found.ID != "" {
		return fmt.Errorf("company %s already exists", company.Name)
	}

	return nil
}

func validateCompanyFields(company Company) error {
	var err error

	if company.Name == "" {
		err = errors.Join(err, errors.New("missing company name"))
	}

	if len(company.Name) > 15 {
		err = errors.Join(err, errors.New("invalid company name length"))
	}

	if len(company.Description) > 3000 {
		err = errors.Join(err, errors.New("invalid company description length"))
	}

	if company.EmployeesNum <= 0 {
		err = errors.Join(err, errors.New("invalid company amount of employees"))
	}

	return err
}
