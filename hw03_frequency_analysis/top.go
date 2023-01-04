package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(text string) []string {
	type word struct {
		word  string
		count uint8
	}
	wordsText := strings.Fields(text)
	wordsMap := make(map[string]uint8)
	wordsStructured := []word{}
	result := make([]string, 0)

	// Вычисление частоты слов при помощи map
	// re := regexp.MustCompile(`([а-яА-Я]+|[а-яА-Я-]{2,})`)
	re := regexp.MustCompile(`([а-яА-Я-]+)`)
	for _, v := range wordsText {
		t := strings.ToLower(string(re.Find([]byte(v))))
		if t == "-" {
			continue
		}
		wordsMap[t]++
	}

	// Конвертация map в slice
	for k, v := range wordsMap {
		wordsStructured = append(wordsStructured, struct {
			word  string
			count uint8
		}{
			k,
			v,
		},
		)
	}

	// Сортировка. При одинаковой частоте сравниваются слова
	sort.Slice(wordsStructured, func(i, j int) bool {
		if wordsStructured[i].count == wordsStructured[j].count {
			return wordsStructured[i].word < wordsStructured[j].word
		}

		return wordsStructured[i].count > wordsStructured[j].count
	})

	// Результат slice слов
	for k := range wordsStructured {
		result = append(result, wordsStructured[k].word)
		if k == 9 {
			break
		}
	}

	return result
}
