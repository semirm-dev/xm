package datastore

import (
	"context"
	"github.com/samber/lo"
	"xm/companies"
)

type InMemoryStore struct {
	Companies []*companies.Company
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Companies: make([]*companies.Company, 0),
	}
}

func (s *InMemoryStore) Save(ctx context.Context, company companies.Company) (companies.Company, error) {
	s.Companies = append(s.Companies, &company)
	return company, nil
}

func (s *InMemoryStore) Delete(ctx context.Context, id string) error {
	for i, c := range s.Companies {
		if c.ID == id {
			copy(s.Companies[i:], s.Companies[i+1:])
			s.Companies[len(s.Companies)-1] = nil
			s.Companies = s.Companies[:len(s.Companies)-1]
			break
		}
	}
	return nil
}

func (s *InMemoryStore) FindByID(ctx context.Context, id string) (companies.Company, error) {
	cmp, found := lo.Find(s.Companies, func(c *companies.Company) bool {
		return c.ID == id
	})

	if !found {
		return companies.Company{}, nil
	}

	return *cmp, nil
}

func (s *InMemoryStore) FindByName(ctx context.Context, name string) (companies.Company, error) {
	cmp, found := lo.Find(s.Companies, func(c *companies.Company) bool {
		return c.Name == name
	})

	if !found {
		return companies.Company{}, nil
	}

	return *cmp, nil
}
