package models

type CreateTodoItemInp struct { // Структура для Задач "todo"
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}
