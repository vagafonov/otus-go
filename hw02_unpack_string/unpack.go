package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	input := []rune(str)
	var builder strings.Builder
	var symbol string
	var prevSymbol string

	for i := 0; i < len(input); i++ {
		symbol += string(input[i])

		if symbol == "\\" {
			continue
		}

		if symbol == "0" && i != 0 {
			continue
		}

		number, err := strconv.Atoi(symbol)

		// Ошибка в случае начинающийся строки с цифры
		if i == 0 && err == nil {
			return "", ErrInvalidString
		}

		if len(input)-1 != i {
			_, nextNumberErr := strconv.Atoi(string(input[i+1]))
			// Текущий и следующий символ - цифры
			if nextNumberErr == nil && err == nil {
				return "", ErrInvalidString
			}

			// Удаление при распаковке
			if string(input[i+1]) == "0" {
				symbol = ""
				continue
			}
		}

		if err == nil { // Текущий сивол число
			// Дублирование последнего символа (распаковка)
			for j := 0; j < number-1; j++ {
				builder.WriteString(prevSymbol)
			}
			continue
		}

		// Текущий символ строка. Добавление в результат
		prevSymbol = string(input[i])
		builder.WriteString(prevSymbol)
		symbol = ""
	}

	return builder.String(), nil
}
