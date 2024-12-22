package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	todo "github.com/speeddem0n/todoapp"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres { // инициалицируем новую структуру AuthPostgres которая принимает подключение в БД
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin() // Begin() начинает sql транзакцию
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListTable) // SQL зарос для создания списка в таблице todo_list, вовращает id созданного списка
	row := tx.QueryRow(createListQuery, list.Title, list.Description)                                                 // Выполняем первый запрос в транзакции
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback() // В случае ошибки откатываем транзакцию методом Rollback()
		return 0, err
	}

	createUserListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable) // SQL зарос что бы связать id Пользователя и id новго списка
	_, err = tx.Exec(createUserListQuery, userId, id)                                                        // Метод Exec для простого выполнения SQL запроса
	if err != nil {
		tx.Rollback() // В случае ошибки откатываем транзакцию методом Rollback()
		return 0, err
	}

	return id, tx.Commit() // Commit() применяет изменения к базе данных и заканчивает транзакцию
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList // Пустой стайс для списков дел

	query := fmt.Sprintf("SELECT todo_lists.id, todo_lists.title, todo_lists.description FROM %s INNER JOIN %s on todo_lists.id = users_list.list_id WHERE users_list.user_id = $1", todoListTable, usersListsTable) // SQL запрос для получения всех списков конкретного юзера

	err := r.db.Select(&lists, query, userId) // Метод Select для выборки N колл-ва элементов из БД

	return lists, err // Возвращаем списки и ошибку
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) { // Метод для получения конкретного списка пользователя по его ID
	var list todo.TodoList // переменная для последующей записи нужного списка

	query := fmt.Sprintf("SELECT todo_lists.id, todo_lists.title, todo_lists.description FROM %s INNER JOIN %s on todo_lists.id = users_list.list_id WHERE users_list.user_id = $1 AND users_list.list_id = $2", todoListTable, usersListsTable) // SQL запрос для получения конкретного списка, конкретного юзера

	err := r.db.Get(&list, query, userId, listId) // Метод Get для выборки 1ой строки из БД

	return list, err // Возвращаем списки и ошибку
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul where tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2", todoListTable, usersListsTable) // SQL для удаления списка пользователя по его id

	_, err := r.db.Exec(query, userId, listId) // Метод Exec для простого выполнения SQL запроса

	return err // Возвращаем ошибку (или ее отсутствие)
}
