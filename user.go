package todo

type User struct { // Структура для пользователя
	Id       int    `json:"-" db :"id"`              // db тэг для базы данных
	Name     string `json:"name" binding:"required"` // Тег building Валидирует наличие данного поля в теле запроса
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
