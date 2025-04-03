package db

import (
	"L1/internal/domain"
	"L1/internal/uow"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type GroupMapper struct {
	uow *uow.UnitOfWork
}

func NewGroup() *GroupMapper {
	uow := uow.New()
	if uow == nil {
		return nil
	}
	return &GroupMapper{uow: uow}
}

func (gm *GroupMapper) Columns() string {
	return "id, name, type"
}

func (gm *GroupMapper) Fields(g *domain.Group) []any {
	return []any{&g.ID, &g.Name, &g.Type}
}

func (gm *GroupMapper) SelectGroup(ctx context.Context, c string, v ...any) (*domain.Group, error) {
	g := &domain.Group{}

	// такой вызов либо использует уже запущенную ранее в uow транзакцию, либо запустит новую
	if err := gm.uow.WithTx(ctx, func(uowtx uow.UnitOfWork) error {
		// внутри этой функции нам ничто не мешает получить транзакцию WithTx ещё раз
		// это полезно, если дальше вызываются какие-то универсальные функции наподобие этой
		return uowtx.Tx().QueryRowContext(ctx,
			fmt.Sprintf(`SELECT %s FROM groups WHERE %s`, gm.Columns(), c),
			v...,
		).Scan(gm.Fields(g)...)
	}); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return g, nil
}

// func (gm *GroupMapper) CreateGroup(ctx context.Context, name string, groupType domain.GroupType) (*domain.Group, error) {
func (gm *GroupMapper) SaveGroup(ctx context.Context, g *domain.Group) error {
	// g := &domain.Group{}

	if err := gm.uow.WithTx(ctx, func(uowtx uow.UnitOfWork) error {
		return uowtx.Tx().QueryRowContext(ctx,
			fmt.Sprintf(`INSERT INTO groups(name, groupType)
					 	VALUES($1, $2)
	 					RETURNING %s`, gm.Columns()), g.Name, g.Type).Scan(gm.Fields(g))
	}); err != nil {
		return err
	}
	return nil
}

func (gm *GroupMapper) AddUser(ctx context.Context, userID, groupID int) error {
	if err := gm.uow.WithTx(ctx, func(uowtx uow.UnitOfWork) error {
		//check if user exist
		//check if group exist
		res, err := uowtx.Tx().Exec(
			`insert into relations (userID, groupID) values($1, $2)`,
			userID, groupID,
		)
		if err != nil {
			return err
		}
		if aff, err := res.RowsAffected(); aff == 0 || err != nil {
			if err != nil {
				return err
			}
			if aff == 0 {
				return errors.New("no affected rows")
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
