package domain

import "os"

type AppI interface {
	//groups
	SaveGroup(g *Group) error
	AddUser(userID int, goupID int) error
	//users
	SaveUser(u *User) error
	FindUserByEmail(email string) (user *User, err error)
}

type Service struct {
	App AppI
}

func NewService(repo AppI) *Service {
	return &Service{App: repo}
}

func (s *Service) NewGroup(_name string, _type GroupType) (*Group, error) {
	g := Group{ID: os.Getegid(), Name: _name, Type: _type}
	err := s.App.SaveGroup(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *Service) AddUser(userID, groupID int) error {
	return s.App.AddUser(userID, groupID)
}

func (s *Service) NewUser(name string) (*User, error) {
	u := User{Name: name, ID: os.Getegid()}
	err := s.App.SaveUser(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) FindUserByEmail(email string) (userID *User, err error) {
	return s.App.FindUserByEmail(email)
}
