package domain

import (
	"context"
	"os"
)

type AppI interface {
	//groups
	SaveGroup(ctx context.Context, g *Group) error
	AddUser(ctx context.Context, userID int, goupID int) error
	//users
	SaveUser(ctx context.Context, u *User) error
	FindUserByEmail(ctx context.Context, email string) (user *User, err error)
}

type Service struct {
	App AppI
}

func NewService(repo AppI) *Service {
	return &Service{App: repo}
}

func (s *Service) NewGroup(ctx context.Context, _name string, _type GroupType) (*Group, error) {
	g := Group{ID: os.Getegid(), Name: _name, Type: _type}
	err := s.App.SaveGroup(ctx, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *Service) AddUser(ctx context.Context, userID, groupID int) error {
	return s.App.AddUser(ctx, userID, groupID)
}

func (s *Service) NewUser(ctx context.Context, name string) (*User, error) {
	u := User{Name: name, ID: os.Getegid()}
	err := s.App.SaveUser(ctx, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Service) FindUserByEmail(ctx context.Context, email string) (userID *User, err error) {
	return s.App.FindUserByEmail(ctx, email)
}
