package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	todo "github.com/speeddem0n/todoapp"
)

type TodoItemPostgres struct { // TodoItemPostgres с полем подключения к BD
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres { // инициалицируем новую структуру TodoItemPostgress которая принимает подключение в БД
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) { // Метод для создания "todo" элемента возвращает id созданного элемента и ошибку
	tx, err := r.db.Begin() // Начинаем SQL транзакцию
	if err != nil {
		return 0, err
	}
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable) // SQL зарос для создания элемента "todo" в таблице todo_items, вовращает id созданного элемента

	row := tx.QueryRow(createItemQuery, item.Title, item.Description) // Выполняем первый запрос в транзакции
	err = row.Scan(&itemId)                                           // Записываем возвращенный id в переменную
	if err != nil {
		tx.Rollback() // В случае ошибки откатываем транзакцию методом Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable) // SQL зарос что бы связать id списка и id нового "todo" элемента списка

	_, err = tx.Exec(createListItemsQuery, listId, itemId) // Метод Exec для простого выполнения SQL запроса
	if err != nil {
		tx.Rollback() // В случае ошибки откатываем транзакцию методом Rollback()
		return 0, err
	}

	return itemId, tx.Commit() // Commit() применяет изменения к базе данных и заканчивает транзакцию
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) { // Метод для возвращения всех элементов списка конкретного пользователя (принимает id пользователя и списка)
	var items []todo.TodoItem // Стурктура для записи ответа
	query := fmt.Sprintf(`SELECT ti.title, ti.description, ti.done FROM %s ti  
	INNER JOIN %s li on li.item_id = ti.id 
	INNER JOIN %s ul on ul.list_id = li.list_id
	WHERE li.list_id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable) // SQL зарос для выборки всех элементов "todo" из списка

	err := r.db.Select(&items, query, listId, userId) // Метод Select для выборки N колл-ва элементов из БД
	if err != nil {
		return nil, err
	}

	return items, nil
}
