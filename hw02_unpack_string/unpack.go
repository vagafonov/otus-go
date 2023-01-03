package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	input := []rune(str)
	output := []string{}
	var symbol string

	for i := 0; i < len(input); i++ {
		symbol += string(input[i])

		if symbol == "\\" {
			continue
		}

		number, err := strconv.Atoi(symbol)

		// Ошибка в случае начинающийся строки с цифры
		if i == 0 && err == nil {
			return "", ErrInvalidString
		}

		// Ошибка в случае двух цифрах подряд
		if len(input)-1 != i {
			_, nextNumberErr := strconv.Atoi(string(input[i+1]))
			// Текущий и следующий символ - цифры
			if nextNumberErr == nil && err == nil {
				return "", ErrInvalidString
			}
		}

		if err == nil { // Текущий сивол число
			if number == 0 { // Удаление последенего символа
				output = output[:len(output)-1]
			} else { // Дублирование последнего символа (распаковка)
				for j := 0; j < number-1; j++ {
					last := output[len(output)-1]
					output = append(output, last)
				}
			}
		} else { // Текущий символ строка. Добавление в результат
			output = append(output, string(input[i]))
		}

		symbol = ""
	}

	return strings.Join(output, ""), nil
}
