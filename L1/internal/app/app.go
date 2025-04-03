package app

import (
	"L1/internal/domain"
	"context"
)

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
	SaveGroup(ctx context.Context, g *domain.Group) error
	AddUser(ctx context.Context, userID, goupID int) error
	//users
	SaveUser(ctx context.Context, u *domain.User) error
	FindUserByEmail(ctx context.Context, email string) (user *domain.User, err error)
}

type App struct {
	DB DBI
}

func NewApp(db DBI) *App {
	return &App{DB: db}
}

func (a *App) SaveGroup(ctx context.Context, g *domain.Group) error {
	return a.DB.SaveGroup(ctx, g)
}

func (a *App) AddUser(ctx context.Context, userID, groupID int) error {
	return a.DB.AddUser(ctx, userID, groupID)
}

func (a *App) SaveUser(ctx context.Context, u *domain.User) error {
	return a.DB.SaveUser(ctx, u)
}

func (a *App) FindUserByEmail(ctx context.Context, email string) (user *domain.User, err error) {
	return a.DB.FindUserByEmail(ctx, email)
}
