package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	todo "github.com/speeddem0n/todoapp"
)

type TodoItemPostgres struct { // TodoItemPostgres с полем подключения к БД
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres { // Инициализируем новую структуру TodoItemPostgress которая принимает подключение к БД
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
	var items []todo.TodoItem // Слайс стурктур для записи ответа
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti  
	INNER JOIN %s li on li.item_id = ti.id 
	INNER JOIN %s ul on ul.list_id = li.list_id
	WHERE li.list_id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable) // SQL зарос для выборки всех элементов "todo" из списка

	err := r.db.Select(&items, query, listId, userId) // Метод Select для выборки N колл-ва элементов из БД
	if err != nil {
		return nil, err
	}

	return items, nil // Возвращаем срез элементов
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) { // Метод для получения элемента списка по его ID
	var item todo.TodoItem // Стурктура для записи ответа
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti  
	INNER JOIN %s li on li.item_id = ti.id 
	INNER JOIN %s ul on ul.list_id = li.list_id
	WHERE ti.id = $1 AND ul.user_id = $2`, todoItemsTable, listsItemsTable, usersListsTable) // SQL зарос для выборки всех элементов "todo" из списка

	err := r.db.Get(&item, query, itemId, userId) // Метод Get для выборки одного элемента из БД
	if err != nil {
		return item, err
	}

	return item, nil // Возвращаем полученный элемент
}

func (r *TodoItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error { // Метод для обновления элемента списка по его id
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

	if input.Done != nil { // Проверка поля Done
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId)) // Добавляем в слайс элементы для формирования запроса в базу. В setValues добавляем присвоение к полю Done а после знака "=" записываем значение для placeholder
		args = append(args, *input.Done)                              // В слайс агрументов добовляется само значения поля Done
		argId++                                                       // инкримент id аргумента
	}

	setQuery := strings.Join(setValues, ", ") // Соеденяем элементы setValues в одну строку через запятую

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
	WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d 
	AND ti.id =$%d`, todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	/* SQL запрос обновления в БД.
	Пояснения для sprintf(1: таблица todo_lists, 2: строка setQuery со значениями которые нужно изменить,
	3: таблица list_items, 4: таблица user_lists, 5: Номер плейсхолдера, 6: Номер плейсхолдера+1)
	*/

	args = append(args, userId, itemId) // Добовляем в args id элемента и id пользователя

	_, err := r.db.Exec(query, args...) // Выполням запрос методом Exec и передаем в него список аргументов

	return err
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error { // Метод для удаления эелемнта списка по его ID
	query := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul WHERE li.item_id = ti.id AND ul.list_id = li.list_id AND ti.id = $1 AND ul.user_id = $2", todoItemsTable, listsItemsTable, usersListsTable) // SQL зарос для удаления эелемнта списка по его ID

	_, err := r.db.Exec(query, itemId, userId) // Метод Exec для простого выполнения SQL запроса

	return err // Возвращаем ошибку (или ее отсутствие)
}
