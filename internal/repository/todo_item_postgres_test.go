package repository

import (
	"database/sql"
	"errors"
	"testing"

	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/internal/models"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db) // Подменяем подключение к бд моком

	type args struct {
		listId int
		item   todo.TodoItem
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test Description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "test Description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("Empty fields"))
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd insert error",
			args: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test Description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnError(errors.New("Insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.id)

			got, err := r.Create(test.args.listId, test.args.item)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.id, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []todo.TodoItem
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true).
					AddRow(2, "title2", "description2", false).
					AddRow(3, "title3", "description3", false)

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_list ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			want: []todo.TodoItem{
				{Id: 1, Title: "title1", Description: "description1", Done: true},
				{Id: 2, Title: "title2", Description: "description2", Done: false},
				{Id: 3, Title: "title3", Description: "description3", Done: false},
			},
		},
		{
			name: "No Records",
			mock: func() {

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_list ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetAll(test.input.userId, test.input.listId)
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

func TestTodoItemPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    todo.TodoItem
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				row := sqlmock.NewRows([]string{"id", "title", "description", "done"}).AddRow(1, "title1", "description1", true)

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_list ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(row)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			want: todo.TodoItem{Id: 1, Title: "title1", Description: "description1", Done: true},
		},
		{
			name: "No Records",
			mock: func() {
				row := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_list ul on (.+) WHERE (.+)").
					WithArgs(404, 1).WillReturnRows(row)
			},
			input: args{
				listId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetById(test.input.userId, test.input.listId)
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

func TestTodoItemPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		itemId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("DELETE FROM todo_items ti USING lists_items li, users_list ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
			},
		},
		{
			name: "No Records",
			mock: func() {
				mock.ExpectExec("DELETE FROM todo_items ti USING lists_items li, users_list ul WHERE (.+)").
					WithArgs(404, 1).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				itemId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := r.Delete(test.input.userId, test.input.itemId)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		itemId int
		userId int
		input  models.UpdateItemInput
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK_ALL_FIELDS",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_list ul WHERE (.+)").
					WithArgs("new title", "new description", true, 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: models.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
					Done:        boolPointer(true),
				},
			},
		},
		{
			name: "OK_WithoutDone",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_list ul WHERE (.+)").
					WithArgs("new title", "new description", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: models.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_WithoutDoneAndDescription",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_list ul WHERE (.+)").
					WithArgs("new title", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: models.UpdateItemInput{
					Title: stringPointer("new title"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET FROM lists_items li, users_list ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := r.Update(test.input.userId, test.input.itemId, test.input.input)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
