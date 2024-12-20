package todo

type User struct { // Структура для пользователя
	Id       int    `json:"-"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
