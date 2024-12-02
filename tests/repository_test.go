package tests_test

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/today2098/testdbs/tests"
)

func TestPersonsRepository_CreatePerson(t *testing.T) {
	t.Parallel()

	td, err := h.Create()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer td.Drop()

	dbx := td.DBx()

	type fields struct {
		dbx *sqlx.DB
	}
	type args struct {
		c      context.Context
		person *tests.Person
	}
	tts := []struct {
		name    string
		fields  *fields
		args    *args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Positive case",
			fields: &fields{
				dbx: dbx,
			},
			args: &args{
				c: context.Background(),
				person: &tests.Person{
					Id:       "new_id",
					Name:     "D",
					Birthday: time.Date(2024, time.December, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Negative case",
			fields: &fields{
				dbx: dbx,
			},
			args: &args{
				c: context.Background(),
				person: &tests.Person{
					Id:       "ana1999",
					Name:     "Ana",
					Birthday: time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			r := tests.NewPersonsRepository(tt.fields.dbx)
			err := r.CreatePerson(tt.args.c, tt.args.person)

			tt.wantErr(t, err)
		})
	}
}

func TestPersonsRepository_GetPerson(t *testing.T) {
	t.Parallel()

	td, err := h.Create()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer td.Drop()

	dbx := td.DBx()

	type fields struct {
		dbx *sqlx.DB
	}
	type args struct {
		c  context.Context
		id string
	}
	type want struct {
		person *tests.Person
	}
	tts := []struct {
		name    string
		fields  *fields
		args    *args
		want    *want
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Positive case",
			fields: &fields{
				dbx: dbx,
			},
			args: &args{
				c:  context.Background(),
				id: "ana1999",
			},
			want: &want{
				person: &tests.Person{
					Id:       "ana1999",
					Name:     "Ana",
					Birthday: time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Negative case",
			fields: &fields{
				dbx: dbx,
			},
			args: &args{
				c:  context.Background(),
				id: "unknow_id",
			},
			want: &want{
				person: nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "Negative case",
			fields: &fields{
				dbx: dbx,
			},
			args: &args{
				c:  context.Background(),
				id: "new_id", // NOTE: This id is used in TestPersonsRepository_CreatePerson but safe.
			},
			want: &want{
				person: nil,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			r := tests.NewPersonsRepository(tt.fields.dbx)
			person, err := r.GetPerson(tt.args.c, tt.args.id)

			tt.wantErr(t, err)

			if err == nil {
				assert.Equal(t, tt.want.person.Id, person.Id)
				assert.Equal(t, tt.want.person.Name, person.Name)
				assert.Equal(t, tt.want.person.Birthday, person.Birthday)
			}
		})
	}
}
