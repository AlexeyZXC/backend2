package db

type DB struct {
	// Um UserMapper
	UserMapper
	// Gm GroupMapper
	GroupMapper
}

func NewDB() *DB {
	return &DB{}
}
