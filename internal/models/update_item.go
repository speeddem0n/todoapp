package models

import "errors"

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
