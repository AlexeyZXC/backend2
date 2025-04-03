package db

import (
	"L1/internal/domain"
	"L1/internal/uow"
	"context"
	"database/sql"
	"fmt"
)

type UserMapper struct {
	uow *uow.UnitOfWork
}

func NewUser() *UserMapper {
	uow := uow.New()
	if uow == nil {
		return nil
	}
	return &UserMapper{uow: uow}
}

func (um *UserMapper) Columns() string {
	return "id, name, email, created_at"
}

func (um *UserMapper) Fields(u *domain.User) []any {
	return []any{&u.ID, &u.Name, &u.Email, &u.CreatedAt}
}

func (um *UserMapper) SelectUser(ctx context.Context, c string, v ...any) (*domain.User, error) {
	u := &domain.User{}

	// такой вызов либо использует уже запущенную ранее в uow транзакцию, либо запустит новую
	if err := um.uow.WithTx(ctx, func(uowtx uow.UnitOfWork) error {
		// внутри этой функции нам ничто не мешает получить транзакцию WithTx ещё раз
		// это полезно, если дальше вызываются какие-то универсальные функции наподобие этой
		return uowtx.Tx().QueryRowContext(ctx,
			fmt.Sprintf(`SELECT %s FROM users WHERE %s`, um.Columns(), c),
			v...,
		).Scan(um.Fields(u)...)
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// func (um *UserMapper) CreateUser(ctx context.Context, name, email, pw string) (*domain.User, error) {
func (um *UserMapper) SaveUser(ctx context.Context, u *domain.User) error {
	// u := &domain.User{}

	if err := um.uow.WithTx(ctx, func(uowtx uow.UnitOfWork) error {
		return uowtx.Tx().QueryRowContext(ctx,
			fmt.Sprintf(`INSERT INTO users(name, email, password)
					 	VALUES($1, $2, crypt($3, gen_salt('bf', 8)))
	 					RETURNING %s`, um.Columns()), u.Name, u.Email, u.PW).Scan(um.Fields(u))
	}); err != nil {
		return err
	}
	return nil
}

func (um *UserMapper) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return um.SelectUser(ctx,
		fmt.Sprintf(`select %s from users where lower(email) = lower($1)`, um.Columns()),
		email)
}
