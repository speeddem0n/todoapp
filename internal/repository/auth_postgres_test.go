package repository

import (
	"testing"

	"github.com/speeddem0n/todoapp/internal/models"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type mockBehavior func(user models.User)

	tests := []struct {
		name    string
		user    models.User
		mock    mockBehavior
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			user: models.User{
				Id:       1,
				Name:     "testname",
				Username: "Testusername",
				Password: "testpassword",
			},
			mock: func(user models.User) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(user.Id)

				mock.ExpectQuery("INSERT INTO users").WithArgs(user.Name, user.Username, user.Password).WillReturnRows(rows)
			},
			want: 1,
		},
		{
			name: "Empty Fields",
			user: models.User{
				Id:       1,
				Name:     "testname",
				Username: "Testusername",
				Password: "",
			},
			mock: func(user models.User) {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO users").WithArgs(user.Name, user.Username, user.Password).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.user)

			got, err := r.CreateUser(test.user)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	type args struct {
		username string
		password string
	}

	type mockBehavior func(args args)

	tests := []struct {
		name    string
		input   args
		mock    mockBehavior
		want    models.User
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				username: "testusername",
				password: "testpassword",
			},
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"}).AddRow(1, "testname", "testusername", "testpassword")

				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(args.username, args.password).WillReturnRows(rows)
			},
			want: models.User{
				Id:       1,
				Name:     "testname",
				Username: "testusername",
				Password: "testpassword",
			},
		},
		{
			name: "User Not Found",
			input: args{
				username: "not",
				password: "found",
			},
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"})

				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs(args.username, args.password).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input)

			got, err := r.GetUser(test.input.username, test.input.password)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
