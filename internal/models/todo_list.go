package models

type TodoList struct { // Структура для Списка Задач "todo"
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}
