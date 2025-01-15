package todo

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
