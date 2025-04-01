package uow

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"
)

type UnitOfWork struct {
	DB *sql.DB
	TX *sql.Tx
}

func New() *UnitOfWork {
	db, err := sql.Open("driver name", "source name")
	if err != nil {
		return nil
	}
	return &UnitOfWork{DB: db}
}

func (uow UnitOfWork) WithTx(ctx context.Context, f func(uow UnitOfWork) error) (err error) {
	// если вызов — внутри UnitOfWork с существующей транзакцией, то выполняем внутри неё
	if uow.TX != nil {
		if e := f(uow); e != nil {
			return e
		}
		return nil
	}
	// иначе — создаём и выполняем транзакцию
	return uow.WithBeginTx(ctx, f)
}

func (uow UnitOfWork) WithBeginTx(ctx context.Context, f func(uow UnitOfWork) error) (err error) {
	var newTx *sql.Tx
	// стартуем новую транзакцию
	if tx, e := uow.DB.BeginTx(ctx, nil); e != nil {
		return e
	} else {
		newTx = tx
	}
	nuow := uow
	nuow.TX = newTx
	// во время транзакции может произойти что угодно — обработаем это
	// в случае ошибки или паники — откатим её
	commit := false
	defer func() {
		if r := recover(); r != nil || !commit {
			if r != nil {
				log.Printf("!!! TRANSACTION PANIC !!! : %s\n%s", r,
					string(debug.Stack()))
			}
			if e := newTx.Rollback(); e != nil {
				err = e
			} else if r != nil {
				err = fmt.Errorf("transaction panic: %s", r)
			}
		} else if commit {
			if e := newTx.Commit(); e != nil {
				err = e
			}
		}
	}()
	if e := f(nuow); e != nil {
		return e
	}
	commit = true
	return nil
}

func (uow UnitOfWork) Tx() *sql.Tx {
	return uow.TX
}
