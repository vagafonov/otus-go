package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(text string) []string {
	type wordFrequency struct {
		word  string
		count uint8
	}
	splitedText := strings.Fields(text)
	frequenciesMap := make(map[string]wordFrequency)
	wordsForSortSlice := []wordFrequency{}
	result := make([]string, 0, 10)

	// Вычисление частоты слов при помощи map
	// re := regexp.MustCompile(`([а-яА-Я]+|[а-яА-Я-]{2,})`)
	re := regexp.MustCompile(`([a-zA-Zа-яА-Я-]+)`)
	for _, v := range splitedText {
		word := strings.ToLower(string(re.Find([]byte(v))))
		if word == "-" {
			continue
		}

		if value, ok := frequenciesMap[word]; ok {
			value.count++
			frequenciesMap[word] = value
		} else {
			frequenciesMap[word] = wordFrequency{word, 1}
		}
	}

	// Конвертация map в slice для сортировки
	for k := range frequenciesMap {
		wordsForSortSlice = append(wordsForSortSlice, frequenciesMap[k])
	}

	// Сортировка. При одинаковой частоте сравниваются слова
	sort.Slice(wordsForSortSlice, func(i, j int) bool {
		if wordsForSortSlice[i].count == wordsForSortSlice[j].count {
			return wordsForSortSlice[i].word < wordsForSortSlice[j].word
		}

		return wordsForSortSlice[i].count > wordsForSortSlice[j].count
	})

	// Конветация из слайса структур в слайс строк.
	for k := range wordsForSortSlice {
		result = append(result, wordsForSortSlice[k].word)
		if k == 9 {
			break
		}
	}

	return result
}
