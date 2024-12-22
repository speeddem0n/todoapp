package todo

type TodoList struct { // Структура для Списка Задач "todo"
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UserList struct { // Структура для списка пользователей
	Id     int
	UserId int
	ListId int
}

type TodoItem struct { // Структура для Задач "todo"
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type ListItem struct {
	Id     int
	ListId int
	ItemId int
}
