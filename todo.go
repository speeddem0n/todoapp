package todo

type Todo struct { // Структура для Списка дел
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UserList struct { // Структура для списка пользователей
	Id     int
	UserId int
	ListId int
}

type TodoItem struct { // Структура для Todo
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
