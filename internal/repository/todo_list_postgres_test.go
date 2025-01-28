package repository

import (
	"database/sql"
	"testing"

	"github.com/speeddem0n/todoapp/internal/models"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoListPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db) // Подменяем подключение к бд моком

	type args struct {
		userId int
		list   models.CreateListInput
	}

	type mockBehavior func(input args, id int)

	testTable := []struct {
		name    string
		mock    mockBehavior
		input   args
		id      int
		wantErr bool
	}{
		{
			name: "OK",
			mock: func(input args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_lists").WithArgs(input.list.Title, input.list.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_list").WithArgs(input.userId, id).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: args{
				userId: 1,
				list: models.CreateListInput{
					Title:       "new title",
					Description: "new description",
				},
			},
			id: 2,
		},
		{
			name: "Empty Fields",
			mock: func(input args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO todo_lists").
					WithArgs(input.list.Title, input.list.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			input: args{
				userId: 1,
				list: models.CreateListInput{
					Title:       "",
					Description: "description",
				},
			},
			wantErr: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input, test.id)

			got, err := r.Create(test.input.userId, test.input.list)
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

func TestTodoListPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type mockBehavior func(userId int)

	tests := []struct {
		name    string
		mock    mockBehavior
		userId  int
		want    []models.TodoList
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1").
					AddRow(2, "title2", "description2").
					AddRow(3, "title3", "description3")

				mock.ExpectQuery("SELECT (.+) FROM todo_lists INNER JOIN users_list on (.+) WHERE (.+)").
					WithArgs(userId).WillReturnRows(rows)
			},
			userId: 2,
			want: []models.TodoList{
				{Id: 1, Title: "title1", Description: "description1"},
				{Id: 2, Title: "title2", Description: "description2"},
				{Id: 3, Title: "title3", Description: "description3"},
			},
		},
		{
			name: "No Records",
			mock: func(userId int) {

				mock.ExpectQuery("SELECT (.+) FROM todo_lists INNER JOIN users_list on (.+) WHERE (.+)").
					WithArgs(userId).WillReturnError(sql.ErrNoRows)
			},
			userId:  404,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.userId)

			got, err := r.GetAll(test.userId)
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

func TestTodoListPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type mockBehavior func(userId, listId int)

	type args struct {
		userId int
		listId int
	}

	tests := []struct {
		name    string
		mock    mockBehavior
		args    args
		want    models.TodoList
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(userId, listId int) {
				row := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(2, "title1", "description1")

				mock.ExpectQuery("SELECT (.+) FROM todo_lists INNER JOIN users_list on (.+) WHERE (.+)").
					WithArgs(userId, listId).WillReturnRows(row)
			},
			args: args{
				userId: 2,
				listId: 2,
			},
			want: models.TodoList{Id: 2, Title: "title1", Description: "description1"},
		},
		{
			name: "No Records",
			mock: func(userId, listId int) {

				mock.ExpectQuery("SELECT (.+) FROM todo_lists INNER JOIN users_list on (.+) WHERE (.+)").
					WithArgs(userId, listId).WillReturnError(sql.ErrNoRows)
			},
			args: args{
				userId: 404,
				listId: 2,
			},
			wantErr: true,
		},
		{
			name: "Not Found",
			mock: func(userId, listId int) {
				row := sqlmock.NewRows([]string{"id", "title", "description"})

				mock.ExpectQuery("SELECT (.+) FROM todo_lists INNER JOIN users_list on (.+) WHERE (.+)").
					WithArgs(userId, listId).WillReturnRows(row)
			},
			args: args{
				userId: 2,
				listId: 404,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.args.userId, test.args.listId)

			got, err := r.GetById(test.args.userId, test.args.listId)
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

func TestTodoListPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type mockBehavior func(userId, listId int)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(userId, listId int) {

				mock.ExpectExec("DELETE FROM todo_lists tl USING users_list ul WHERE (.+)").
					WithArgs(userId, listId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
			},
		},
		{
			name: "No Records",
			mock: func(userId, listId int) {
				mock.ExpectExec("DELETE FROM todo_lists tl USING users_list ul WHERE (.+)").
					WithArgs(userId, listId).WillReturnError(sql.ErrNoRows)
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
			test.mock(test.input.userId, test.input.listId)

			err := r.Delete(test.input.userId, test.input.listId)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoListPostgres(db)

	type args struct {
		listId    int
		userId    int
		userInput models.UpdateListInput
	}

	type mockBehavior func(userId, listId int, input models.UpdateListInput)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		wantErr bool
	}{
		{
			name: "OK_ALL_FIELDS",
			mock: func(userId, listId int, input models.UpdateListInput) {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_list ul WHERE (.+)").
					WithArgs(input.Title, input.Description, userId, listId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				userInput: models.UpdateListInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_WithoutDescription",
			mock: func(userId, listId int, input models.UpdateListInput) {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_list ul WHERE (.+)").
					WithArgs(input.Title, userId, listId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				userInput: models.UpdateListInput{
					Title: stringPointer("new title"),
				},
			},
		},
		{
			name: "OK_WithoutTItle",
			mock: func(userId, listId int, input models.UpdateListInput) {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_list ul WHERE (.+)").
					WithArgs(input.Description, userId, listId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				userInput: models.UpdateListInput{
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func(userId, listId int, input models.UpdateListInput) {
				mock.ExpectExec("UPDATE todo_lists tl SET FROM users_list ul WHERE (.+)").
					WithArgs(userId, listId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input.userId, test.input.listId, test.input.userInput)

			err := r.Update(test.input.userId, test.input.listId, test.input.userInput)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
