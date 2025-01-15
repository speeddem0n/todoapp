package models

import "errors"

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
