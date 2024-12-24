package todo

import "errors"

type TodoList struct { // Структура для Списка Задач "todo"
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UserList struct { // Структура для связавания пользователя и его списков
	Id     int
	UserId int
	ListId int
}

type TodoItem struct { // Структура для Задач "todo"
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ListItem struct { // Структура для связавания списков и его "todo"
	Id     int
	ListId int
	ItemId int
}

type UpdateListInput struct { // Структура для обновления списка list
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateListInput) Validate() error { // Метод для проверки вадлиности структуры для обнавления списка UpdateListInput
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type UpdateItemInput struct { // Структура для обновления элемента списка
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateItemInput) Validate() error { // Метод для проверки вадлиности структуры для обнавления списка UpdateItemInput
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
