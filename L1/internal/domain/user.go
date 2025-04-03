package domain

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	PW        string
	CreatedAt time.Time
}
