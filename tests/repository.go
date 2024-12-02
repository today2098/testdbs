package tests

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type PersonsRepositoryInterface interface {
	CreatePerson(c context.Context, person *Person) error
	GetPerson(c context.Context, id string) (*Person, error)
}

type PersonsRepository struct {
	dbx *sqlx.DB
}

var _ PersonsRepositoryInterface = (*PersonsRepository)(nil)

func NewPersonsRepository(dbx *sqlx.DB) *PersonsRepository {
	return &PersonsRepository{dbx}
}

func (r *PersonsRepository) CreatePerson(c context.Context, person *Person) error {
	if _, err := r.dbx.NamedExecContext(c, "INSERT INTO users (id, name, birthday) VALUES (:id, :name, :birthday)", person); err != nil {
		return err
	}
	return nil
}

func (r *PersonsRepository) GetPerson(c context.Context, id string) (*Person, error) {
	var person Person
	if err := r.dbx.GetContext(c, &person, "SELECT id, name, birthday FROM users WHERE id = ?", id); err != nil {
		return nil, err
	}
	return &person, nil
}
