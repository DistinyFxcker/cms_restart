package domain

import (
	"sync"
	"upper.io/db.v3/lib/sqlbuilder"
)

type TnxAction func(ctx Context) error

type Context interface{
	Copy() Context

	DB() sqlbuilder.SQLBuilder

	Tx(fn TnxAction) error

	Set(key string , val interface{})

	Get(key string) (interface{}, bool)

	DBType()string
}

func NewContext(db sqlbuilder.Database, dbType string) *dbContext {
	return &dbContext{
		db:		db,
		data:	make(map[string]interface{}),
		dbType: dbType,
		locker: &sync.Mutex{},
	}
}

//事务上下文
type txContext struct {
	db   sqlbuilder.Database
	tx   sqlbuilder.Tx
	data map[string]interface{}
	dbType string
}

func (ctx *txContext) Copy() Context {
	return NewContext(ctx.db,ctx.dbType)
}

func (ctx *txContext) DB() sqlbuilder.SQLBuilder {
	return ctx.tx
}

func (ctx *txContext) Tx(action TnxAction) error {
	return action(ctx)
}

func (ctx *txContext) Set(key string, val interface{}) {
	ctx.data[key] = val
}

func (ctx *txContext) Get(key string) (interface{}, bool) {
	val, ok := ctx.data[key]
	if ok {
		return val, true
	}
	return nil, false
}
func (ctx *txContext)DBType()string{
	return ctx.dbType
}

//数据库上下文
type dbContext struct {
	db sqlbuilder.Database
	data map[string]interface{}
	locker sync.Locker
	dbType string
}



