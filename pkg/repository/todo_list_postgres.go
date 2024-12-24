package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	todo "github.com/speeddem0n/todoapp"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres { // инициалицируем новую структуру TodoListPostgres которая принимает подключение в БД
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) { // Метод для создания списка возвращает id созданного списка и ошибку
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

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) { // Метод для возвращения всех списков "todo" конкретного пользователя (принимает id пользователя)
	var lists []todo.TodoList // Пустой стайс для списков дел

	query := fmt.Sprintf("SELECT todo_lists.id, todo_lists.title, todo_lists.description FROM %s INNER JOIN %s on todo_lists.id = users_list.list_id WHERE users_list.user_id = $1", todoListTable, usersListsTable) // SQL запрос для получения всех списков конкретного юзера

	err := r.db.Select(&lists, query, userId) // Метод Select для выборки N колл-ва элементов из БД

	return lists, err // Возвращаем списки и ошибку
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) { // Метод для получения списка пользователя по его ID
	var list todo.TodoList // переменная для последующей записи нужного списка

	query := fmt.Sprintf("SELECT todo_lists.id, todo_lists.title, todo_lists.description FROM %s INNER JOIN %s on todo_lists.id = users_list.list_id WHERE users_list.user_id = $1 AND users_list.list_id = $2", todoListTable, usersListsTable) // SQL запрос для получения конкретного списка, конкретного юзера

	err := r.db.Get(&list, query, userId, listId) // Метод Get для выборки 1ой строки из БД

	return list, err // Возвращаем списки и ошибку
}

func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error { // Метод для обновления списка по его id
	setValues := make([]string, 0) // Слайс строк
	args := make([]interface{}, 0) // Слайс interface
	argId := 1                     // Id аргументов

	if input.Title != nil { // Проверка поля Title
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId)) // Добавляем в слайс элементы для формирования запроса в базу. В setValues добавляем присвоение к полю title а после знака "=" записываем значение для placeholder
		args = append(args, *input.Title)                              // В слайс агрументов добовляется само значения поля Title
		argId++                                                        // инкримент id аргумента
	}

	if input.Description != nil { // Проверка поля Description
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId)) // Добавляем в слайс элементы для формирования запроса в базу. В setValues добавляем присвоение к полю Description а после знака "=" записываем значение для placeholder
		args = append(args, *input.Description)                              // В слайс агрументов добовляется само значения поля Description
		argId++                                                              // инкримент id аргумента
	}

	setQuery := strings.Join(setValues, ", ") // Соеденяем элементы setValues в одну строку через запятую

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id = $%d AND ul.user_id = $%d", todoListTable, setQuery, usersListsTable, argId, argId+1) /* SQL запрос обновления в БД.
	Пояснения для sprintf(1: таблица todo_lists, 2: строка setQuery со значениями которые нужно изменить, 3: таблица user_lists, 4: Номер плейсхолдера, 5: Номер плейсхолдера+1)
	*/

	args = append(args, listId, userId) // Добовляем в args id списка и id пользователя

	_, err := r.db.Exec(query, args...) // Выполням запрос методом Exec и передаем в него список аргументов

	return err
}

func (r *TodoListPostgres) Delete(userId, listId int) error { // Метод для удаления списка по его ID
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul where tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2", todoListTable, usersListsTable) // SQL для удаления списка пользователя по его id

	_, err := r.db.Exec(query, userId, listId) // Метод Exec для простого выполнения SQL запроса

	return err // Возвращаем ошибку (или ее отсутствие)
}
