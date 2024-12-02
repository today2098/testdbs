package tests

import "time"

type Person struct {
	Id       string    `db:"id"`
	Name     string    `db:"name"`
	Birthday time.Time `db:"birthday"`
}
