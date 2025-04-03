package db

type DB struct {
	UserMapper
	GroupMapper
}

func NewDB() *DB {
	return &DB{}
}
