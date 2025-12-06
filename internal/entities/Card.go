package entities

import "errors"

type Card struct {
	//Photos
	Name      string
	Condition string
	Modifiers string
	Tags      string
}

func CardValidating(c Card) error {
	if len(c.Name) == 0 {
		return errors.New("ошибка заполнения карточки, пустое имя")
	} else if len(c.Condition) == 0 {
		return errors.New("ошибка заполнения карточки, пустое состояние")
	} else if len(c.Modifiers) == 0 {
		return errors.New("ошибка заполнения карточки, пустые модификации")
	} else if len(c.Tags) == 0 {
		return errors.New("ошибка заполнения карточки, пустые теги")
	} else {
		return nil
	}
}
