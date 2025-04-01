package app

import "L1/internal/domain"

// type Repo interface {
// 	//groups
// 	SaveGroup(g *Group) error
// 	AddUser(userID int, goupID int) error
// 	//users
// 	SaveUser(u *User) error
// 	FindUserByEmail(email string) (user *User, err error)
// }

type DBI interface {
	//groups
	SaveGroup(g *domain.Group) error
	AddUser(userID, goupID int) error
	//users
	CreateUser(u *domain.User) error
	FindUserByEmail(email string) (user *domain.User, err error)
}

type App struct {
	DB DBI
}

func NewApp(db DBI) *App {
	return &App{DB: db}
}

func (a *App) SaveGroup(g *domain.Group) error {
	return a.DB.SaveGroup(g)
}

func (a *App) AddUser(userID, groupID int) error {
	return a.DB.AddUser(userID, groupID)
}

func (a *App) SaveUser(u *domain.User) error {
	return a.DB.CreateUser(u)
}

func (a *App) FindUserByEmail(email string) (user *domain.User, err error) {
	return a.DB.FindUserByEmail(email)
}
